# FAQ

## Bisa dipakai dari browser laptop?

Tidak, memang disengaja. Server mendengarkan di `127.0.0.1` milik HP, jadi hanya halaman yang dibuka di HP itu yang
bisa menjangkaunya.

## Perlu root?

Tidak. Perintah berjalan dengan hak server AxManager: shell secara default, dan root kalau AxManager dalam mode root.

## Bisa mematikan token atau allowlist?

Tidak. Dua hal itu yang membuat server shell yang dibiarkan menyala cukup aman. Kalau situsmu hanya untuk dipakai
sendiri, menambah satu origin hanya butuh beberapa detik.

## Bisa mengizinkan semua website?

Tidak. Wildcard sengaja tidak didukung. Tambahkan origin persis yang kamu pakai.

## Kenapa harus menempel token?

Token membuktikan bahwa pemilik HP menyetujui situs ini. Tanpanya, situs yang diizinkan bisa disalahgunakan lewat
script yang disusupkan ke dalamnya. Menempel cukup sekali per situs, karena client menyimpannya di `localStorage`.

## Jalan di Firefox?

AxBridge dikembangkan dan diuji di browser Chromium. Browser lain bergantung pada cara mereka memperlakukan panggilan
dari halaman https ke `http://127.0.0.1`. Safari tidak mengizinkannya.

## Bisa dari PWA atau aplikasi web yang diinstal?

Bisa, kalau berjalan di mesin Chromium pada HP yang sama dan origin-nya diizinkan.

## Client jalan di Node?

Ya, dengan Node 18 atau lebih baru (memakai `fetch`). `fullScreen` dan `toast` butuh browser. Token disimpan di
memori di sana.

## Bagaimana melihat daftar aplikasi terpasang?

Pakai `exec("pm list packages -3")`. Tidak ada `listPackages` karena jembatan native tidak tersedia di browser.
Lihat [Pindah dari kernelsu](../api/migrate-from-kernelsu).

## Apakah ini resmi dari AxManager?

Bukan. Ini plugin komunitas yang memakai sistem plugin AxManager.
