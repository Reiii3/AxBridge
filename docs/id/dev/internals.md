# Isi plugin

## Struktur repositori

```
server/    Server HTTP Go (+ halaman tes tertanam) ........ inti
client/    axbridge.js / axbridge.d.ts ................... client browser, dipublish sebagai axbridge-api
plugin/    Plugin AxManager: module.prop, skrip, WebUI (tanpa binary)
scripts/   package.sh (membuat ZIP yang bisa di-flash)
docs/      situs ini (VitePress)
```

## Skrip plugin

| File | Kapan jalan | Fungsinya |
| --- | --- | --- |
| `module.prop` | | `id=axbridge`, versi, `axeronPlugin=19100` |
| `customize.sh` | instalasi | Menyimpan binary sesuai CPU perangkat, menghapus yang lain, mengatur izin |
| `service.sh` | boot, dan saat restart | Menghentikan server lama, menyusun daftar origin, menjalankan binary, memastikan ia naik |
| `action.sh` | tombol Action | Menjalankan `service.sh`, lalu membuka halaman tes dengan token |
| `ctl.sh` | WebUI | Satu-satunya pintu WebUI ke shell (lihat di bawah) |
| `uninstall.sh` | pencopotan | Menghentikan server, menghapus folder datanya |

`service.sh` menguji apakah binary bisa dieksekusi dari folder plugin. Kalau shell mengembalikan 126 atau 127, binary
disalin ke `/data/local/tmp` (ubah dengan `AXBRIDGE_TMP`) dan dijalankan dari sana. Server dijalankan dengan `setsid`
jika tersedia, pid file diverifikasi, dan log ditulis ke `$AXERONDIR/logs/axbridge.log`.

## ctl.sh

WebUI tidak pernah menyusun string shell. Ia memanggil `ctl.sh <perintah>` dan membaca outputnya.

| Perintah | |
| --- | --- |
| `token` | Mencetak token |
| `status` | `running` atau `stopped` |
| `restart` | Restart server, mencetak ekor log kalau gagal |
| `diag` | ABI, izin, uji eksekusi, ekor log |
| `open` | Membuka halaman tes dengan `#token=…` |
| `origins-list` | Satu origin per baris |
| `origins-add <origin>` | Memvalidasi dan menambahkan (huruf kecil) |
| `origins-del <origin>` | Menghapus satu |

`origins-add` menerima `https://host[:port]`, atau `http://localhost` / `http://127.0.0.1` dengan port opsional.
Selain itu ditolak, jadi file tidak bisa berisi sesuatu yang akan ditolak server.

## File di perangkat

| Path | |
| --- | --- |
| `$AXERONDIR/axbridge/token` | Token (mode 600, dibuat saat pertama menyala) |
| `$AXERONDIR/axbridge/origins` | Origin yang diizinkan, satu per baris |
| `$AXERONDIR/axbridge/pid` | Pid server |
| `$AXERONDIR/logs/axbridge.log` | Log server dan startup |

## WebUI

HTML, CSS, dan ES module polos. Tanpa build step, tanpa framework.

```
plugin/webroot/
├── index.html          markup saja
├── css/
│   ├── tokens.css      warna, radius, font (terang + gelap)
│   ├── base.css        reset, layout halaman, top bar, judul
│   └── components.css  kartu status, row, field form, tombol, snackbar
└── js/
    ├── main.js         entry point, menyambungkan semua modul
    ├── ctl.js          menjalankan ctl.sh lewat ksu.exec
    ├── status.js · server.js · origins.js · token.js · diagnostics.js
    └── layout.js · icons.js · dom.js
```

`ctl.js` memakai bentuk callback `ksu.exec`, jadi halaman tidak pernah terblokir selama perintah berjalan.
AxManager mengembalikan `moduleDir` sebagai string di sebagian versi dan sebagai `{ path }` di versi lain, dan
`ctl.js` menerima keduanya.

## Server

Satu file Go yang hanya memakai standard library. `index.html` (halaman tes) ditanam dengan `go:embed`.
Penanganan request: `readRequest` (method, content type, token, body, nama env), `buildCmd` (shell atau langsung),
lalu `doExec` atau `doSpawn`.
