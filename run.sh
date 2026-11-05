#!/usr/bin/env bash

set -x

cd "$(dirname "${0}")"

PASSWORD=secret ./bin/webdav-server -l 8000 -u alice -s ~

