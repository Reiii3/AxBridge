# axbridge-api

Browser client for [AxBridge](https://github.com/Reiii3/axbridge), the AxManager plugin that lets a normal
browser (or a hosted https site) run shell commands on the same Android device.

The API has the same shape as the [`kernelsu`](https://www.npmjs.com/package/kernelsu) package, so code written for
it keeps working.

```sh
npm i axbridge-api
```

```js
import { exec, spawn, setToken, ping } from "axbridge-api";

if (!(await ping())) throw new Error("AxBridge is not running");
setToken("<token from the AxBridge WebUI>");   // saved in localStorage

const { errno, stdout } = await exec("getprop ro.product.model");

const child = spawn("logcat", ["-v", "brief"]);
child.stdout.on("data", (text) => console.log(text));
child.on("exit", (code) => console.log("exited", code));
// child.kill();
```

Use it in place of `kernelsu` without touching your imports (Vite):

```js
// vite.config.js
export default { resolve: { alias: { kernelsu: "axbridge-api" } } };
```

| Export | |
| --- | --- |
| `exec(cmd, { cwd?, env?, timeout? })` | `Promise<{ errno, stdout, stderr }>` |
| `spawn(cmd, args?, { cwd?, env?, timeout? })` | `ChildProcess` with `stdout`/`stderr` `on("data")`, `on("exit")`, `on("error")`, `kill()` |
| `fullScreen(enable)`, `toast(message)` | Browser stand-ins |
| `ping()`, `setToken()`, `getToken()`, `clearToken()`, `setBaseUrl()` | Connection helpers |
| `BridgeError` | `error.code`: `network`, `unauthorized`, `forbidden`, `http` |

This package is only the client. You also need the AxBridge plugin installed on the phone and your site added to its
**Allowed websites** list. Setup, security notes and the HTTP protocol are in the
[main README](https://github.com/Reiii3/axbridge#readme).

Requires a server with `/spawn` (plugin v0.1.6 or newer) for `spawn`.

MIT. Based on RiProG-ADB-Panel (MIT, © 2025 RiProG-id).
