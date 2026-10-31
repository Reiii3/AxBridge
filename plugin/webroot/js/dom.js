/** Small DOM helpers shared by all modules. */

export const $ = (id) => document.getElementById(id);

const SNACKBAR_MS = 2800;
let snackbarTimer;

/** Shows a short message at the bottom of the screen. */
export function showSnackbar(message) {
  const el = $("snackbar");
  el.textContent = message;
  el.classList.add("snackbar--visible");
  clearTimeout(snackbarTimer);
  snackbarTimer = setTimeout(() => el.classList.remove("snackbar--visible"), SNACKBAR_MS);
}

/** Copies text to the clipboard. Rejects when the WebView does not allow it. */
export function copyText(text) {
  if (!navigator.clipboard) return Promise.reject(new Error("clipboard unavailable"));
  return navigator.clipboard.writeText(text);
}
