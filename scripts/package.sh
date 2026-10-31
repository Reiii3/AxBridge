#!/bin/sh
# Builds dist/axbridge-<version>.zip (flashable in AxManager). Run `make zip`, which builds the binaries first.
set -e
cd "$(dirname "$0")/.."
VERSION=${1:?usage: package.sh <version>}
OUT="dist/axbridge-$VERSION.zip"
STAGE=$(mktemp -d)
trap 'rm -rf "$STAGE"' EXIT
cp -r plugin/. "$STAGE/"
rm -f "$STAGE/system/bin/.gitkeep"
mkdir -p "$STAGE/client" && cp client/axbridge.js "$STAGE/client/"
mkdir -p dist && rm -f "$OUT"
(cd "$STAGE" && zip -qr -X "$OLDPWD/$OUT" .)
echo "built $OUT"
