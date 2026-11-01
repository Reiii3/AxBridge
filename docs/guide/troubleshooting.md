# Troubleshooting

## The WebUI says *Stopped*

1. Tap **Restart server**.
2. Still stopped? Tap **Run diagnostics** and read the output. It shows the CPU ABI, file permissions, a test run of
   the binary, and the tail of the log.
3. If the binary cannot be executed from the plugin folder, the start script retries automatically from
   `/data/local/tmp`.
4. Use **Copy** on the diagnostics box and include it if you open an issue.

## `Failed to fetch` from your site

Check these in order:

1. The server is running (WebUI shows *Running*).
2. Your origin is in **Allowed websites**, **exactly**: include `https://`, no trailing slash, no path.
3. You opened the page in a browser **on the same phone**.
4. Chrome did not block "local network access". Open the site settings for your page and allow it.
5. You are not on Safari or on an in-app browser that blocks calls to `127.0.0.1` from https pages.

## `BridgeError: unauthorized`

The token is wrong or missing. Copy it again from the WebUI (**Token** → copy) and call `setToken(...)`.
Restarting the server does **not** change the token.

## `BridgeError: forbidden`

The site is not in the allowlist. Add it, wait for the server to restart, and try again.

## A command returns `errno: 124`

It hit the timeout. `exec` allows at most 120 s. For long jobs use `spawn` (up to 3600 s).

## `spawn` does not exist or returns 404

Your plugin is older than v0.1.6. Install a newer ZIP, because `/spawn` is a server feature.

## Output looks cut off

`exec` keeps at most 1 MB of output. Use `spawn` to stream large output.

## The server stops after a while

Some ROMs kill background processes. Disable battery optimisation for AxManager and keep the plugin enabled.
The server restarts at boot, and **Restart server** brings it back at any time.

## Reading the log

```sh
cat "$AXERONDIR/logs/axbridge.log"
```
