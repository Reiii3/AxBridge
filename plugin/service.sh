#!/system/bin/sh
MODDIR=$(cd "${0%/*}" && pwd)   # absolute, even when called as ./script.sh
DATA="$AXERONDIR/axbridge"
LOG="$AXERONDIR/logs/axbridge.log"
PIDF="$DATA/pid"
PORT=8787
FALLBACK="${AXBRIDGE_TMP:-/data/local/tmp}/axbridge"

mkdir -p "$DATA" "$AXERONDIR/logs"
log() { echo "$(date '+%F %T') [service] $*" >> "$LOG"; }

# stop previous instance (it writes its own pid file)
if [ -f "$PIDF" ]; then
  kill "$(cat "$PIDF")" 2>/dev/null
  rm -f "$PIDF"
  sleep 1
fi

# allowed browser origins: one per line in $DATA/origins (managed from the WebUI)
ORIGINS=""
[ -f "$DATA/origins" ] && ORIGINS=$(grep -v '^[[:space:]]*#' "$DATA/origins" | tr -d ' \r' | tr '\n' ',')

# 126/127 = "cannot execute" / "not found" (permission, noexec mount, wrong ABI ...)
cannot_exec() { "$1" -h >/dev/null 2>&1; rc=$?; [ "$rc" -eq 126 ] || [ "$rc" -eq 127 ]; }

BIN="$MODDIR/system/bin/axbridge"
chmod 755 "$BIN" 2>/dev/null
if cannot_exec "$BIN"; then
  log "cannot exec $BIN (rc=$rc), trying $FALLBACK"
  mkdir -p "${FALLBACK%/*}" 2>/dev/null
  if cp "$BIN" "$FALLBACK" 2>>"$LOG" && chmod 755 "$FALLBACK" 2>>"$LOG"; then BIN="$FALLBACK"; fi
  if cannot_exec "$BIN"; then
    log "FAILED: $BIN cannot be executed either (rc=$rc). Wrong ABI or exec blocked by SELinux/noexec."
    exit 1
  fi
fi

cd "$DATA" || { log "FAILED: cannot cd to $DATA"; exit 1; }
log "starting $BIN origins=$(echo "$ORIGINS" | tr ',' ' ')"
# setsid (when available) detaches the daemon from the caller's process group
if command -v setsid >/dev/null 2>&1; then
  setsid "$BIN" -addr "127.0.0.1:$PORT" -origins "$ORIGINS" -token-file "$DATA/token" -pid-file "$PIDF" >>"$LOG" 2>&1 &
else
  "$BIN" -addr "127.0.0.1:$PORT" -origins "$ORIGINS" -token-file "$DATA/token" -pid-file "$PIDF" >>"$LOG" 2>&1 &
fi

sleep 1
if kill -0 "$(cat "$PIDF" 2>/dev/null)" 2>/dev/null; then
  log "running, pid $(cat "$PIDF")"
else
  log "FAILED: server exited right after start (port $PORT busy? see lines above)"
  exit 1
fi
