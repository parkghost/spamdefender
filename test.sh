#!/bin/bash
set -e
# a helper script to run tests in the appropriate directories

for dir in html ; do
    echo "testing $dir"
    pushd $dir >/dev/null
    go test -test.v -timeout 15s
    popd >/dev/null
done

echo "testing mailfile"
pushd mailfile >/dev/null
go test -v -test.run=".*POP3|.*RFC2047"
popd >/dev/null 

# no tests, but a build is something
echo "build spamdefender"
go build
go clean

echo "build tools"
pushd tools >/dev/null
go build testing.go 
go build training.go 
go clean

echo "build mailfetcher"
pushd mailfetecher >/dev/null
go build fetcher.go
go clean
popd >/dev/null

echo "build mailautoclassifier"
pushd mailautoclassifier >/dev/null
go build classifier.go
go build explain.go
go build termfreq.go
go build training.go
go clean
popd >/dev/null
popd >/dev/null 