import { $ } from "./dom.js";
import { ctl } from "./ctl.js";

const VIEWS = {
  running: { title: "Running", subtitle: "127.0.0.1:8787 · reachable from this device only" },
  stopped: { title: "Stopped", subtitle: "Tap Restart server to start it." },
  error: { title: "Script error", subtitle: "The plugin script could not run. Open diagnostics." },
};

/** Asks ctl.sh whether the server process is alive and updates the status card. */
export async function refreshStatus() {
  const { errno, out } = await ctl("status");
  const state = out === "running" ? "running" : errno ? "error" : "stopped";

  $("status").dataset.state = state;
  $("statusTitle").textContent = VIEWS[state].title;
  $("statusSubtitle").textContent = VIEWS[state].subtitle;
}
