#!/usr/bin/env bash

set -euo pipefail

# Generate a self-signed TLS certificate and private key

# Check if -n flag is provided for non-interactive mode
if [[ $# -gt 0 && "$1" == "-n" ]] ; then
    echo "Generating certificate in non-interactive mode..."
    openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes \
      -subj "/C=TN/ST=Gafsa/L=Gafsa_Southern/O=WebdavServer/CN=localhost"
    echo "✓ Certificate and private key generated successfully"
    echo "  - server.crt (certificate)"
    echo "  - server.key (private key)"
else
    echo "Generating certificate in interactive mode..."
    echo "Follow the prompts to enter certificate details:"
    openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes
    echo "✓ Certificate and private key generated successfully"
    echo "  - server.crt (certificate)"
    echo "  - server.key (private key)"
fi

