# Build dari source

Kamu butuh **Go 1.24 atau lebih baru** dan `zip`. Untuk client dan docs ini, **Node 18 atau lebih baru**.

```sh
make test      # go vet + unit test
make zip       # cross-compile arm64 + arm, menulis dist/axbridge-<versi>.zip
make host      # binary native di dist/ untuk percobaan lokal
```

`make zip` mengambil versi dari `plugin/module.prop`. Binary tidak di-commit, ia dibangun ke dalam ZIP.

## Jalankan server di komputermu

```sh
make host
./dist/axbridge -origins http://localhost:5173      # mencetak lokasi file token
curl -s http://127.0.0.1:8787/health
```

| Flag / env | |
| --- | --- |
| `-addr` | Alamat listen. Default `127.0.0.1:8787`. **Biarkan di loopback.** |
| `-origins`, `AXBRIDGE_ORIGINS` | Origin yang diizinkan, dipisah koma |
| `-token-file` | Lokasi penyimpanan token. Default `token` di folder saat ini |
| `-pid-file` | Tulis process id ke sini (dipakai plugin) |
| `AXBRIDGE_TOKEN` | Tentukan token sendiri alih-alih dibuat otomatis |

## Tes

`make test` menjalankan tes Go: matriks keamanan (token, origin, content type, host), header preflight dan Private
Network Access, `exec` dan timeout, serta aliran `spawn` (output, exit code, `args` tanpa shell, `env`, kill
seluruh process group saat timeout, UTF-8 terbelah).

## Build docs ini

```sh
cd docs
npm install
npm run dev        # http://localhost:5173/AxBridge/
npm run build      # hasil di docs/.vitepress/dist
```

Setel `DOCS_BASE=/` kalau situs disajikan dari root sebuah domain.
