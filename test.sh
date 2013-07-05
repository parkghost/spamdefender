#!/bin/bash
set -e
# a helper script to run tests in the appropriate directories

for dir in html mailfile; do
    echo "testing $dir"
    pushd $dir >/dev/null
    go test -test.v -timeout 15s
    popd >/dev/null
done

# no tests, but a build is something
go build
rm spamdefender

echo "build tools"
pushd tools >/dev/null
go build testing.go
rm testing
go build training.go
rm training

echo "build mailfetcher"
pushd mailfetecher >/dev/null
go build fetcher.go
rm fetcher
popd >/dev/null

echo "build mailautoclassifier"
pushd mailautoclassifier >/dev/null
go build classifier.go
rm classifier
go build explain.go
rm explain
go build termfreq.go
rm termfreq
go build training.go
rm training
popd >/dev/null

popd >/dev/null 