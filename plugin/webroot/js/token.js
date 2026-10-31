import { $, copyText, showSnackbar } from "./dom.js";
import { ctl } from "./ctl.js";
import { svg } from "./icons.js";

const MASK = "•".repeat(16);

let token = "";
let revealed = false;

/** Reads the token file through ctl.sh and re-renders the row. */
export async function refreshToken() {
  token = (await ctl("token")).out;
  render();
}

function render() {
  $("token").textContent = !token ? "(not created yet)" : revealed ? token : MASK;
  $("toggleToken").innerHTML = svg(revealed ? "eyeOff" : "eye");
  $("toggleToken").setAttribute("aria-label", revealed ? "Hide token" : "Show token");
}

function copyToken() {
  if (!token) return showSnackbar("No token yet");

  copyText(token).then(
    () => showSnackbar("Token copied"),
    () => showSnackbar("Copy blocked, tap the eye and copy manually"),
  );
}

export function initToken() {
  $("toggleToken").onclick = () => {
    revealed = !revealed;
    render();
  };
  $("copyToken").onclick = copyToken;
  render();
}
