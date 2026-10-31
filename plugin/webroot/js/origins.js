import { $, showSnackbar } from "./dom.js";
import { ctl } from "./ctl.js";
import { svg } from "./icons.js";

/** https://host[:port], or http only for localhost / 127.0.0.1. Mirrors the check in ctl.sh. */
const ORIGIN_PATTERN = /^(https:\/\/[a-z0-9]([a-z0-9.-]*[a-z0-9])?|http:\/\/(localhost|127\.0\.0\.1))(:[0-9]{1,5})?$/;

const DEFAULT_HINT = "Paste a URL, the path is dropped automatically.";

let busy = false;
let onChanged = async () => {};

/**
 * Turns what the user typed into an origin.
 * @returns {{ origin: string } | { error: string }}
 */
export function normalizeOrigin(input) {
  const text = input.trim();
  if (!text) return { error: "Enter a website" };

  try {
    const hasScheme = /^[a-z][a-z0-9+.-]*:\/\//i.test(text);
    const url = new URL(hasScheme ? text : `https://${text}`);
    const origin = url.origin.toLowerCase();

    if (!ORIGIN_PATTERN.test(origin)) {
      return { error: url.protocol === "http:" ? "Plain http is only allowed for localhost" : "Not a valid origin" };
    }
    return { origin };
  } catch {
    return { error: "Not a valid URL" };
  }
}

function showError(message) {
  $("originField").classList.toggle("field--error", Boolean(message));
  $("originSupport").textContent = message || DEFAULT_HINT;
}

function renderList(origins) {
  const list = $("originList");
  list.textContent = "";

  if (origins.length === 0) {
    const empty = document.createElement("li");
    empty.className = "empty";
    empty.textContent = "No websites yet. The built-in test page always works.";
    list.append(empty);
    return;
  }

  for (const origin of origins) {
    list.append(createRow(origin));
  }
}

function createRow(origin) {
  const row = document.createElement("li");
  row.className = "row";
  row.innerHTML = `
    <div class="row__icon">${svg("globe")}</div>
    <div class="row__text">
      <div class="row__title"></div>
      <div class="row__summary"></div>
    </div>
    <button class="icon-button" aria-label="Remove website">${svg("trash")}</button>`;

  // textContent, never innerHTML, for anything that comes from the origins file
  row.querySelector(".row__title").textContent = origin;
  row.querySelector(".row__summary").textContent = origin.startsWith("https") ? "HTTPS" : "Local development";
  row.querySelector("button").onclick = () => remove(origin);
  return row;
}

/** Reloads the list from `$AXERONDIR/axbridge/origins`. */
export async function refreshOrigins() {
  const { out } = await ctl("origins-list");
  renderList(out ? out.split("\n").filter(Boolean) : []);
}

async function add() {
  if (busy) return;

  const result = normalizeOrigin($("originInput").value);
  if (result.error) return showError(result.error);
  showError("");

  busy = true;
  $("addOrigin").disabled = true;

  const { errno, err } = await ctl(`origins-add ${result.origin}`);
  if (errno) {
    showError(err || "Could not save");
  } else {
    $("originInput").value = "";
    await onChanged("Added.");
  }

  busy = false;
  $("addOrigin").disabled = false;
}

async function remove(origin) {
  if (busy) return;
  busy = true;

  const { errno, err } = await ctl(`origins-del ${origin}`);
  if (errno) showSnackbar(err || "Could not remove");
  else await onChanged("Removed.");

  busy = false;
}

/** @param {{ onChanged: (note: string) => Promise<void> }} options called after the list was saved */
export function initOrigins(options) {
  onChanged = options.onChanged;

  $("addOrigin").onclick = add;
  $("originInput").onkeydown = (event) => {
    if (event.key === "Enter") add();
  };
  $("originInput").oninput = () => showError("");
}
