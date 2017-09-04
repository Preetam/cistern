#!/bin/sh

set -e
cd ~/.go_workspace/src/github.com/Cistern/cistern
go build ./cmd/cistern -O cistern-linux-amd64 && mv cistern-linux-amd64 $CIRCLE_ARTIFACTS
go build ./cmd/cistern -O cistern-darwin-amd64 && mv cistern-darwin-amd64 $CIRCLE_ARTIFACTS
cd ui
npm i
npm run build
tar czvf cistern-ui-assets.tar.gz static && mv cistern-ui-assets.tar.gz $CIRCLE_ARTIFACTS
