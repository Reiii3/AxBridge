# Allowed websites

Browser mengirim header `Origin` di setiap request lintas situs. AxBridge hanya menjawab origin yang sudah ditambahkan
pemilik perangkat ke allowlist, jadi situs acak yang kamu buka tidak bisa memanggil jembatan, walau menebak port-nya.

## Tambah situs

Di WebUI plugin, **Allowed websites** → ketik situsnya → **Add**. Server restart sendiri supaya perubahannya berlaku.

Kamu boleh menempel URL lengkap. Path, query, dan hash dibuang, dan host dijadikan huruf kecil:

| Yang kamu ketik | Tersimpan sebagai |
| --- | --- |
| `New-Site.Example.com/some/path?x=1` | `https://new-site.example.com` |
| `https://app.example.com:8443/x` | `https://app.example.com:8443` |
| `http://localhost:5173/` | `http://localhost:5173` |

## Aturan

- `https://host` atau `https://host:port`.
- `http` biasa hanya diterima untuk `localhost` dan `127.0.0.1`, untuk pengembangan lokal.
- Origin yang persis. Tanpa wildcard, tanpa path, tanpa garis miring di akhir.
- `https://app.example.com` dan `https://www.app.example.com` adalah origin berbeda. Tambahkan keduanya kalau
  dipakai dua-duanya.

## Selalu diizinkan

Halaman tes bawaan (`http://127.0.0.1:8787` dan `http://localhost:8787`) selalu boleh memanggil server. Itu yang
dipakai tombol **Open in browser**.

## Hapus situs

Tekan ikon tempat sampah di barisnya. Server restart lagi.

## Disimpan di mana

Satu origin per baris di `$AXERONDIR/axbridge/origins`. Cara yang didukung untuk mengeditnya adalah WebUI. Kamu juga
bisa memakai helper `ctl.sh` dari shell, lihat [Isi plugin](../dev/internals#ctl-sh).
