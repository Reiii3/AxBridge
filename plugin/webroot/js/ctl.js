/**
 * Bridge to the plugin's shell side.
 *
 * Every action in this WebUI ends up as `sh <plugin>/ctl.sh <command>` run through AxManager's
 * `ksu.exec`. Keeping that in one place means the rest of the code never builds shell strings.
 */

const PLUGIN_ID = "axbridge";

/** Raw JSON string from AxManager (kept for the diagnostics report). */
export const moduleInfoRaw = readModuleInfo();

/** Absolute plugin folder. Falls back to $AXERONDIR, which AxManager sets for every command. */
export const MODULE_DIR = parseModuleDir(moduleInfoRaw) ?? `$AXERONDIR/plugins/${PLUGIN_ID}`;

/** Plugin version from module.prop, or "" when unknown. */
export const moduleVersion = parseField(moduleInfoRaw, "version");

function readModuleInfo() {
  try {
    return ksu.moduleInfo();
  } catch (error) {
    return String(error);
  }
}

function parseField(raw, field) {
  try {
    return JSON.parse(raw)[field] ?? "";
  } catch {
    return "";
  }
}

/**
 * AxManager returns `moduleDir` as a string in some versions and as an object like
 * `{ path: "/data/..." }` in others. Only a plain absolute path is accepted.
 */
function parseModuleDir(raw) {
  try {
    const info = JSON.parse(raw);
    let dir = info.moduleDir !== undefined ? info.moduleDir : info.path || info.dir;
    if (dir && typeof dir === "object") {
      dir = dir.path || dir.absolutePath || dir.canonicalPath;
    }
    if (typeof dir === "string" && /^\/[A-Za-z0-9_.\/-]+$/.test(dir)) {
      return dir.replace(/\/+$/, "");
    }
  } catch {
    /* not JSON: fall through to the default */
  }
  return null;
}

/**
 * Runs `ctl.sh <args>` and resolves with `{ errno, out, err }` (output is trimmed).
 * Uses the callback form of ksu.exec so the page never blocks while the command runs.
 */
export function ctl(args) {
  return new Promise((resolve) => {
    const callback = `ctl_${Math.random().toString(36).slice(2)}`;
    window[callback] = (errno, stdout, stderr) => {
      delete window[callback];
      resolve({ errno, out: (stdout || "").trim(), err: (stderr || "").trim() });
    };
    ksu.exec(`sh ${MODULE_DIR}/ctl.sh ${args}`, "{}", callback);
  });
}
