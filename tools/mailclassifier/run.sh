#!/bin/bash
./clean.sh
go run training.go
go run classifier.go
