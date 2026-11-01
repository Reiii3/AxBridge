# Pasang plugin

1. Unduh `axbridge-<versi>.zip` dari [Releases](https://github.com/Reiii3/AxBridge/releases),
   atau [build sendiri](../dev/build).
2. Di AxManager, buka **Plugin** lalu pasang dari ZIP.
3. Buka **WebUI** plugin. Kalau kartu status bertuliskan *Stopped*, tekan **Restart server**.
4. Tekan **Open in browser**. Ini membuka halaman tes bawaan dengan token yang sudah terisi, jadi kamu bisa langsung
   mencoba perintah.

## Yang dilakukan installer

- Memilih binary yang cocok dengan CPU-mu (`arm64` atau `arm`) dan menghapus yang lain.
- Mengatur izin file.
- Saat pertama kali menyala, server membuat token acak di `$AXERONDIR/axbridge/token` (mode 600).

Server menyala lagi setiap boot. Log ada di `$AXERONDIR/logs/axbridge.log`.

::: warning
Plugin mendeklarasikan `axeronPlugin=19100`, dan AxManager mensyaratkan angka itu sama dengan atau lebih rendah dari
versi servernya. Kalau installer menolak ZIP, turunkan angka itu di `module.prop`.
:::

## Tombol Action

Tombol **Action** milik plugin me-restart server dan membuka halaman tes di browser, sama seperti
**Open in browser** di WebUI.

## Update

Pasang ZIP yang lebih baru di atas yang lama, lalu restart server. Token dan daftar situs ada di luar folder plugin,
jadi biasanya tetap tersimpan (lihat [Isi plugin](../dev/internals#files-on-the-device)).

## Copot

Hapus plugin di AxManager. Skrip uninstall menghentikan server dan menghapus folder datanya, termasuk token dan
daftar situs yang diizinkan.
