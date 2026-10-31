import { $, showSnackbar } from "./dom.js";
import { ctl } from "./ctl.js";
import { refreshStatus } from "./status.js";
import { refreshToken } from "./token.js";
import { refreshOrigins } from "./origins.js";
import { runDiagnostics } from "./diagnostics.js";

/**
 * Stops and starts the server (service.sh), then refreshes everything that depends on it.
 * @param {string} [note] short text shown before the progress message, e.g. "Added."
 * @returns {Promise<boolean>} true when the server came back up
 */
export async function restartServer(note = "") {
  const prefix = note ? `${note} ` : "";

  showSnackbar(`${prefix}Restarting…`);
  const { errno } = await ctl("restart");
  await Promise.all([refreshStatus(), refreshOrigins(), refreshToken()]);

  if (errno) {
    showSnackbar("Server failed to start. See diagnostics.");
    await runDiagnostics();
    return false;
  }
  showSnackbar(`${prefix}Server restarted.`);
  return true;
}

async function openInBrowser() {
  const { errno } = await ctl("open");
  showSnackbar(errno ? "Could not open a browser" : "Opening browser…");
}

export function initServer() {
  $("restart").onclick = () => restartServer();
  $("openBrowser").onclick = openInBrowser;
}
