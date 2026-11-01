---
layout: home

hero:
  name: AxBridge
  text: Shell AxManager dari browser mana pun
  tagline: Plugin kecil yang membuat tab browser biasa, atau situs https milikmu sendiri, bisa menjalankan perintah di HP tempat halaman itu dibuka.
  image:
    src: /logo.svg
    alt: AxBridge
  actions:
    - theme: brand
      text: Mulai
      link: /id/guide/what-is-axbridge
    - theme: alt
      text: API client
      link: /id/api/client
    - theme: alt
      text: GitHub
      link: https://github.com/Reiii3/AxBridge

features:
  - icon: 🔌
    title: API mirip kernelsu
    details: exec, spawn, fullScreen, dan toast berbentuk sama dengan paket kernelsu. Sering kali cukup satu baris alias.
  - icon: 📡
    title: spawn streaming
    details: Output datang selagi proses berjalan, dan kill() menghentikannya di HP beserta semua proses anaknya.
  - icon: 🔒
    title: Aman secara default
    details: Hanya loopback, token per perangkat, allowlist origin yang kamu atur dari WebUI, plus pengecekan DNS-rebinding dan CSRF.
  - icon: 📦
    title: Pasang dari npm
    details: "npm i axbridge-api. Lima kilobyte, tanpa dependensi, sudah ada tipe TypeScript."
  - icon: 🎨
    title: WebUI gaya Settings
    details: Atur server, situs yang diizinkan, dan token dari layar Material 3 di dalam AxManager, terang dan gelap.
  - icon: 🧩
    title: Kode polos
    details: Server Go dan WebUI dari HTML, CSS, dan ES module. Plugin tanpa build step, mudah dibaca dan di-fork.
---
