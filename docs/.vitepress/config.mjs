import { defineConfig } from "vitepress";

// Change these two when the repository is published.
const GITHUB = "https://github.com/Reiii3/AxBridge";
const NPM = "https://www.npmjs.com/package/axbridge-api";

// GitHub Pages serves a project site under /<repo>/. Override with DOCS_BASE=/ for a custom domain.
const base = process.env.DOCS_BASE || "/AxBridge/";

// Public address of the site (no trailing slash). Used for the sitemap and canonical links.
const SITE = process.env.DOCS_SITE || "https://reiii3.github.io/AxBridge";

/** Sidebar for one language. `prefix` is "" for English and "/id" for Indonesian. */
function sidebar(prefix, t) {
  return [
    {
      text: t.start,
      items: [
        { text: t.what, link: `${prefix}/guide/what-is-axbridge` },
        { text: t.install, link: `${prefix}/guide/install` },
        { text: t.quick, link: `${prefix}/guide/quick-start` },
      ],
    },
    {
      text: t.using,
      items: [
        { text: t.allowed, link: `${prefix}/guide/allowed-websites` },
        { text: t.security, link: `${prefix}/guide/security` },
        { text: t.trouble, link: `${prefix}/guide/troubleshooting` },
        { text: "FAQ", link: `${prefix}/guide/faq` },
      ],
    },
    {
      text: "API",
      items: [
        { text: t.client, link: `${prefix}/api/client` },
        { text: t.http, link: `${prefix}/api/http` },
        { text: t.migrate, link: `${prefix}/api/migrate-from-kernelsu` },
      ],
    },
    {
      text: t.dev,
      items: [
        { text: t.build, link: `${prefix}/dev/build` },
        { text: t.internals, link: `${prefix}/dev/internals` },
        { text: t.publish, link: `${prefix}/dev/publish` },
      ],
    },
  ];
}

function nav(prefix, t) {
  return [
    { text: t.navGuide, link: `${prefix}/guide/what-is-axbridge`, activeMatch: "/guide/" },
    { text: "API", link: `${prefix}/api/client`, activeMatch: "/api/" },
    { text: t.navDev, link: `${prefix}/dev/build`, activeMatch: "/dev/" },
    { text: t.navDownload, link: `${GITHUB}/releases` },
  ];
}

const en = {
  start: "Getting started",
  what: "What is AxBridge?",
  install: "Install the plugin",
  quick: "Use it in your website",
  using: "Using AxBridge",
  allowed: "Allowed websites",
  security: "Security model",
  trouble: "Troubleshooting",
  client: "Client API",
  http: "HTTP protocol",
  migrate: "Migrate from kernelsu",
  dev: "For developers",
  build: "Build from source",
  internals: "Plugin internals",
  publish: "Publish the npm package",
  navGuide: "Guide",
  navDev: "Develop",
  navDownload: "Download",
};

const id = {
  start: "Mulai",
  what: "Apa itu AxBridge?",
  install: "Pasang plugin",
  quick: "Pakai di website-mu",
  using: "Memakai AxBridge",
  allowed: "Allowed websites",
  security: "Model keamanan",
  trouble: "Pemecahan masalah",
  client: "API client",
  http: "Protokol HTTP",
  migrate: "Pindah dari kernelsu",
  dev: "Untuk developer",
  build: "Build dari source",
  internals: "Isi plugin",
  publish: "Publish paket npm",
  navGuide: "Panduan",
  navDev: "Develop",
  navDownload: "Unduh",
};

export default defineConfig({
  title: "AxBridge",
  base,
  cleanUrls: false,
  lastUpdated: false,
  sitemap: { hostname: `${SITE}/` },

  // One canonical address per page, so search engines do not treat /x and /x.html as different pages.
  transformPageData(pageData) {
    const path = pageData.relativePath.replace(/index\.md$/, "").replace(/\.md$/, ".html");
    pageData.frontmatter.head ??= [];
    pageData.frontmatter.head.push(["link", { rel: "canonical", href: `${SITE}/${path}` }]);
  },
  head: [
    ["link", { rel: "icon", type: "image/svg+xml", href: `${base}logo.svg` }],
    ["meta", { name: "theme-color", content: "#0b57d0" }],
    ["meta", { property: "og:type", content: "website" }],
    ["meta", { property: "og:site_name", content: "AxBridge" }],
  ],

  locales: {
    root: {
      label: "English",
      lang: "en",
      description: "Run AxManager shell commands from any browser.",
      themeConfig: {
        nav: nav("", en),
        sidebar: sidebar("", en),
        editLink: { pattern: `${GITHUB}/edit/main/docs/:path`, text: "Edit this page on GitHub" },
      },
    },
    id: {
      label: "Bahasa Indonesia",
      lang: "id",
      link: "/id/",
      description: "Jalankan shell AxManager dari browser biasa.",
      themeConfig: {
        nav: nav("/id", id),
        sidebar: sidebar("/id", id),
        editLink: { pattern: `${GITHUB}/edit/main/docs/:path`, text: "Edit halaman ini di GitHub" },
        outline: { label: "Di halaman ini" },
        docFooter: { prev: "Sebelumnya", next: "Berikutnya" },
        darkModeSwitchLabel: "Tampilan",
        sidebarMenuLabel: "Menu",
        returnToTopLabel: "Kembali ke atas",
        langMenuLabel: "Ganti bahasa",
        lightModeSwitchTitle: "Ganti ke mode terang",
        darkModeSwitchTitle: "Ganti ke mode gelap",
        skipToContentLabel: "Lewati ke konten",
        notFound: {
          title: "HALAMAN TIDAK DITEMUKAN",
          quote: "Halaman yang kamu cari tidak ada atau sudah dipindahkan.",
          linkLabel: "Ke beranda",
          linkText: "Ke beranda",
        },
      },
    },
  },

  themeConfig: {
    logo: "/logo.svg",
    socialLinks: [
      { icon: "github", link: GITHUB },
      { icon: "npm", link: NPM },
    ],
    search: {
      provider: "local",
      options: {
        locales: {
          id: {
            translations: {
              button: { buttonText: "Cari", buttonAriaLabel: "Cari" },
              modal: {
                noResultsText: "Tidak ada hasil untuk",
                resetButtonTitle: "Hapus pencarian",
                footer: { selectText: "pilih", navigateText: "pindah", closeText: "tutup" },
              },
            },
          },
        },
      },
    },
    footer: {
      message: "Released under the MIT License. Based on RiProG-ADB-Panel (MIT).",
      copyright: "AxBridge is a community plugin for AxManager.",
    },
  },
});
