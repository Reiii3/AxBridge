# Use it in your website

## 1. Allow your site

In the plugin WebUI, go to **Allowed websites** and add your site's origin, for example
`https://myapp.vercel.app`. For local development on the same phone use `http://localhost:5173`.
Details in [Allowed websites](./allowed-websites).

## 2. Install the client

::: code-group

```sh [npm]
npm i axbridge-api
```

```sh [pnpm]
pnpm add axbridge-api
```

```sh [yarn]
yarn add axbridge-api
```

:::

No bundler? Copy [`axbridge.js`](https://github.com/Reiii3/AxBridge/blob/main/client/axbridge.js)
into your project and import it by path.

## 3. Give your site the token

In the WebUI, open **Token** and tap the copy button. Paste it into your site once. The client saves it in
`localStorage`.

```js
import { setToken } from "axbridge-api";

setToken("<token from the AxBridge WebUI>");
```

## 4. Run commands

```js
import { exec, spawn, ping } from "axbridge-api";

if (!(await ping())) {
  alert("AxBridge is not running. Open the plugin and tap Restart server.");
}

// One-shot command
const { errno, stdout, stderr } = await exec("getprop ro.product.model");
console.log(stdout.trim());

// Long-running command, output arrives while it runs
const child = spawn("logcat", ["-v", "brief"]);
child.stdout.on("data", (text) => console.log(text));
child.on("exit", (code) => console.log("exited", code));

// later
child.kill();
```

`errno` is the command's exit code. A command that fails is **not** an exception. Exceptions
(`BridgeError`) mean the bridge itself could not be reached or refused you.
See [Client API](../api/client#bridgeerror).

## Good to know

- The page must be opened **in the browser on the same phone**. `127.0.0.1` is the phone itself.
- Chrome may ask for "local network access" the first time. Allow it.
- Keep pages that use AxBridge free of untrusted third-party scripts: any script on the page can read the token.
- Prefer exposing **specific actions** from your own code instead of passing user input to `exec`.
