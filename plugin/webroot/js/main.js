/**
 * AxBridge WebUI entry point. Runs inside AxManager's plugin window (needs `window.ksu`).
 *
 *   ctl.js          talks to ctl.sh through ksu.exec
 *   status.js       server status card
 *   server.js       restart / open in browser
 *   origins.js      allowed-websites list and form
 *   token.js        token row
 *   diagnostics.js  "why will it not start" report
 *   layout.js       large title + footer
 *   icons.js, dom.js  small shared helpers
 */
import { mountIcons } from "./icons.js";
import { initLargeTitle, initFooter } from "./layout.js";
import { refreshStatus } from "./status.js";
import { initToken, refreshToken } from "./token.js";
import { initOrigins, refreshOrigins } from "./origins.js";
import { initDiagnostics } from "./diagnostics.js";
import { initServer, restartServer } from "./server.js";

mountIcons();
initLargeTitle();
initFooter();
initToken();
initDiagnostics();
initServer();
// Origins only take effect after a restart, so saving the list triggers one.
initOrigins({ onChanged: (note) => restartServer(note) });

refreshStatus();
refreshToken();
refreshOrigins();
