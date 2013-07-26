#!/bin/bash

source set_go_path.sh

echo "build server"
pushd server/src >/dev/null
go build -o "../../bin/server"
popd >/dev/null