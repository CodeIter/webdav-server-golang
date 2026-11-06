#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "${0}")"

export PASSWORD="${PASSWORD:-secret password QWE+/601}"

# Default values
USERNAME="codeiter"
LISTEN="8000"
SERVE=~
TLS_MODE=false
CERT_FILE="server.crt"
KEY_FILE="server.key"

# Parse arguments
while [[ $# -gt 0 ]]; do
    case "$1" in
        -t|--tls)
            TLS_MODE=true
            shift
            ;;
        -c|--cert)
            CERT_FILE="$2"
            shift 2
            ;;
        -k|--key)
            KEY_FILE="$2"
            shift 2
            ;;
        -u|--username)
            USERNAME="$2"
            shift 2
            ;;
        -l|--listen)
            LISTEN="$2"
            shift 2
            ;;
        -s|--serve)
            SERVE="$2"
            shift 2
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Build the command
CMD=(
    "./bin/webdav-server"
    "--username" "${USERNAME}"
    "--listen" "${LISTEN}"
    "--serve" "${SERVE}"
)

if [[ "${TLS_MODE}" == true ]]; then
    CMD+=("--cert" "${CERT_FILE}" "--key" "${KEY_FILE}")
fi

echo
echo "Running webdav server with the following configuration:"
echo
echo "Listen   : ${LISTEN}"
echo "Serve    : ${SERVE}"
echo "TLS      : ${TLS_MODE}"
echo "CertFile : ${CERT_FILE}"
echo "KeyFile  : ${KEY_FILE}"
echo "Username : ${USERNAME}"
echo "Password : '${PASSWORD}'"
echo

set -x

"${CMD[@]}"
