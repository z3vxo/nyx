#!/bin/bash
set -e

CERT_DIR="$HOME/.kronos/certs"
mkdir -p "$CERT_DIR"

openssl req -x509 -newkey rsa:4096 -days 365 -nodes \
  -keyout "$CERT_DIR/server.key" \
  -out "$CERT_DIR/cert.crt" \
  -subj "/CN=kronos"

echo "Generated:"
echo "  $CERT_DIR/cert.crt"
echo "  $CERT_DIR/server.key"
