#!/usr/bin/env bash

set -euo pipefail

echo "Generating certificate in interactive mode..."
echo "Follow the prompts to enter certificate details:"
echo

# Generate a self-signed TLS certificate and private key
openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes

echo
echo "✓ Certificate and private key generated successfully"
echo "  - server.crt (certificate)"
echo "  - server.key (private key)"

