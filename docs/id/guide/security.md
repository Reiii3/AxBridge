# Model keamanan

Ini adalah shell yang bisa dijangkau lewat HTTP. Ia hanya aman karena lapisan di bawah ini, jadi tolong jangan
dihapus.

| Lapisan | Yang diblokir |
| --- | --- |
| Hanya bind ke `127.0.0.1` | Perangkat lain di Wi-Fi yang sama |
| **Token** bearer (acak, 192-bit, per perangkat) | Situs mana pun yang tidak tahu token |
| **Allowlist origin**, diatur pemilik perangkat di WebUI | Website sembarang yang memanggil jembatan dari browser pengguna |
| `Content-Type: application/json` wajib | CSRF "simple request" dengan `text/plain` |
| Pengecekan header **Host** | DNS rebinding |
| Timeout dan batas ukuran | Perintah yang lepas kendali |
| Kill seluruh process group | Proses anak yang tertinggal setelah timeout atau `kill()` |

## Batasan

| | `/exec` | `/spawn` |
| --- | --- | --- |
| Timeout default | 30 dtk | 300 dtk |
| Timeout maksimum | 120 dtk | 3600 dtk |
| Body request | 64 KB | 64 KB |
| Output | 1 MB (kelebihannya dibuang) | 16 MB (proses dimatikan jika lewat) |

## Token

- Dibuat saat pertama kali menyala, disimpan di `$AXERONDIR/axbridge/token` dengan mode 600.
- Dibandingkan dalam waktu konstan.
- Dikirim sebagai `Authorization: Bearer <token>`.
- Kamu bisa menentukan sendiri lewat environment variable `AXBRIDGE_TOKEN`.

::: danger Perlakukan token seperti password
Siapa pun yang memegang token **dan** berjalan dari origin yang diizinkan punya akses shell ke HP.
:::

## Untuk developer aplikasi

- Token ada di `localStorage` halamanmu. **Script pihak ketiga mana pun di halaman itu bisa membacanya** dan
  menjalankan perintah di HP pengguna. Jaga halaman itu bebas dari script tak tepercaya, kode iklan, dan analytics
  yang tidak kamu kendalikan.
- Bug cross-site scripting di situsmu menjadi remote shell di HP penggunamu. Anggap XSS sebagai masalah kritis.
- Jangan teruskan input bebas dari pengguna ke `exec`. Sediakan **aksi spesifik** dari kodemu sendiri
  (`restartGame()`, `clearCache()`) dan susun perintahnya sendiri.
- Dengan `spawn`, kirim array `args` supaya tidak ada shell yang mem-parse argumen.

## Yang tidak dilindungi AxBridge

- Aplikasi jahat yang sudah berjalan di HP yang sama dan memegang token. Aplikasi juga bisa menjangkau `127.0.0.1`,
  dan pengecekan origin hanya membatasi browser. Karena itu token ada dan harus dirahasiakan.
- Pemilik perangkat yang menambahkan situs tak tepercaya ke allowlist.

## Melaporkan masalah

Tolong buka security advisory privat di repositori GitHub, bukan issue publik.
