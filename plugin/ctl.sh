#!/system/bin/sh
# Helper for the WebUI (keeps quoting out of JavaScript).
MODDIR=$(cd "${0%/*}" && pwd)   # absolute, even when called as ./script.sh
DATA="$AXERONDIR/axbridge"
OFILE="$DATA/origins"
PORT=8787

# https://host[:port]   or   http://localhost|127.0.0.1[:port]
valid_origin() {
  printf '%s' "$1" | grep -Eq '^(https://[A-Za-z0-9]([A-Za-z0-9.-]*[A-Za-z0-9])?|http://(localhost|127\.0\.0\.1))(:[0-9]{1,5})?$'
}

case "$1" in
  token)   cat "$DATA/token" 2>/dev/null ;;
  status)  if kill -0 "$(cat "$DATA/pid" 2>/dev/null)" 2>/dev/null; then echo running; else echo stopped; fi ;;
  restart)
    sh "$MODDIR/service.sh" || { tail -n 4 "$AXERONDIR/logs/axbridge.log" >&2; exit 1; } ;;
  diag)
    echo "== env";   echo "AXERONDIR=$AXERONDIR"; echo "AXERON=$AXERON AXERONVER=$AXERONVER"; id
    echo "abi=$(getprop ro.product.cpu.abi) uname=$(uname -m)"
    echo "== files"; ls -l "$MODDIR/system/bin/" 2>&1; ls -ld "$DATA" 2>&1; ls -l "$DATA" 2>&1
    echo "== exec test"; eo=$("$MODDIR/system/bin/axbridge" -h 2>&1); erc=$?; echo "$eo" | head -n 3; echo "rc=$erc (126/127 = cannot execute)"
    echo "setsid: $(command -v setsid || echo none)"
    echo "== log (last 15)"; tail -n 15 "$AXERONDIR/logs/axbridge.log" 2>&1 ;;
  open)    am start -a android.intent.action.VIEW -d "http://127.0.0.1:$PORT/#token=$(cat "$DATA/token")" ;;
  origins-list)
    [ -f "$OFILE" ] && grep -v '^[[:space:]]*#' "$OFILE" | tr -d '\r' | grep -v '^[[:space:]]*$'
    exit 0 ;;
  origins-add)
    o=$(printf '%s' "$2" | tr 'A-Z' 'a-z')
    valid_origin "$o" || { echo "invalid origin" >&2; exit 3; }
    mkdir -p "$DATA"; touch "$OFILE"
    grep -qxF "$o" "$OFILE" || echo "$o" >> "$OFILE" ;;
  origins-del)
    o=$(printf '%s' "$2" | tr 'A-Z' 'a-z')
    valid_origin "$o" || { echo "invalid origin" >&2; exit 3; }
    [ -f "$OFILE" ] || exit 0
    grep -vxF "$o" "$OFILE" > "$OFILE.tmp"; mv "$OFILE.tmp" "$OFILE" ;;
  *) echo "usage: ctl.sh token|status|restart|diag|open|origins-list|origins-add <o>|origins-del <o>"; exit 2 ;;
esac
