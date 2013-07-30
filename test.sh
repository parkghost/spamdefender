#!/bin/bash
# a helper script to run tests in the appropriate directories

set -e
source set_go_path.sh

echo "testing pkg"
pushd pkg >/dev/null
go test -v "./..."
popd >/dev/null

echo "testing jworld"
pushd jworld >/dev/null
go test -v "./..."
popd >/dev/null

# no tests, but a build is something
echo "build server"
pushd server/src >/dev/null
go build
go clean
popd >/dev/null

echo "build tools"
pushd tools >/dev/null
go build testing.go 
go build training.go 
go build builddictionarydata.go
go clean

echo "build tools/mailfetcher"
pushd mailfetecher >/dev/null
go build fetcher.go
go clean
popd >/dev/null

echo "build tools/mailclassifier"
pushd mailclassifier >/dev/null
go build classifier.go
go build explain.go
go build termfreq.go
go build training.go
go clean
popd >/dev/null

echo "build tools/mailsender"
pushd mailsender >/dev/null
go build sender.go
go clean
popd >/dev/null
popd >/dev/null 
