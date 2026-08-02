#!/bin/bash
set -e

# macOS Code Signing & Gatekeeper Notarization Helper Script
# This script strips quarantine attributes and signs the macOS app bundle so Gatekeeper recognizes it as a trusted application.

APP_PATH="${1:-app/vpsm-desktop/build/bin/VPSM Desktop.app}"

if [ ! -d "$APP_PATH" ]; then
    echo "[!] App bundle not found at '$APP_PATH'. Build first with 'make build-desktop'."
    exit 1
fi

echo "[info] Clearing quarantine attributes from '$APP_PATH'..."
xattr -cr "$APP_PATH"

echo "[info] Signing app bundle with hardened runtime..."
if [ -n "$APPLE_SIGNING_IDENTITY" ]; then
    echo "[info] Using signing identity: $APPLE_SIGNING_IDENTITY"
    codesign --force --deep --options runtime --timestamp -s "$APPLE_SIGNING_IDENTITY" "$APP_PATH"
else
    echo "[info] No APPLE_SIGNING_IDENTITY found. Applying ad-hoc signature with hardened runtime..."
    codesign --force --deep --options runtime -s - "$APP_PATH"
fi

echo "[info] Verification details:"
codesign -dv --verbose=2 "$APP_PATH"
