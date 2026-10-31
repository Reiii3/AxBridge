# AxBridge

**Jalankan perintah shell [AxManager](https://axmanager.app) dari browser biasa, atau dari website hosting-mu sendiri.**

[English](README.md)

Objek JavaScript `ksu` / `Axeron` milik AxManager hanya ada di jendela WebUI plugin. AxBridge adalah plugin AxManager
yang menyalakan server HTTP kecil di perangkat (`127.0.0.1:8787`), sehingga halaman web apa pun yang dibuka di
browser HP yang sama bisa menjalankan perintah shell lewat server itu, dengan hak akses yang sama dengan server
AxManager (shell, atau root di mode root).

```
 Chrome / situsmu (https)               AxBridge (di HP)                      Shell Android
 ───────────────────────────  POST   ───────────────────────────  sh -c   ──────────────
 exec("getprop ro.build...")  ─────▶  127.0.0.1:8787/exec        ───────▶  jalankan perintah
                              ◀─────  { errno, stdout, stderr }  ◀───────
```

> Berbasis [RiProG-ADB-Panel](https://github.com/RiProG-id) (MIT). AxBridge menambahkan autentikasi, allowlist origin,
> penanganan CORS/Private-Network-Access, batasan, dan plugin AxManager di sekelilingnya.

## Model keamanan: wajib dibaca

Ini adalah shell yang bisa dijangkau lewat HTTP. Ia aman hanya karena lapisan-lapisan ini, jadi jangan dihapus:

| Lapisan | Yang diblokir |
| --- | --- |
| Hanya bind ke `127.0.0.1` | Perangkat lain di WiFi yang sama |
| **Token** bearer (acak, 192-bit, per perangkat) | Situs yang tidak tahu token |
| **Allowlist origin** (diatur pemilik perangkat lewat WebUI) | Website sembarang yang memanggil bridge dari browser user |
| Wajib `Content-Type: application/json` | CSRF "simple request" dengan `text/plain` |
| Cek header **Host** | DNS-rebinding |
| Timeout 30 dtk (maks 120), request 64 KB, output 1 MB | Perintah yang lepas kendali |

Catatan untuk developer aplikasi:

- Token disimpan di `localStorage` halamanmu. Script pihak ketiga di halaman itu bisa membacanya dan menjalankan
  perintah di HP user. Jaga halaman yang memakai AxBridge tetap bebas dari script tak tepercaya.
- Siapa pun yang memegang token (dan origin yang diizinkan) punya akses shell. Perlakukan seperti password.
- Kalau kamu merilis situs berbasis ini, pertimbangkan membuka **aksi spesifik** dari kodemu sendiri, bukan
  meneruskan perintah bebas dari input user.

## Pasang

1. Unduh `axbridge-<versi>.zip` dari [Releases](../../releases) (atau build sendiri, lihat bawah).
2. Di AxManager: **Plugin → pasang dari ZIP**.
3. Buka **WebUI** plugin. Tap **Restart** kalau chip status menunjukkan *Stopped*.
4. Tap **Open in browser** untuk mencoba halaman test bawaan.

Syarat: Android dengan AxManager (arm64 atau arm). Plugin memakai `axeronPlugin=19100`; turunkan di `module.prop`
kalau server AxManager-mu lebih lama.

## Pakai di website-mu

1. Di WebUI plugin → **Allowed websites**, tambah origin situsmu, misalnya `https://myapp.vercel.app`
   (pakai `http://localhost:5173` untuk dev lokal di HP yang sama).
2. Salin [`client/axbridge.js`](client/axbridge.js) (dan opsional [`axbridge.d.ts`](client/axbridge.d.ts)) ke project-mu.
3. Di WebUI → **Token → Copy**, tempel ke situsmu sekali saja (disimpan di `localStorage`).

```js
import { exec, ping, setToken } from "./axbridge.js";

if (!(await ping())) alert("AxBridge belum jalan");
setToken("<token dari WebUI plugin>");

const { errno, stdout, stderr } = await exec("getprop ro.product.model");
console.log(stdout.trim());
```

Halaman harus dibuka **di browser pada perangkat yang sama** (`127.0.0.1` adalah perangkat itu sendiri).
Disarankan Chrome/Chromium Android; Chrome bisa menampilkan prompt izin "akses jaringan lokal" saat pertama kali.
Safari tidak mengizinkan halaman https memanggil `http://127.0.0.1`.

### API client

| Fungsi | |
| --- | --- |
| `exec(cmd, { cwd?, timeout? })` | Menjalankan `cmd` dengan `sh -c`. Mengembalikan `{ errno, stdout, stderr }` (bentuk sama dengan `exec` di `kernelsu`). `timeout` dalam detik. |
| `ping()` | `true` kalau server menjawab. |
| `setToken(t)` / `getToken()` | Penyimpanan token (`localStorage`). |
| `setBaseUrl(url)` | Ganti `http://127.0.0.1:8787`. |

### API HTTP

```
GET  /health                      -> {"ok":true}                     (tanpa auth)
POST /exec                        Authorization: Bearer <token>
     Content-Type: application/json
     {"cmd":"id","cwd":"/","timeout":30}  -> {"errno":0,"stdout":"…","stderr":""}
```
`errno` adalah exit code; `124` berarti perintah kena timeout.

## File di perangkat

| Path | |
| --- | --- |
| `$AXERONDIR/axbridge/token` | Token (mode 600, dibuat saat start pertama) |
| `$AXERONDIR/axbridge/origins` | Origin yang diizinkan, satu per baris (diubah dari WebUI) |
| `$AXERONDIR/logs/axbridge.log` | Log server dan proses start |

## Struktur repo

```
server/    Server HTTP Go (+ halaman test tertanam) ........ inti
client/    axbridge.js / axbridge.d.ts ..................... client browser
plugin/    Plugin AxManager: module.prop, script, WebUI (tanpa binary)
scripts/   package.sh (membuat ZIP yang bisa di-flash)
```

WebUI plugin berupa HTML, CSS, dan ES module biasa, tanpa build step:

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
    ├── status.js · server.js · origins.js · token.js · diagnostics.js   satu file per bagian layar
    └── layout.js · icons.js · dom.js                                    helper kecil bersama
```

## Build dari source

Butuh Go 1.24+ dan `zip`.

```sh
make test      # go vet + unit test
make zip       # cross-compile arm64 + arm, hasil dist/axbridge-<versi>.zip
make host      # binary native di dist/ untuk coba-coba di PC
```

Jalankan lokal (tanpa HP):

```sh
make host
./dist/axbridge -origins http://localhost:5173
curl -s http://127.0.0.1:8787/health
```

Flag: `-addr` (default `127.0.0.1:8787`, jangan keluar dari loopback), `-origins` (dipisah koma, atau env
`AXBRIDGE_ORIGINS`), `-token-file`, `-pid-file`. Env `AXBRIDGE_TOKEN` menentukan token secara eksplisit.

## Pemecahan masalah

- **WebUI menunjukkan *Stopped*:** tap **Restart**, lalu **Run diagnostics** dan baca hasilnya. Isinya ABI, izin file,
  tes eksekusi, dan akhir log. Kalau binary tidak bisa dieksekusi dari folder plugin, script start otomatis mencoba
  lagi dari `/data/local/tmp`.
- **`Failed to fetch` dari situsmu:** server mati, origin belum ada di *Allowed websites* (harus persis, termasuk
  `https://` dan tanpa slash di akhir), atau browser memblokir akses jaringan lokal.
- **`wrong or missing token`:** salin ulang token dari WebUI.

## Kredit & lisensi

MIT. Konsep server awal: [RiProG-ADB-Panel](https://github.com/RiProG-id) oleh RiProG-id. Lihat [LICENSE](LICENSE).
