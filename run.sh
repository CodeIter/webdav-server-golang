#!/usr/bin/env bash

set -x

cd "$(dirname "${0}")"

export PASSWORD="${PASSWORD:-secret password QWE+/601}"

if [[ $# -gt 0 && "$1" == "-t" ]] ; then
    ./bin/webdav-server --cert server.crt --key server.key --username codeiter --listen 8000 --serve ~
else
    ./bin/webdav-server --username codeiter --listen 8000 --serve ~
fi

