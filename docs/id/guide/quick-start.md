# Pakai di website-mu

## 1. Izinkan situsmu

Di WebUI plugin, buka **Allowed websites** dan tambahkan origin situsmu, misalnya `https://myapp.vercel.app`.
Untuk dev lokal di HP yang sama pakai `http://localhost:5173`.
Detailnya ada di [Allowed websites](./allowed-websites).

## 2. Pasang client

::: code-group

```sh [npm]
npm i axbridge-api
```

```sh [pnpm]
pnpm add axbridge-api
```

```sh [yarn]
yarn add axbridge-api
```

:::

Tanpa bundler? Salin [`axbridge.js`](https://github.com/Reiii3/AxBridge/blob/main/client/axbridge.js)
ke project-mu dan import lewat path.

## 3. Beri token ke situsmu

Di WebUI, buka **Token** lalu tekan tombol salin. Tempel ke situsmu sekali saja. Client menyimpannya di
`localStorage`.

```js
import { setToken } from "axbridge-api";

setToken("<token dari WebUI AxBridge>");
```

## 4. Jalankan perintah

```js
import { exec, spawn, ping } from "axbridge-api";

if (!(await ping())) {
  alert("AxBridge belum jalan. Buka plugin lalu tekan Restart server.");
}

// Perintah sekali jalan
const { errno, stdout, stderr } = await exec("getprop ro.product.model");
console.log(stdout.trim());

// Perintah panjang, output datang selagi berjalan
const child = spawn("logcat", ["-v", "brief"]);
child.stdout.on("data", (text) => console.log(text));
child.on("exit", (code) => console.log("selesai", code));

// nanti
child.kill();
```

`errno` adalah exit code perintah. Perintah yang gagal **bukan** exception. Exception
(`BridgeError`) berarti jembatannya sendiri tidak terjangkau atau menolakmu.
Lihat [API client](../api/client#bridgeerror).

## Perlu diketahui

- Halaman harus dibuka **di browser pada HP yang sama**. `127.0.0.1` adalah HP itu sendiri.
- Chrome bisa meminta izin "akses jaringan lokal" saat pertama kali. Izinkan.
- Jaga halaman yang memakai AxBridge bebas dari script pihak ketiga yang tidak tepercaya: script apa pun di halaman
  itu bisa membaca token.
- Lebih baik sediakan **aksi spesifik** dari kodemu sendiri daripada meneruskan input pengguna ke `exec`.
