# Plugin internals

## Repository layout

```
server/    Go HTTP server (+ embedded test page) .......... the core
client/    axbridge.js / axbridge.d.ts ................... browser client, published as axbridge-api
plugin/    AxManager plugin: module.prop, scripts, WebUI (no binaries)
scripts/   package.sh (builds the flashable ZIP)
docs/      this site (VitePress)
```

## Plugin scripts

| File | When it runs | What it does |
| --- | --- | --- |
| `module.prop` | | `id=axbridge`, version, `axeronPlugin=19100` |
| `customize.sh` | install | Keeps the binary for the device's CPU, removes the other, sets permissions |
| `service.sh` | boot, and on restart | Stops an old server, builds the origin list, starts the binary, checks it came up |
| `action.sh` | Action button | Runs `service.sh`, then opens the test page with the token |
| `ctl.sh` | WebUI | The WebUI's only entry point into the shell (see below) |
| `uninstall.sh` | removal | Stops the server, deletes its data folder |

`service.sh` tests whether the binary can be executed from the plugin folder. If the shell returns 126 or 127 it
copies the binary to `/data/local/tmp` (override with `AXBRIDGE_TMP`) and runs it from there. It starts the server
with `setsid` when available, verifies the pid file and writes to `$AXERONDIR/logs/axbridge.log`.

## ctl.sh

The WebUI never builds shell strings. It calls `ctl.sh <command>` and reads the output.

| Command | |
| --- | --- |
| `token` | Prints the token |
| `status` | `running` or `stopped` |
| `restart` | Restarts the server, prints the log tail if it fails |
| `diag` | ABI, permissions, exec test, log tail |
| `open` | Opens the test page with `#token=…` |
| `origins-list` | One origin per line |
| `origins-add <origin>` | Validates and appends (lowercased) |
| `origins-del <origin>` | Removes one |

`origins-add` accepts `https://host[:port]`, or `http://localhost` / `http://127.0.0.1` with an optional port.
Anything else is refused, so the file cannot hold something the server would reject.

## Files on the device

| Path | |
| --- | --- |
| `$AXERONDIR/axbridge/token` | The token (mode 600, created on first start) |
| `$AXERONDIR/axbridge/origins` | Allowed origins, one per line |
| `$AXERONDIR/axbridge/pid` | Server pid |
| `$AXERONDIR/logs/axbridge.log` | Server and start-up log |

## WebUI

Plain HTML, CSS and ES modules. There is no build step, and no framework.

```
plugin/webroot/
├── index.html          markup only
├── css/
│   ├── tokens.css      colors, radii, fonts (light + dark)
│   ├── base.css        reset, page layout, top bar, titles
│   └── components.css  status card, rows, form field, buttons, snackbar
└── js/
    ├── main.js         entry point, wires the modules together
    ├── ctl.js          runs ctl.sh through ksu.exec
    ├── status.js · server.js · origins.js · token.js · diagnostics.js
    └── layout.js · icons.js · dom.js
```

`ctl.js` uses the callback form of `ksu.exec`, so the page never blocks while a command runs. AxManager returns
`moduleDir` as a string in some versions and as `{ path }` in others, and `ctl.js` accepts both.

## Server

A single Go file using only the standard library. `index.html` (the test page) is embedded with `go:embed`.
Request handling is `readRequest` (method, content type, token, body, env names), `buildCmd` (shell or direct),
then `doExec` or `doSpawn`.
