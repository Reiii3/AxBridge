// AxBridge browser client (ES module, no dependencies).
//
// Mirrors the `kernelsu` package, so code written for it keeps working:
//
//   import { exec, spawn, fullScreen, toast } from "axbridge";
//
// Extra, AxBridge only: setBaseUrl, getToken, setToken, clearToken, ping, BridgeError.

let BASE = "http://127.0.0.1:8787";
const TOKEN_KEY = "axbridge_token";

// ---------------------------------------------------------------------------
// Errors
// ---------------------------------------------------------------------------

/**
 * Thrown (or emitted as the 'error' event of spawn) when the bridge itself fails.
 * `code` is one of: "network" | "unauthorized" | "forbidden" | "http".
 * A command that runs and fails is NOT an error, it resolves with a non-zero `errno`.
 */
export class BridgeError extends Error {
  constructor(code, message) {
    super(message);
    this.name = "BridgeError";
    this.code = code;
  }
}

// ---------------------------------------------------------------------------
// Connection settings
// ---------------------------------------------------------------------------

/** Point the client at another address (default http://127.0.0.1:8787). */
export const setBaseUrl = (url) => {
  BASE = url.replace(/\/+$/, "");
};

// localStorage can be blocked (private mode, sandboxed iframe) or missing (Node), so keep a fallback.
let memoryToken = "";

export const getToken = () => {
  try {
    return localStorage.getItem(TOKEN_KEY) || memoryToken;
  } catch {
    return memoryToken;
  }
};

export const setToken = (token) => {
  memoryToken = String(token).trim();
  try {
    localStorage.setItem(TOKEN_KEY, memoryToken);
  } catch {
    /* storage unavailable, the in-memory copy still works for this page load */
  }
};

export const clearToken = () => {
  memoryToken = "";
  try {
    localStorage.removeItem(TOKEN_KEY);
  } catch {
    /* nothing to clear */
  }
};

/** true when the AxBridge server answers on this device. */
export async function ping() {
  try {
    return (await fetch(`${BASE}/health`)).ok;
  } catch {
    return false;
  }
}

// ---------------------------------------------------------------------------
// Requests
// ---------------------------------------------------------------------------

async function post(path, body, signal) {
  let res;
  try {
    res = await fetch(`${BASE}${path}`, {
      method: "POST",
      headers: { "Content-Type": "application/json", Authorization: `Bearer ${getToken()}` },
      body: JSON.stringify(body),
      signal,
    });
  } catch (error) {
    if (error && error.name === "AbortError") throw error;
    throw new BridgeError("network", "AxBridge is not reachable");
  }

  if (res.status === 401) throw new BridgeError("unauthorized", "Wrong or missing AxBridge token");
  if (res.status === 403) throw new BridgeError("forbidden", "This website is not in AxBridge's allowed list");
  if (!res.ok) throw new BridgeError("http", `AxBridge HTTP ${res.status}: ${await res.text()}`);
  return res;
}

/**
 * Runs a shell command and waits for it to finish.
 * @param {string} command
 * @param {{cwd?: string, env?: Record<string,string>, timeout?: number}} [options]
 *   timeout is in seconds (default 30, max 120).
 * @returns {Promise<{errno: number, stdout: string, stderr: string}>}
 */
export async function exec(command, options = {}) {
  const { cwd, env, timeout } = options;
  const res = await post("/exec", { cmd: command, cwd, env, timeout });
  return res.json();
}

// ---------------------------------------------------------------------------
// spawn
// ---------------------------------------------------------------------------

/** Minimal event emitter, enough for 'data', 'exit' and 'error'. */
class Emitter {
  constructor() {
    this.listeners = {};
  }
  on(event, listener) {
    (this.listeners[event] ||= []).push(listener);
    return this;
  }
  emit(event, ...args) {
    (this.listeners[event] || []).forEach((listener) => listener(...args));
  }
}

class ChildProcess extends Emitter {
  constructor(abort) {
    super();
    this.stdout = new Emitter();
    this.stderr = new Emitter();
    this._abort = abort;
  }

  /** Stops the process (the server kills its whole process group). 'exit' fires with code null. */
  kill() {
    this._abort.abort();
  }
}

/**
 * Starts a process and streams its output.
 *
 *   const child = spawn("ls", ["-lh", "/data"]);
 *   child.stdout.on("data", (text) => ...);
 *   child.stderr.on("data", (text) => ...);
 *   child.on("exit", (code) => ...);   // code is null when killed or timed out by the client
 *   child.on("error", (error) => ...); // BridgeError, e.g. wrong token
 *
 * With a non-empty `args` array the command is executed directly (no shell, no quoting issues).
 * Without args the command is run as a shell line, like exec().
 *
 * @param {string} command
 * @param {string[] | object} [args] or the options object
 * @param {{cwd?: string, env?: Record<string,string>, timeout?: number}} [options]
 *   timeout is in seconds (default 300, max 3600).
 * @returns {ChildProcess}
 */
export function spawn(command, args, options) {
  if (args !== undefined && !Array.isArray(args)) {
    options = args;
    args = undefined;
  }
  const { cwd, env, timeout } = options || {};
  const body = { cmd: command, args: args && args.length > 0 ? args : undefined, cwd, env, timeout };

  const abort = new AbortController();
  const child = new ChildProcess(abort);
  run(child, body, abort.signal);
  return child;
}

async function run(child, body, signal) {
  const decoders = { out: new TextDecoder(), err: new TextDecoder() };
  const streams = { out: child.stdout, err: child.stderr };
  let exited = false;

  const exit = (code) => {
    if (exited) return;
    exited = true;
    child.emit("exit", code);
  };

  try {
    const res = await post("/spawn", body, signal);
    const reader = res.body.getReader();
    // Cancelling the reader closes the connection right away; abort() alone can lag in some runtimes.
    signal.addEventListener("abort", () => reader.cancel().catch(() => {}));
    const lines = new TextDecoder();
    let pending = "";

    const handle = (line) => {
      if (!line) return;
      const event = JSON.parse(line);

      if (event.t === "exit") return exit(event.code);

      const bytes = Uint8Array.from(atob(event.b), (c) => c.charCodeAt(0));
      // stream:true keeps a multi-byte character that was split across chunks intact
      const text = decoders[event.t].decode(bytes, { stream: true });
      if (text) streams[event.t].emit("data", text);
    };

    for (;;) {
      const { done, value } = await reader.read();
      if (done) break;
      pending += lines.decode(value, { stream: true });
      const parts = pending.split("\n");
      pending = parts.pop(); // last piece may be an unfinished line
      parts.forEach(handle);
    }
    handle(pending);
    exit(null); // stream ended without an exit line (connection cut)
  } catch (error) {
    if (error && error.name === "AbortError") return exit(null);
    child.emit("error", error);
    exit(null);
  }
}

// ---------------------------------------------------------------------------
// Browser stand-ins for the WebView-only functions of kernelsu
// ---------------------------------------------------------------------------

/** Enters or leaves browser full screen. Resolves quietly when the browser refuses. */
export async function fullScreen(enable) {
  try {
    if (enable) await document.documentElement.requestFullscreen();
    else if (document.fullscreenElement) await document.exitFullscreen();
  } catch {
    /* needs a user gesture, or the browser does not allow it */
  }
}

const TOAST_MS = 2500;
let toastTimer;

/** Shows a short message at the bottom of the page (a plain DOM element, no CSS needed). */
export function toast(message) {
  let el = document.getElementById("axbridge-toast");
  if (!el) {
    el = document.createElement("div");
    el.id = "axbridge-toast";
    el.setAttribute("role", "status");
    el.style.cssText =
      "position:fixed;left:50%;bottom:24px;transform:translateX(-50%);max-width:80vw;" +
      "padding:10px 16px;border-radius:12px;background:#2e3135;color:#f0f0f4;" +
      "font:14px/1.4 system-ui,sans-serif;z-index:2147483647;pointer-events:none;";
    document.body.append(el);
  }
  el.textContent = String(message);
  el.hidden = false;
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => (el.hidden = true), TOAST_MS);
}
