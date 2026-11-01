# Apa itu AxBridge?

AxManager menyediakan objek JavaScript `ksu` dan `Axeron`, tetapi **hanya di dalam jendela WebUI milik plugin**.
Halaman yang dibuka di Chrome, atau situsmu di Vercel atau GitHub Pages, tidak bisa melihatnya, jadi tidak bisa
menjalankan perintah apa pun.

AxBridge adalah plugin AxManager yang menutup celah itu. Ia menjalankan server HTTP kecil di perangkat,
`127.0.0.1:8787`, dan halaman apa pun yang dibuka di browser HP bisa memintanya menjalankan perintah shell:

```
 Chrome / situsmu (https)              AxBridge (di HP)                      Shell Android
 ───────────────────────────  POST   ───────────────────────────  sh -c   ──────────────
 exec("getprop ro.build...")  ─────▶  127.0.0.1:8787/exec        ───────▶  jalankan perintah
                              ◀─────  { errno, stdout, stderr }  ◀───────
```

Perintah berjalan dengan hak yang sama seperti server AxManager: **shell**, atau **root** kalau AxManager dalam
mode root.

## Kenapa tidak dibuka saja port-nya?

Shell lewat HTTP yang terbuka berarti situs mana pun yang kamu kunjungi bisa menjalankan perintah di HP-mu.
AxBridge menambahkan lapisan yang membuatnya aman dibiarkan menyala: token, allowlist situs, dan pengecekan request.
Baca [Model keamanan](./security) sebelum membangun sesuatu di atasnya.

## Yang kamu dapat

| | |
| --- | --- |
| **Plugin** | Dipasang dari ZIP di AxManager. Menyalakan server saat boot, dengan WebUI untuk restart, mengatur situs yang diizinkan, dan menyalin token. |
| **Client** | [`axbridge-api`](https://www.npmjs.com/package/axbridge-api) di npm. `exec`, `spawn`, `fullScreen`, `toast`, bentuknya sama dengan `kernelsu`. |
| **API HTTP** | Dua endpoint, `/exec` dan `/spawn`, kalau kamu lebih suka memanggilnya dari bahasa lain. |

## Syarat

- Android dengan **AxManager** (arm64 atau arm).
- Browser **di HP yang sama**. Disarankan Chrome atau browser Chromium lain. Safari tidak bisa memanggil
  `http://127.0.0.1` dari halaman https.
- Plugin mendeklarasikan `axeronPlugin=19100`. Turunkan di `module.prop` kalau server AxManager-mu lebih lama.

::: info BEDA DENGAN KERNELSU
`kernelsu` bicara ke jembatan native WebView, jadi hanya jalan di dalam aplikasi manager.
AxBridge bicara lewat HTTP, jadi jalan di browser mana pun, tetapi plugin harus menyala dan kamu butuh token.
Fungsi yang butuh WebView itu sendiri (`listPackages`, `getPackagesInfo`, `exit`, `moduleInfo`,
`enableEdgeToEdge`) tidak tersedia. Lihat [Pindah dari kernelsu](../api/migrate-from-kernelsu).
:::

## Kredit

Berdasarkan [RiProG-ADB-Panel](https://github.com/RiProG-id) (MIT, © 2025 RiProG-id).
AxBridge menambahkan autentikasi, allowlist origin, penanganan CORS dan Private Network Access, batasan,
`spawn` streaming, dan plugin AxManager di sekelilingnya.
