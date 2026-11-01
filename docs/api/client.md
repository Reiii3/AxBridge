# Client API

```sh
npm i axbridge-api
```

```js
import { exec, spawn, fullScreen, toast, ping, setToken } from "axbridge-api";
```

The package is ESM only, has no dependencies and ships its own TypeScript types.

## exec

Runs a command with `sh -c` and waits for it to finish.

```ts
exec(command: string, options?: {
  cwd?: string;
  env?: Record<string, string>;
  timeout?: number; // seconds, default 30, max 120
}): Promise<{ errno: number; stdout: string; stderr: string }>
```

```js
const { errno, stdout, stderr } = await exec("ls -l", { cwd: "/data/local/tmp" });
if (errno === 0) console.log(stdout);
```

- `errno` is the exit code. `124` means the command timed out.
- A command that fails still **resolves**. Only a bridge problem rejects (see [BridgeError](#bridgeerror)).
- `env` sets environment variables for this command. Names must match `[A-Za-z_][A-Za-z0-9_]*`.

## spawn

Starts a process and streams its output.

```ts
spawn(command: string, args?: string[], options?: {
  cwd?: string;
  env?: Record<string, string>;
  timeout?: number; // seconds, default 300, max 3600
}): ChildProcess

spawn(command: string, options?: SpawnOptions): ChildProcess
```

```js
const ls = spawn("ls", ["-lh", "/data"]);

ls.stdout.on("data", (text) => console.log(`stdout: ${text}`));
ls.stderr.on("data", (text) => console.log(`stderr: ${text}`));
ls.on("exit", (code) => console.log(`exited with ${code}`));
ls.on("error", (error) => console.error(error.code));
```

### Shell line or direct

| Call | Runs as |
| --- | --- |
| `spawn("logcat -d \| grep Zygote")` | A shell line (`sh -c`), pipes and redirects work |
| `spawn("logcat", ["-d"])` | `logcat` **directly**, no shell. Arguments are never parsed, so no quoting or injection issues |

Use the `args` form whenever part of the command comes from user input.

### ChildProcess

| Member | |
| --- | --- |
| `stdout.on("data", (text: string) => …)` | Chunks of stdout, already decoded as UTF-8. A character split across two chunks is rejoined correctly. |
| `stderr.on("data", (text: string) => …)` | Same for stderr. |
| `on("exit", (code: number \| null) => …)` | Fires once, always last. `code` is `null` when you called `kill()` or the connection was cut. A server timeout gives `124`. |
| `on("error", (error) => …)` | The bridge failed (see [BridgeError](#bridgeerror)). `exit` follows with `null`. |
| `kill()` | Stops the process on the phone, **including everything it started**. |

::: tip Output arrives while the process runs
You do not have to wait for `exit`. A progress bar, a log viewer or a live `logcat` works as you would expect.
:::

Closing the tab, or any other way the connection drops, also kills the process.

## fullScreen

```ts
fullScreen(enable: boolean): Promise<void>
```

Enters or leaves browser full screen. Browsers require a user gesture (a click or tap). If the browser refuses, it
resolves quietly without throwing.

## toast

```ts
toast(message: string): void
```

Shows a short message at the bottom of the page. It is a plain DOM element with inline styles, so it needs no CSS.
Browser only.

## Connection helpers

| Function | |
| --- | --- |
| `ping(): Promise<boolean>` | `true` if the server answers. Use it to show "AxBridge is not running". |
| `setToken(token: string)` | Saves the token. |
| `getToken(): string` | Reads it back (empty string if none). |
| `clearToken()` | Forgets it. Call this when you get `unauthorized` and want to ask the user again. |
| `setBaseUrl(url: string)` | Use another address than `http://127.0.0.1:8787`. |

The token is stored in `localStorage`. If storage is blocked (private mode, sandboxed iframe) or does not exist
(Node), the client keeps it in memory for the life of the page.

## BridgeError

Thrown by `exec` and emitted as the `error` event of `spawn` when the **bridge** could not do the job.

```ts
class BridgeError extends Error {
  code: "network" | "unauthorized" | "forbidden" | "http";
}
```

| `code` | Meaning | Usual fix |
| --- | --- | --- |
| `network` | Server not reachable | Start the plugin, or the browser blocked local-network access |
| `unauthorized` | Wrong or missing token | `setToken(...)` with the token from the WebUI |
| `forbidden` | Your site is not in **Allowed websites** | Add it in the plugin WebUI |
| `http` | Any other HTTP error | The message has the status and body |

```js
try {
  await exec("id");
} catch (error) {
  if (error.code === "unauthorized") {
    clearToken();
    // ask the user to paste the token again
  }
}
```

## TypeScript

```ts
import { exec, spawn, type ExecResult, BridgeError } from "axbridge-api";

const result: ExecResult = await exec("id", { env: { A: "1" }, timeout: 5 });
```

## Node

Works in Node 18 or newer (it uses the global `fetch`), for scripts that run **on a computer that can reach the
phone's port**, for example through `adb forward tcp:8787 tcp:8787`. `fullScreen` and `toast` need a browser.
