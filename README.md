# AxBridge

**Run [AxManager](https://axmanager.app) shell commands from a normal browser, or from your own hosted website.**

[Bahasa Indonesia](README.id.md)

AxManager's `ksu` / `Axeron` JavaScript objects only exist inside a plugin's own WebUI window. AxBridge is an
AxManager plugin that starts a tiny HTTP server on the device (`127.0.0.1:8787`), so any web page opened in the
phone's browser can run shell commands through it, with the same privileges as the AxManager server
(shell, or root in root mode).

```
 Chrome / your site (https)            AxBridge (on the phone)               Android shell
 ───────────────────────────  POST   ───────────────────────────  sh -c   ──────────────
 exec("getprop ro.build...")  ─────▶  127.0.0.1:8787/exec        ───────▶  runs command
                              ◀─────  { errno, stdout, stderr }  ◀───────
```

> Based on [RiProG-ADB-Panel](https://github.com/RiProG-id) (MIT). AxBridge adds authentication, an origin
> allowlist, CORS/Private-Network-Access handling, limits, and an AxManager plugin around it.

## Security model: please read

This is a shell reachable over HTTP. It is only safe because of these layers, so don't remove them:

| Layer | What it blocks |
| --- | --- |
| Binds to `127.0.0.1` only | Other devices on the Wi-Fi |
| Bearer **token** (random, 192-bit, per device) | Any site that does not know the token |
| **Origin allowlist** (managed by the device owner in the WebUI) | Arbitrary websites calling the bridge from the user's browser |
| `Content-Type: application/json` required | "Simple request" CSRF with `text/plain` |
| **Host** header check | DNS-rebinding |
| 30 s default timeout (max 120 s), 64 KB request, 1 MB output | Runaway commands |

Things to keep in mind as an app developer:

- The token lives in your page's `localStorage`. Any third-party script on that page can read it and run
  commands on the user's phone. Keep pages that use AxBridge free of untrusted scripts.
- Anyone holding the token (and an allowed origin) has shell access. Treat it like a password.
- If you publish a site built on this, consider exposing **specific actions** from your own code instead of passing
  free-form commands from user input.

## Install

1. Download `axbridge-<version>.zip` from [Releases](../../releases) (or build it, see below).
2. In AxManager: **Plugin → install from ZIP**.
3. Open the plugin's **WebUI**. Tap **Restart** if the status chip says *Stopped*.
4. Tap **Open in browser** to try the built-in test page.

Requirements: Android with AxManager (arm64 or arm). The plugin declares `axeronPlugin=19100`; lower it in
`module.prop` if your AxManager server is older.

## Use it in your website

1. In the plugin WebUI → **Allowed websites**, add your site's origin, e.g. `https://myapp.vercel.app`
   (use `http://localhost:5173` for local dev on the same phone).
2. Install the client: `npm i axbridge-api` (or copy [`client/axbridge.js`](client/axbridge.js) and optionally [`axbridge.d.ts`](client/axbridge.d.ts) into your project).
3. In the WebUI → **Token → Copy**, paste it into your site once (it is saved in `localStorage`).

```js
import { exec, spawn, ping, setToken } from "axbridge-api";

if (!(await ping())) alert("AxBridge is not running");
setToken("<token from the plugin WebUI>");

// one-shot command
const { errno, stdout, stderr } = await exec("getprop ro.product.model");
console.log(stdout.trim());

// long-running command, output arrives while it runs
const child = spawn("logcat", ["-v", "brief"]);
child.stdout.on("data", (text) => console.log(text));
child.on("exit", (code) => console.log("exited", code));
// child.kill();   // stops it on the device too
```

The client has the same shape as the [`kernelsu`](https://www.npmjs.com/package/kernelsu) package
(`exec`, `spawn`, `fullScreen`, `toast`), so code written for it keeps working. In a Vite project you can
point the import at AxBridge with an alias: `resolve: { alias: { kernelsu: "axbridge-api" } }`.

The page must be opened **in the browser on the same device** (`127.0.0.1` is the device itself).
Chrome/Chromium on Android is recommended; Chrome may show a "local network access" permission prompt the first time.
Safari does not allow https pages to call `http://127.0.0.1`.

### Client API

| Function | |
| --- | --- |
| `exec(cmd, { cwd?, env?, timeout? })` | Runs `cmd` with `sh -c`. Returns `{ errno, stdout, stderr }` (same shape as `kernelsu`'s `exec`). `timeout` in seconds (default 30, max 120). |
| `spawn(cmd, args?, { cwd?, env?, timeout? })` | Starts a process and streams its output. Returns a `ChildProcess` with `stdout.on("data")`, `stderr.on("data")`, `on("exit", code)`, `on("error", err)` and `kill()`. With a non-empty `args` array the command runs directly (no shell); otherwise `cmd` is a shell line. `timeout` default 300 s, max 3600 s. `code` is `null` when killed or disconnected. |
| `fullScreen(enable)` | Browser full screen (needs a user gesture). Does not throw if refused. |
| `toast(message)` | Small message at the bottom of the page. |
| `ping()` | `true` if the server answers. |
| `setToken(t)` / `getToken()` / `clearToken()` | Token storage (`localStorage`, with an in-memory fallback). |
| `setBaseUrl(url)` | Override `http://127.0.0.1:8787`. |
| `BridgeError` | Thrown by `exec` and emitted by `spawn`. `error.code` is `network`, `unauthorized`, `forbidden` or `http`. |

Not included: `listPackages`, `getPackagesInfo`, `exit`, `moduleInfo`, `enableEdgeToEdge` (they need the AxManager WebView).
Use `exec("pm list packages -3")` for package lists.

### HTTP API

```
GET  /health                      -> {"ok":true}                     (no auth)

POST /exec                        Authorization: Bearer <token>
     Content-Type: application/json
     {"cmd":"id","cwd":"/","env":{"A":"1"},"timeout":30}
                                  -> {"errno":0,"stdout":"…","stderr":""}

POST /spawn                       same headers; extra field "args":["-l"] (runs cmd directly, no shell)
                                  -> application/x-ndjson, one JSON object per line:
     {"t":"out","b":"<base64>"}   stdout chunk
     {"t":"err","b":"<base64>"}   stderr chunk
     {"t":"exit","code":0}        always the last line
```
`errno` / `code` is the exit code; `124` means the command timed out. Closing the connection kills the process
(and everything it started). `timeout` defaults: `/exec` 30 s (max 120), `/spawn` 300 s (max 3600).
Chunks are base64 so characters split across reads survive; the client decodes them for you.

## Files on the device

| Path | |
| --- | --- |
| `$AXERONDIR/axbridge/token` | The token (mode 600, created on first start) |
| `$AXERONDIR/axbridge/origins` | Allowed origins, one per line (edited from the WebUI) |
| `$AXERONDIR/logs/axbridge.log` | Server and start-up log |

## Repository layout

```
server/    Go HTTP server (+ embedded test page) .......... the core
client/    axbridge.js / axbridge.d.ts ................... browser client, published to npm as axbridge-api
plugin/    AxManager plugin: module.prop, scripts, WebUI (no binaries)
scripts/   package.sh (builds the flashable ZIP)
```

The plugin's WebUI is plain HTML, CSS and ES modules, no build step:

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
    ├── status.js · server.js · origins.js · token.js · diagnostics.js   one file per screen section
    └── layout.js · icons.js · dom.js                                    small shared helpers
```

## Build from source

Needs Go 1.24+ and `zip`.

```sh
make test      # go vet + unit tests
make zip       # cross-compiles arm64 + arm, writes dist/axbridge-<version>.zip
make host      # native binary in dist/ for local experiments
```

Local run (without a phone):

```sh
make host
./dist/axbridge -origins http://localhost:5173      # prints where the token file is
curl -s http://127.0.0.1:8787/health
```

Flags: `-addr` (default `127.0.0.1:8787`, keep it on loopback), `-origins` (comma-separated, or env
`AXBRIDGE_ORIGINS`), `-token-file`, `-pid-file`. Env `AXBRIDGE_TOKEN` sets the token explicitly.

## Troubleshooting

- **WebUI shows *Stopped*:** tap **Restart**, then **Run diagnostics** and read the output. It shows the ABI,
  file permissions, an exec test and the log tail. If the binary cannot be executed from the plugin folder, the
  start script automatically retries from `/data/local/tmp`.
- **`Failed to fetch` from your site:** the server is off, or your origin is not in *Allowed websites* (must match
  exactly, including `https://` and no trailing slash), or the browser blocked local-network access.
- **`wrong or missing token`:** copy the token again from the WebUI.

## Credits & license

MIT. Original server concept: [RiProG-ADB-Panel](https://github.com/RiProG-id) by RiProG-id. See [LICENSE](LICENSE).

## Publish the npm package

The client in `client/` is published as `axbridge-api`.

```sh
cd client
npm login
npm publish --access public      # bump "version" in client/package.json first
```

Or let GitHub do it: add an `NPM_TOKEN` repository secret, then push a tag that matches the version, e.g.
`git tag npm-v0.1.0 && git push origin npm-v0.1.0` ([workflow](.github/workflows/publish-npm.yml)).
