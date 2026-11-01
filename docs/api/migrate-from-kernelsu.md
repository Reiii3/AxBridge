# Migrate from kernelsu

AxBridge's client has the same shape as the [`kernelsu`](https://www.npmjs.com/package/kernelsu) package, so a web
project that already uses `kernelsu` usually needs one line of config.

## What maps one to one

| `kernelsu` | `axbridge-api` | Notes |
| --- | --- | --- |
| `exec(cmd, { cwd, env })` | `exec(cmd, { cwd, env, timeout })` | Same result `{ errno, stdout, stderr }` |
| `spawn(cmd, args, { cwd, env })` | `spawn(cmd, args, { cwd, env, timeout })` | Same events, plus a real `kill()` |
| `fullScreen(enable)` | `fullScreen(enable)` | Browser full screen, needs a user gesture |
| `toast(message)` | `toast(message)` | A DOM toast instead of an Android toast |

## What is not available

These need the manager app's WebView, which a browser does not have.

| `kernelsu` | In a browser |
| --- | --- |
| `listPackages(type)` | `exec("pm list packages -3")` for user apps, `-s` for system apps, no flag for all. Strip the `package:` prefix. |
| `getPackagesInfo(packages)` | Parse `dumpsys package <name>` for `versionName`, `versionCode` and `userId`. Run it once for many packages inside one command. App labels and icons are not available. |
| `moduleInfo()`, `exit()`, `enableEdgeToEdge()` | No equivalent |

## Option 1: alias (no code changes)

::: code-group

```js [vite.config.js]
export default {
  resolve: { alias: { kernelsu: "axbridge-api" } },
};
```

```js [webpack.config.js]
module.exports = {
  resolve: { alias: { kernelsu: "axbridge-api" } },
};
```

:::

Every `import … from "kernelsu"` now resolves to AxBridge in that build.

## Option 2: one app, two environments

Use `kernelsu` when your page runs inside the manager's WebView and AxBridge in a normal browser:

```js
const api = window.ksu
  ? await import("kernelsu")
  : await import("axbridge-api");

const { stdout } = await api.exec("getprop ro.product.model");
```

Both modules have the same `exec` and `spawn`, so the rest of your code does not care.

## Differences to remember

- In a browser the user must allow your site and paste the token once. See [Use it in your website](../guide/quick-start).
- Handle `BridgeError` (server off, wrong token, site not allowed). `kernelsu` has no equivalent, because inside the
  WebView there is nothing to fail.
- `exec` stops after 120 s at most. Use `spawn` for long jobs.
