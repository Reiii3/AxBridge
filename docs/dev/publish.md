# Publish the npm package

The client in `client/` is published as **`axbridge-api`**.

## First time on npm

npm refuses to publish from an account that is not ready. In order:

1. **Verify your email.** Run `npm profile get`. If the email shows `(unverified)`, open
   [npmjs.com/email-edit](https://www.npmjs.com/email-edit), resend the verification mail and click the link.
   Check spam and the Promotions tab.
2. **Turn on two-factor authentication** for *Authorization and publishing*, or use a granular access token with
   *Bypass two-factor authentication*.
3. `npm login`.

Typical errors:

| Message | Meaning |
| --- | --- |
| `Two-factor authentication or granular access token with bypass 2fa enabled is required` | Do step 2 |
| `you must verify your email before publishing a new package` | Do step 1 |

## Publish by hand

```sh
cd client
npm pack --dry-run                  # check the file list: 5 files, about 5 kB
npm publish --access public         # raise "version" in package.json first
```

## Publish from GitHub

1. Create an npm **Automation** (or granular, read-and-write) token and save it as the repository secret `NPM_TOKEN`.
2. Make sure `client/package.json` has the version you want.
3. Push a tag that matches it:

```sh
git tag npm-v0.1.0
git push origin npm-v0.1.0
```

The [workflow](https://github.com/Reiii3/AxBridge/blob/main/.github/workflows/publish-npm.yml) refuses to run if
the tag and `package.json` disagree, then publishes with provenance.

## Versioning

- The client and the server are versioned separately.
- `spawn`, `env`, `fullScreen` and `toast` need the plugin **v0.1.6 or newer**.
- Never reuse a version number. npm will not let you republish it.
