# Build from source

You need **Go 1.24 or newer** and `zip`. For the client and these docs, **Node 18 or newer**.

```sh
make test      # go vet + unit tests
make zip       # cross-compiles arm64 + arm, writes dist/axbridge-<version>.zip
make host      # native binary in dist/ for local experiments
```

`make zip` takes the version from `plugin/module.prop`. Binaries are not committed, they are built into the ZIP.

## Run the server on your computer

```sh
make host
./dist/axbridge -origins http://localhost:5173      # prints where the token file is
curl -s http://127.0.0.1:8787/health
```

| Flag / env | |
| --- | --- |
| `-addr` | Listen address. Default `127.0.0.1:8787`. **Keep it on loopback.** |
| `-origins`, `AXBRIDGE_ORIGINS` | Comma-separated allowed origins |
| `-token-file` | Where the token is stored. Default `token` in the current folder |
| `-pid-file` | Write the process id here (the plugin uses it) |
| `AXBRIDGE_TOKEN` | Set the token explicitly instead of generating one |

## Tests

`make test` runs the Go tests: the security matrix (token, origin, content type, host), preflight and Private Network
Access headers, `exec` and timeouts, and the `spawn` stream (output, exit code, no-shell `args`, `env`, timeout
kill of the whole process group, split UTF-8).

## Build these docs

```sh
cd docs
npm install
npm run dev        # http://localhost:5173/AxBridge/
npm run build      # output in docs/.vitepress/dist
```

Set `DOCS_BASE=/` if you serve the site from the root of a domain.
