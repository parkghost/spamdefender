#!/bin/bash
set -e
# a helper script to run tests in the appropriate directories

for dir in html analyzer/goseg mailfile postfix ; do
    echo "testing $dir"
    pushd $dir >/dev/null
    go test -test.v -timeout 15s
    popd >/dev/null
done

# no tests, but a build is something
echo "build github.com/parkghost/spamdefender"
go build
go clean

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