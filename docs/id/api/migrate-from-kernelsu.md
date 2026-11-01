# Pindah dari kernelsu

Client AxBridge berbentuk sama dengan paket [`kernelsu`](https://www.npmjs.com/package/kernelsu), jadi project web
yang sudah memakai `kernelsu` biasanya hanya butuh satu baris konfigurasi.

## Yang sepadan satu lawan satu

| `kernelsu` | `axbridge-api` | Catatan |
| --- | --- | --- |
| `exec(cmd, { cwd, env })` | `exec(cmd, { cwd, env, timeout })` | Hasil sama `{ errno, stdout, stderr }` |
| `spawn(cmd, args, { cwd, env })` | `spawn(cmd, args, { cwd, env, timeout })` | Event sama, plus `kill()` sungguhan |
| `fullScreen(enable)` | `fullScreen(enable)` | Layar penuh browser, butuh interaksi pengguna |
| `toast(message)` | `toast(message)` | Toast DOM, bukan toast Android |

## Yang tidak tersedia

Ini butuh WebView milik aplikasi manager, yang tidak dimiliki browser.

| `kernelsu` | Di browser |
| --- | --- |
| `listPackages(type)` | `exec("pm list packages -3")` untuk aplikasi user, `-s` untuk sistem, tanpa flag untuk semua. Buang awalan `package:`. |
| `getPackagesInfo(packages)` | Parse `dumpsys package <nama>` untuk `versionName`, `versionCode`, dan `userId`. Jalankan sekali untuk banyak paket di dalam satu perintah. Label dan ikon aplikasi tidak tersedia. |
| `moduleInfo()`, `exit()`, `enableEdgeToEdge()` | Tidak ada padanannya |

## Opsi 1: alias (tanpa ubah kode)

::: code-group

```js [vite.config.js]
export default {
  resolve: { alias: { kernelsu: "axbridge-api" } },
};
```

```js [webpack.config.js]
module.exports = {
  resolve: { alias: { kernelsu: "axbridge-api" } },
};
```

:::

Setiap `import … from "kernelsu"` kini mengarah ke AxBridge di build itu.

## Opsi 2: satu aplikasi, dua lingkungan

Pakai `kernelsu` saat halamanmu berjalan di dalam WebView manager dan AxBridge di browser biasa:

```js
const api = window.ksu
  ? await import("kernelsu")
  : await import("axbridge-api");

const { stdout } = await api.exec("getprop ro.product.model");
```

Kedua modul punya `exec` dan `spawn` yang sama, jadi sisa kodemu tidak perlu peduli.

## Perbedaan yang perlu diingat

- Di browser, pengguna harus mengizinkan situsmu dan menempel token sekali. Lihat [Pakai di website-mu](../guide/quick-start).
- Tangani `BridgeError` (server mati, token salah, situs tidak diizinkan). `kernelsu` tidak punya padanannya, karena
  di dalam WebView tidak ada yang bisa gagal.
- `exec` berhenti paling lama setelah 120 dtk. Pakai `spawn` untuk pekerjaan panjang.
