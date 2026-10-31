import { $ } from "./dom.js";
import { moduleVersion } from "./ctl.js";

/** Shows the small top-bar title once the large title has scrolled out of view. */
export function initLargeTitle() {
  if (typeof IntersectionObserver === "undefined") return;

  new IntersectionObserver(([entry]) => {
    $("top").classList.toggle("top--scrolled", !entry.isIntersecting);
  }).observe($("largeTitle"));
}

/** Footer line with the plugin version (read from module.prop through AxManager). */
export function initFooter() {
  if (!moduleVersion) return;
  $("footer").textContent = `AxBridge ${moduleVersion} · based on RiProG-ADB-Panel (MIT)`;
}
