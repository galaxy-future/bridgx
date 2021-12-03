#!/bin/bash
RUN_NAME="gf.bridgx.scheduler"
mkdir -p output/conf output/bin
sudo mkdir -p /tmp/bridgx/logs/

find conf/ -type f ! -name "*local*" -print0 | xargs -0 -I{} cp {} output/conf/
cp scripts/run_scheduler.sh output/

go fmt ./...
go vet ./...
export GO111MODULE="on"
export GOPRIVATE="code.galaxy-future.com"

go build -o output/bin/${RUN_NAME} ./cmd/scheduler