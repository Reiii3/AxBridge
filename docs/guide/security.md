# Security model

This is a shell reachable over HTTP. It is only safe because of the layers below, so please do not remove them.

| Layer | What it blocks |
| --- | --- |
| Binds to `127.0.0.1` only | Other devices on the same Wi-Fi |
| Bearer **token** (random, 192-bit, per device) | Any site that does not know the token |
| **Origin allowlist**, managed by the device owner in the WebUI | Arbitrary websites calling the bridge from the user's browser |
| `Content-Type: application/json` required | "Simple request" CSRF with `text/plain` |
| **Host** header check | DNS rebinding |
| Timeouts and size limits | Runaway commands |
| Process-group kill | Leftover children after a timeout or `kill()` |

## Limits

| | `/exec` | `/spawn` |
| --- | --- | --- |
| Default timeout | 30 s | 300 s |
| Maximum timeout | 120 s | 3600 s |
| Request body | 64 KB | 64 KB |
| Output | 1 MB (extra is dropped) | 16 MB (the process is killed past that) |

## The token

- Created on first start, stored in `$AXERONDIR/axbridge/token` with mode 600.
- Compared in constant time.
- Sent as `Authorization: Bearer <token>`.
- You can set your own with the `AXBRIDGE_TOKEN` environment variable.

::: danger Treat the token like a password
Anyone who holds the token **and** runs from an allowed origin has shell access to the phone.
:::

## For app developers

- The token lives in your page's `localStorage`. **Any third-party script on that page can read it** and run
  commands on the user's phone. Keep such pages free of untrusted scripts, ad code and analytics you do not control.
- A cross-site scripting bug in your site becomes a remote shell on your users' phones. Treat XSS as critical.
- Do not forward free-form user input to `exec`. Expose **specific actions** from your own code
  (`restartGame()`, `clearCache()`) and build the command yourself.
- With `spawn`, pass an `args` array so no shell parses the arguments.

## What AxBridge does not protect against

- A malicious app already running on the same phone with the token. Apps can reach `127.0.0.1` too, and the origin
  check only constrains browsers. That is why the token exists and why it must stay private.
- A device owner who adds an untrusted website to the allowlist.

## Reporting a problem

Please open a private security advisory on the GitHub repository instead of a public issue.
