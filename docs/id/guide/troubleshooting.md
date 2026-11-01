# Pemecahan masalah

## WebUI menampilkan *Stopped*

1. Tekan **Restart server**.
2. Masih berhenti? Tekan **Run diagnostics** dan baca hasilnya. Isinya ABI CPU, izin file, uji jalan binary, dan
   ekor log.
3. Kalau binary tidak bisa dieksekusi dari folder plugin, skrip start otomatis mencoba lagi dari
   `/data/local/tmp`.
4. Tekan **Copy** di kotak diagnostik dan sertakan kalau kamu membuka issue.

## `Failed to fetch` dari situsmu

Periksa berurutan:

1. Server berjalan (WebUI menampilkan *Running*).
2. Origin-mu ada di **Allowed websites**, **persis**: pakai `https://`, tanpa garis miring di akhir, tanpa path.
3. Kamu membuka halaman di browser **pada HP yang sama**.
4. Chrome tidak memblokir "akses jaringan lokal". Buka pengaturan situs untuk halamanmu dan izinkan.
5. Kamu tidak memakai Safari atau browser dalam aplikasi yang memblokir panggilan ke `127.0.0.1` dari halaman https.

## `BridgeError: unauthorized`

Token salah atau belum diisi. Salin lagi dari WebUI (**Token** → salin) lalu panggil `setToken(...)`.
Restart server **tidak** mengganti token.

## `BridgeError: forbidden`

Situs belum ada di allowlist. Tambahkan, tunggu server restart, lalu coba lagi.

## Perintah mengembalikan `errno: 124`

Kena timeout. `exec` maksimal 120 dtk. Untuk pekerjaan panjang pakai `spawn` (sampai 3600 dtk).

## `spawn` tidak ada atau mengembalikan 404

Pluginmu lebih lama dari v0.1.6. Pasang ZIP yang lebih baru, karena `/spawn` adalah fitur server.

## Output terpotong

`exec` menyimpan paling banyak 1 MB output. Pakai `spawn` untuk output besar.

## Server berhenti setelah beberapa lama

Sebagian ROM mematikan proses latar belakang. Matikan optimasi baterai untuk AxManager dan biarkan plugin aktif.
Server menyala lagi saat boot, dan **Restart server** mengembalikannya kapan saja.

## Membaca log

```sh
cat "$AXERONDIR/logs/axbridge.log"
```
