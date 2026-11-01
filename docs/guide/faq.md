# FAQ

## Can I use it from my laptop's browser?

No, on purpose. The server listens on `127.0.0.1` of the phone, so only pages opened on that phone can reach it.

## Do I need root?

No. Commands run with the privileges of the AxManager server, shell by default and root when AxManager is in root
mode.

## Can I turn off the token or the allowlist?

No. They are what make leaving a shell server running reasonably safe. If your site is only for your own use,
adding one origin takes a few seconds.

## Can I allow every website?

No. Wildcards are not supported on purpose. Add the exact origins you use.

## Why do I have to paste the token?

The token proves that the person who owns the phone approved this site. Without it, any allowed site could be
abused by a script injected into it. Pasting is a one-time step per site, because the client saves it in
`localStorage`.

## Does it work in Firefox?

Chromium browsers are what AxBridge is developed and tested against. Other browsers depend on how they treat calls from https pages to `http://127.0.0.1`.
Safari does not allow it.

## Can I use it from a PWA or an installed web app?

Yes if it runs in a Chromium engine on the same phone and its origin is allowed.

## Does the client work in Node?

Yes, with Node 18 or newer (it uses `fetch`). `fullScreen` and `toast` need a browser. The token is kept in memory
there.

## How do I list installed apps?

Use `exec("pm list packages -3")`. There is no `listPackages` because the native bridge is not available in a
browser. See [Migrate from kernelsu](../api/migrate-from-kernelsu).

## Is it affiliated with AxManager?

No. It is a community plugin that uses AxManager's plugin system.
