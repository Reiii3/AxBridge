# Publish paket npm

Client di `client/` dipublish sebagai **`axbridge-api`**.

## Pertama kali di npm

npm menolak publish dari akun yang belum siap. Berurutan:

1. **Verifikasi email.** Jalankan `npm profile get`. Kalau email tertulis `(unverified)`, buka
   [npmjs.com/email-edit](https://www.npmjs.com/email-edit), kirim ulang email verifikasi, lalu klik linknya.
   Cek spam dan tab Promotions.
2. **Aktifkan two-factor authentication** untuk *Authorization and publishing*, atau pakai granular access token
   dengan *Bypass two-factor authentication*.
3. `npm login`.

Error yang umum:

| Pesan | Artinya |
| --- | --- |
| `Two-factor authentication or granular access token with bypass 2fa enabled is required` | Kerjakan langkah 2 |
| `you must verify your email before publishing a new package` | Kerjakan langkah 1 |

## Publish manual

```sh
cd client
npm pack --dry-run                  # cek daftar file: 5 file, sekitar 5 kB
npm publish --access public         # naikkan "version" di package.json dulu
```

## Publish dari GitHub

1. Buat token npm **Automation** (atau granular, read-and-write) dan simpan sebagai secret repositori `NPM_TOKEN`.
2. Pastikan `client/package.json` berisi versi yang kamu mau.
3. Push tag yang cocok dengannya:

```sh
git tag npm-v0.1.0
git push origin npm-v0.1.0
```

[Workflow](https://github.com/Reiii3/AxBridge/blob/main/.github/workflows/publish-npm.yml) menolak jalan kalau tag
dan `package.json` berbeda, lalu publish dengan provenance.

## Versi

- Client dan server diberi versi terpisah.
- `spawn`, `env`, `fullScreen`, dan `toast` butuh plugin **v0.1.6 atau lebih baru**.
- Jangan pakai ulang nomor versi. npm tidak mengizinkan publish ulang.
