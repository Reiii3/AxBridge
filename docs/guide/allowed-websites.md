# Allowed websites

A browser sends an `Origin` header with every cross-site request. AxBridge only answers origins that the device
owner has added to the allowlist, so a random website you open cannot call the bridge, even if it guessed the port.

## Add a site

In the plugin WebUI, **Allowed websites** → type the site → **Add**. The server restarts by itself so the change
takes effect.

You can paste a full URL. The path, query and hash are dropped and the host is lowercased:

| You type | Saved as |
| --- | --- |
| `New-Site.Example.com/some/path?x=1` | `https://new-site.example.com` |
| `https://app.example.com:8443/x` | `https://app.example.com:8443` |
| `http://localhost:5173/` | `http://localhost:5173` |

## Rules

- `https://host` or `https://host:port`.
- Plain `http` is accepted for `localhost` and `127.0.0.1` only, for local development.
- Exact origin only. No wildcards, no paths, no trailing slash.
- `https://app.example.com` and `https://www.app.example.com` are different origins. Add both if you use both.

## Always allowed

The built-in test page (`http://127.0.0.1:8787` and `http://localhost:8787`) can always call the server. That is
what **Open in browser** uses.

## Remove a site

Tap the trash icon on its row. The server restarts again.

## Where it is stored

One origin per line in `$AXERONDIR/axbridge/origins`. The WebUI is the supported way to edit it. You can also use
the `ctl.sh` helper from a shell, see [Plugin internals](../dev/internals#ctl-sh).
