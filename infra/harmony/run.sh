#!/usr/bin/env bash
set -e

cd "$GOPATH/src/github.com/harmony-one/harmony"
bash ./test/deploy.sh -B -D 60000 ./test/configs/local-resharding.txt