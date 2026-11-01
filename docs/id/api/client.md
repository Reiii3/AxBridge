# API client

```sh
npm i axbridge-api
```

```js
import { exec, spawn, fullScreen, toast, ping, setToken } from "axbridge-api";
```

Paket ini hanya ESM, tanpa dependensi, dan membawa tipe TypeScript sendiri.

## exec

Menjalankan perintah dengan `sh -c` dan menunggu sampai selesai.

```ts
exec(command: string, options?: {
  cwd?: string;
  env?: Record<string, string>;
  timeout?: number; // detik, default 30, maks 120
}): Promise<{ errno: number; stdout: string; stderr: string }>
```

```js
const { errno, stdout, stderr } = await exec("ls -l", { cwd: "/data/local/tmp" });
if (errno === 0) console.log(stdout);
```

- `errno` adalah exit code. `124` berarti perintah kena timeout.
- Perintah yang gagal tetap **resolve**. Hanya masalah pada jembatan yang reject (lihat [BridgeError](#bridgeerror)).
- `env` mengatur environment variable untuk perintah ini. Nama harus cocok dengan `[A-Za-z_][A-Za-z0-9_]*`.

## spawn

Menjalankan proses dan mengalirkan outputnya.

```ts
spawn(command: string, args?: string[], options?: {
  cwd?: string;
  env?: Record<string, string>;
  timeout?: number; // detik, default 300, maks 3600
}): ChildProcess

spawn(command: string, options?: SpawnOptions): ChildProcess
```

```js
const ls = spawn("ls", ["-lh", "/data"]);

ls.stdout.on("data", (text) => console.log(`stdout: ${text}`));
ls.stderr.on("data", (text) => console.log(`stderr: ${text}`));
ls.on("exit", (code) => console.log(`selesai dengan ${code}`));
ls.on("error", (error) => console.error(error.code));
```

### Baris shell atau langsung

| Panggilan | Dijalankan sebagai |
| --- | --- |
| `spawn("logcat -d \| grep Zygote")` | Baris shell (`sh -c`), pipe dan redirect berfungsi |
| `spawn("logcat", ["-d"])` | `logcat` **langsung**, tanpa shell. Argumen tidak pernah di-parse, jadi tidak ada masalah quoting atau injeksi |

Pakai bentuk `args` setiap kali sebagian perintah berasal dari input pengguna.

### ChildProcess

| Anggota | |
| --- | --- |
| `stdout.on("data", (text: string) => …)` | Potongan stdout, sudah didekode sebagai UTF-8. Karakter yang terbelah di dua potongan disambung dengan benar. |
| `stderr.on("data", (text: string) => …)` | Sama untuk stderr. |
| `on("exit", (code: number \| null) => …)` | Terpanggil sekali, selalu terakhir. `code` bernilai `null` kalau kamu memanggil `kill()` atau koneksi terputus. Timeout dari server memberi `124`. |
| `on("error", (error) => …)` | Jembatan gagal (lihat [BridgeError](#bridgeerror)). `exit` menyusul dengan `null`. |
| `kill()` | Menghentikan proses di HP, **termasuk semua yang dijalankannya**. |

::: tip Output datang selagi proses berjalan
Kamu tidak perlu menunggu `exit`. Progress bar, penampil log, atau `logcat` langsung berfungsi seperti yang kamu harapkan.
:::

Menutup tab, atau cara lain apa pun yang memutus koneksi, juga mematikan prosesnya.

## fullScreen

```ts
fullScreen(enable: boolean): Promise<void>
```

Masuk atau keluar dari layar penuh browser. Browser mewajibkan interaksi pengguna (klik atau tap). Kalau browser
menolak, fungsi ini selesai diam-diam tanpa melempar error.

## toast

```ts
toast(message: string): void
```

Menampilkan pesan singkat di bagian bawah halaman. Berupa elemen DOM biasa dengan style inline, jadi tidak butuh CSS.
Hanya di browser.

## Helper koneksi

| Fungsi | |
| --- | --- |
| `ping(): Promise<boolean>` | `true` kalau server menjawab. Pakai untuk menampilkan "AxBridge belum jalan". |
| `setToken(token: string)` | Menyimpan token. |
| `getToken(): string` | Membacanya kembali (string kosong kalau belum ada). |
| `clearToken()` | Melupakannya. Panggil saat menerima `unauthorized` dan kamu ingin meminta ulang ke pengguna. |
| `setBaseUrl(url: string)` | Pakai alamat selain `http://127.0.0.1:8787`. |

Token disimpan di `localStorage`. Kalau penyimpanan diblokir (mode privat, iframe sandbox) atau tidak ada (Node),
client menyimpannya di memori selama halaman hidup.

## BridgeError

Dilempar oleh `exec` dan di-emit sebagai event `error` oleh `spawn` saat **jembatan** tidak bisa mengerjakan tugasnya.

```ts
class BridgeError extends Error {
  code: "network" | "unauthorized" | "forbidden" | "http";
}
```

| `code` | Arti | Perbaikan umum |
| --- | --- | --- |
| `network` | Server tidak terjangkau | Nyalakan plugin, atau browser memblokir akses jaringan lokal |
| `unauthorized` | Token salah atau belum diisi | `setToken(...)` dengan token dari WebUI |
| `forbidden` | Situsmu belum ada di **Allowed websites** | Tambahkan di WebUI plugin |
| `http` | Error HTTP lainnya | Pesannya berisi status dan body |

```js
try {
  await exec("id");
} catch (error) {
  if (error.code === "unauthorized") {
    clearToken();
    // minta pengguna menempel token lagi
  }
}
```

## TypeScript

```ts
import { exec, spawn, type ExecResult, BridgeError } from "axbridge-api";

const result: ExecResult = await exec("id", { env: { A: "1" }, timeout: 5 });
```

## Node

Jalan di Node 18 atau lebih baru (memakai `fetch` global), untuk script yang berjalan **di komputer yang bisa
menjangkau port HP**, misalnya lewat `adb forward tcp:8787 tcp:8787`. `fullScreen` dan `toast` butuh browser.
