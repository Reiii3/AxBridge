# Install the plugin

1. Download `axbridge-<version>.zip` from [Releases](https://github.com/Reiii3/AxBridge/releases),
   or [build it yourself](../dev/build).
2. In AxManager, open **Plugin** and install from ZIP.
3. Open the plugin's **WebUI**. If the status card says *Stopped*, tap **Restart server**.
4. Tap **Open in browser**. This opens the built-in test page with the token already filled in, so you can try a
   command straight away.

## What the installer does

- Picks the right binary for your CPU (`arm64` or `arm`) and removes the other.
- Sets file permissions.
- On first start the server creates a random token in `$AXERONDIR/axbridge/token` (mode 600).

The server starts again at every boot. Logs go to `$AXERONDIR/logs/axbridge.log`.

::: warning
The plugin declares `axeronPlugin=19100`, and AxManager requires that to be equal to or lower than its server
version. If the installer refuses the ZIP, lower that number in `module.prop`.
:::

## The Action button

The plugin's **Action** button restarts the server and opens the test page in your browser, the same as
**Open in browser** in the WebUI.

## Update

Install the newer ZIP over the old one and restart the server. Your token and allowed websites live outside the
plugin folder, so they are normally kept (see [Plugin internals](../dev/internals#files-on-the-device)).

## Uninstall

Remove the plugin in AxManager. The uninstall script stops the server and deletes its data folder, which includes
the token and the allowed websites list.
