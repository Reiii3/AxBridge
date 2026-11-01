# Sourced by the AxManager installer (not executed).
ui_print "-------------------------------------------------------"
ui_print "- AxBridge $(sed -n 's/^version=//p' "$MODPATH/module.prop")"
ui_print "- Developer: $(sed -n 's/^author=//p' "$MODPATH/module.prop")"
ui_print "- Credit: $(sed -n 's/^credit=//p' "$MODPATH/module.prop")"
ui_print "- ARCH: $ARCH"
ui_print "- API: $API"
ui_print "- AXERONVER: $AXERONVER"
ui_print "-------------------------------------------------------"

case "$ARCH" in
  arm64) KEEP=axbridge-arm64; DROP=axbridge-arm ;;
  arm)   KEEP=axbridge-arm;   DROP=axbridge-arm64 ;;
  *) abort "! Unsupported architecture: $ARCH" ;;
esac

mv "$MODPATH/system/bin/$KEEP" "$MODPATH/system/bin/axbridge" || abort "! Binary missing"
rm -f "$MODPATH/system/bin/$DROP"

set_perm_recursive "$MODPATH" 0 0 0755 0644
set_perm "$MODPATH/system/bin/axbridge" 0 0 0755
for f in service.sh action.sh uninstall.sh ctl.sh; do
  set_perm "$MODPATH/$f" 0 0 0755
done

ui_print "- Done. Tap Action to open the test page in your browser."