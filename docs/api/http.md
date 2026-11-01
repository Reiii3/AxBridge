# HTTP protocol

You do not need the JavaScript client. Any language that can send HTTP can use AxBridge.

All requests except `/health` need:

```
Authorization: Bearer <token>
Content-Type: application/json
```

Browsers also send an `Origin` header, which must be in the allowlist. Non-browser clients (curl, scripts) send
none and are checked by the token only.

## GET /health

No authentication. Returns `{"ok":true}`. Use it to see whether the server is up.

## POST /exec

Runs a command and returns when it ends.

```http
POST /exec
Authorization: Bearer <token>
Content-Type: application/json

{"cmd":"echo hi; echo oops >&2; exit 3","cwd":"/","env":{"A":"1"},"timeout":30}
```

```json
{"errno":3,"stdout":"hi\n","stderr":"oops\n"}
```

| Field | | |
| --- | --- | --- |
| `cmd` | string, **required** | Run with `sh -c` |
| `cwd` | string | Working directory |
| `env` | object | Extra environment variables. Names: `[A-Za-z_][A-Za-z0-9_]*` |
| `timeout` | number | Seconds. Default 30, max 120 |

`errno` is the exit code, `124` for a timeout, `-1` if the process could not be started.
Output is capped at 1 MB per stream.

## POST /spawn

Same fields as `/exec`, plus `args`, and the response **streams**.

```http
POST /spawn
Authorization: Bearer <token>
Content-Type: application/json

{"cmd":"echo","args":["a; echo injected","$HOME"]}
```

| Field | | |
| --- | --- | --- |
| `args` | string array | When present, `cmd` is executed **directly**, without a shell |
| `timeout` | number | Seconds. Default 300, max 3600 |

The response is `application/x-ndjson`, one JSON object per line:

```
{"t":"out","b":"YTsgZWNobyBpbmplY3RlZCAkSE9NRQo="}
{"t":"exit","code":0}
```

| Line | Meaning |
| --- | --- |
| `{"t":"out","b":"<base64>"}` | A chunk of stdout |
| `{"t":"err","b":"<base64>"}` | A chunk of stderr |
| `{"t":"exit","code":N}` | Always the last line |

Chunks are base64 so that a multi-byte character split across two reads survives intact. Decode every chunk to
bytes first, and decode text from the byte stream, not chunk by chunk.

If the process cannot be started you get a single `exit` line with code `-1`.
If the timeout is reached the process group is killed, a line with `[timeout after …]` is sent on `err`, and
`exit` has code `124`.

**Closing the connection kills the process** and everything it started. That is how `kill()` works.

### Try it with curl

```sh
TOKEN=$(cat "$AXERONDIR/axbridge/token")

curl -N -X POST http://127.0.0.1:8787/spawn \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"cmd":"sh","args":["-c","echo one; sleep 1; echo two"]}' \
  | jq -r 'select(.t=="out") | .b | @base64d'
```

## Status codes

| Code | Reason |
| --- | --- |
| 200 | Request accepted (a failing command is still 200) |
| 400 | Invalid JSON, empty `cmd`, or invalid `env` name |
| 401 | Missing or wrong token |
| 403 | Origin not in the allowlist, or a bad `Host` header |
| 405 | Wrong method |
| 415 | `Content-Type` is not `application/json` |

## CORS

For allowed origins the server answers preflight requests and sends
`Access-Control-Allow-Origin` (the exact origin), `Access-Control-Allow-Headers: Authorization, Content-Type`,
`Access-Control-Allow-Methods` and `Access-Control-Allow-Private-Network: true`, which Chrome needs when a public
https page calls a local address.
