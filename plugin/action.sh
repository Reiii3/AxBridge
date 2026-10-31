#!/system/bin/sh
# Action button: (re)start the server and open the test page with the token pre-filled.
MODDIR=$(cd "${0%/*}" && pwd)   # absolute, even when called as ./script.sh
DATA="$AXERONDIR/axbridge"
PORT=8787

sh "$MODDIR/service.sh"
sleep 1
TOKEN=$(cat "$DATA/token" 2>/dev/null)
[ -z "$TOKEN" ] && { echo "Token not created yet, check $AXERONDIR/logs/axbridge.log"; exit 1; }

echo "AxBridge running on http://127.0.0.1:$PORT"
echo "Token file: $DATA/token"
am start -a android.intent.action.VIEW -d "http://127.0.0.1:$PORT/#token=$TOKEN" >/dev/null 2>&1 \
  && echo "Opened in browser." \
  || echo "Could not open a browser; open the URL manually."
