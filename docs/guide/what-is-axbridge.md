# What is AxBridge?

AxManager exposes `ksu` and `Axeron` JavaScript objects, but **only inside a plugin's own WebUI window**.
A page opened in Chrome, or your own site on Vercel or GitHub Pages, cannot see them, so it cannot run any command.

AxBridge is an AxManager plugin that fills that gap. It starts a tiny HTTP server on the device,
`127.0.0.1:8787`, and any page opened in the phone's browser can ask it to run shell commands:

```
 Chrome / your site (https)            AxBridge (on the phone)               Android shell
 ───────────────────────────  POST   ───────────────────────────  sh -c   ──────────────
 exec("getprop ro.build...")  ─────▶  127.0.0.1:8787/exec        ───────▶  runs command
                              ◀─────  { errno, stdout, stderr }  ◀───────
```

Commands run with the same privileges as the AxManager server: **shell**, or **root** when AxManager runs in root mode.

## Why not just open the port?

An open shell over HTTP would let any website you visit run commands on your phone. AxBridge adds the layers that
make it safe to leave running: a token, an allowlist of websites, and request checks.
Read [Security model](./security) before you build on it.

## What you get

| | |
| --- | --- |
| **Plugin** | Installs from a ZIP in AxManager. Starts the server at boot, and has a WebUI to restart it, manage allowed websites and copy the token. |
| **Client** | [`axbridge-api`](https://www.npmjs.com/package/axbridge-api) on npm. `exec`, `spawn`, `fullScreen`, `toast`, same shape as `kernelsu`. |
| **HTTP API** | Two endpoints, `/exec` and `/spawn`, if you would rather call it from another language. |

## Requirements

- Android with **AxManager** (arm64 or arm).
- A browser **on the same phone**. Chrome or another Chromium browser is recommended. Safari cannot call
  `http://127.0.0.1` from an https page.
- The plugin declares `axeronPlugin=19100`. Lower it in `module.prop` if your AxManager server is older.

::: info DIFFERENCE WITH KERNELSU
`kernelsu` talks to a native WebView bridge, so it only works inside the manager app.
AxBridge talks HTTP, so it works in any browser, but you need the plugin running and a token.
Functions that need the WebView itself (`listPackages`, `getPackagesInfo`, `exit`, `moduleInfo`,
`enableEdgeToEdge`) are not available. See [Migrate from kernelsu](../api/migrate-from-kernelsu).
:::

## Credits

Based on [RiProG-ADB-Panel](https://github.com/RiProG-id) (MIT, © 2025 RiProG-id).
AxBridge adds authentication, an origin allowlist, CORS and Private Network Access handling, limits, streaming
`spawn`, and the AxManager plugin around it.
