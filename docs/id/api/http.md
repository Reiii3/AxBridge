# Protokol HTTP

Kamu tidak harus memakai client JavaScript. Bahasa apa pun yang bisa mengirim HTTP bisa memakai AxBridge.

Semua request kecuali `/health` membutuhkan:

```
Authorization: Bearer <token>
Content-Type: application/json
```

Browser juga mengirim header `Origin`, yang harus ada di allowlist. Client non-browser (curl, script) tidak
mengirimnya dan hanya dicek lewat token.

## GET /health

Tanpa autentikasi. Mengembalikan `{"ok":true}`. Pakai untuk melihat apakah server menyala.

## POST /exec

Menjalankan perintah dan mengembalikan hasil saat selesai.

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
| `cmd` | string, **wajib** | Dijalankan dengan `sh -c` |
| `cwd` | string | Direktori kerja |
| `env` | object | Environment variable tambahan. Nama: `[A-Za-z_][A-Za-z0-9_]*` |
| `timeout` | number | Detik. Default 30, maks 120 |

`errno` adalah exit code, `124` untuk timeout, `-1` kalau proses tidak bisa dimulai.
Output dibatasi 1 MB per aliran.

## POST /spawn

Field sama seperti `/exec`, ditambah `args`, dan responsnya **mengalir**.

```http
POST /spawn
Authorization: Bearer <token>
Content-Type: application/json

{"cmd":"echo","args":["a; echo injected","$HOME"]}
```

| Field | | |
| --- | --- | --- |
| `args` | array string | Kalau ada, `cmd` dijalankan **langsung**, tanpa shell |
| `timeout` | number | Detik. Default 300, maks 3600 |

Responsnya `application/x-ndjson`, satu objek JSON per baris:

```
{"t":"out","b":"YTsgZWNobyBpbmplY3RlZCAkSE9NRQo="}
{"t":"exit","code":0}
```

| Baris | Arti |
| --- | --- |
| `{"t":"out","b":"<base64>"}` | Potongan stdout |
| `{"t":"err","b":"<base64>"}` | Potongan stderr |
| `{"t":"exit","code":N}` | Selalu baris terakhir |

Potongan memakai base64 supaya karakter multi-byte yang terbelah di dua pembacaan tetap utuh. Dekode setiap potongan
menjadi byte dulu, lalu dekode teks dari aliran byte, bukan potongan demi potongan.

Kalau proses tidak bisa dimulai, kamu mendapat satu baris `exit` dengan kode `-1`.
Kalau timeout tercapai, process group dimatikan, baris `[timeout after …]` dikirim lewat `err`, dan `exit` berkode `124`.

**Menutup koneksi mematikan proses** beserta semua yang dijalankannya. Begitulah `kill()` bekerja.

### Coba dengan curl

```sh
TOKEN=$(cat "$AXERONDIR/axbridge/token")

curl -N -X POST http://127.0.0.1:8787/spawn \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"cmd":"sh","args":["-c","echo one; sleep 1; echo two"]}' \
  | jq -r 'select(.t=="out") | .b | @base64d'
```

## Kode status

| Kode | Alasan |
| --- | --- |
| 200 | Request diterima (perintah yang gagal tetap 200) |
| 400 | JSON tidak valid, `cmd` kosong, atau nama `env` tidak valid |
| 401 | Token tidak ada atau salah |
| 403 | Origin tidak ada di allowlist, atau header `Host` salah |
| 405 | Method salah |
| 415 | `Content-Type` bukan `application/json` |

## CORS

Untuk origin yang diizinkan, server menjawab preflight dan mengirim
`Access-Control-Allow-Origin` (origin persis), `Access-Control-Allow-Headers: Authorization, Content-Type`,
`Access-Control-Allow-Methods`, dan `Access-Control-Allow-Private-Network: true`, yang dibutuhkan Chrome saat halaman
https publik memanggil alamat lokal.
