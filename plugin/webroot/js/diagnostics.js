import { $, copyText, showSnackbar } from "./dom.js";
import { ctl, MODULE_DIR, moduleInfoRaw } from "./ctl.js";

let report = "";

/** Collects ctl.sh's diagnostics plus what the WebUI itself sees, and shows it on screen. */
export async function runDiagnostics() {
  const { errno, out, err } = await ctl("diag");

  report = [
    `webui: MODULE_DIR=${MODULE_DIR}`,
    `moduleInfo=${String(moduleInfoRaw).slice(0, 300)}`,
    `ctl errno=${errno}`,
    out,
    err ? `[stderr]\n${err}` : "",
  ].join("\n");

  $("diagnosticsOutput").textContent = report;
  $("diagnostics").hidden = false;
  $("diagnostics").scrollIntoView({ behavior: "smooth", block: "nearest" });
}

export function initDiagnostics() {
  $("diagnose").onclick = runDiagnostics;
  $("copyDiagnostics").onclick = () =>
    copyText(report).then(
      () => showSnackbar("Copied"),
      () => showSnackbar("Copy blocked, long-press the text to select"),
    );
}
