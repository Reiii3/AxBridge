#!/system/bin/sh
DATA="$AXERONDIR/axbridge"
[ -f "$DATA/pid" ] && kill "$(cat "$DATA/pid")" 2>/dev/null
rm -rf "$DATA"
