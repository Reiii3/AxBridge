// AxBridge browser client (ES module, no dependencies).
// Same result shape as kernelsu's exec(): { errno, stdout, stderr }.
let BASE = "http://127.0.0.1:8787";
const KEY = "axbridge_token";

/** Point the client at another address (default http://127.0.0.1:8787). */
export const setBaseUrl = (url) => { BASE = url.replace(/\/+$/, ""); };

export const getToken = () => { try { return localStorage.getItem(KEY) || ""; } catch { return ""; } };
export const setToken = (t) => { try { localStorage.setItem(KEY, t.trim()); } catch {} };

/** true when the AxBridge server answers on this device. */
export async function ping() {
  try { return (await fetch(`${BASE}/health`)).ok; } catch { return false; }
}

/**
 * Run a shell command on the device.
 * @param {string} cmd
 * @param {{cwd?: string, timeout?: number}} [opts] timeout in seconds (default 30, max 120)
 * @returns {Promise<{errno:number, stdout:string, stderr:string}>}
 */
export async function exec(cmd, { cwd, timeout } = {}) {
  const res = await fetch(`${BASE}/exec`, {
    method: "POST",
    headers: { "Content-Type": "application/json", Authorization: `Bearer ${getToken()}` },
    body: JSON.stringify({ cmd, cwd, timeout }),
  });
  if (res.status === 401) throw new Error("AxBridge: wrong or missing token");
  if (res.status === 403) throw new Error("AxBridge: this website is not in the allowed list");
  if (!res.ok) throw new Error(`AxBridge: HTTP ${res.status}: ${await res.text()}`);
  return res.json();
}
