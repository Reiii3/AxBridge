---
layout: home

hero:
  name: AxBridge
  text: AxManager shell, from any browser
  tagline: A small plugin that lets a normal browser tab, or your own hosted https site, run commands on the phone it is opened on.
  image:
    src: /logo.svg
    alt: AxBridge
  actions:
    - theme: brand
      text: Get started
      link: /guide/what-is-axbridge
    - theme: alt
      text: Client API
      link: /api/client
    - theme: alt
      text: GitHub
      link: https://github.com/Reiii3/AxBridge

features:
  - icon: 🔌
    title: kernelsu-compatible API
    details: exec, spawn, fullScreen and toast have the same shape as the kernelsu package. Often a one-line alias is all you need.
  - icon: 📡
    title: Streaming spawn
    details: Output arrives while the process runs, and kill() stops it on the device, including everything it started.
  - icon: 🔒
    title: Secure by default
    details: Loopback only, a per-device token, an origin allowlist you manage from the WebUI, DNS-rebinding and CSRF checks.
  - icon: 📦
    title: Install from npm
    details: "npm i axbridge-api. Five kilobytes, no dependencies, TypeScript types included."
  - icon: 🎨
    title: Settings-style WebUI
    details: Manage the server, allowed websites and the token from a Material 3 screen inside AxManager, light and dark.
  - icon: 🧩
    title: Plain code
    details: A Go server and a WebUI made of HTML, CSS and ES modules. No build step for the plugin, easy to read and fork.
---
