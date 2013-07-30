#!/bin/bash

source set_go_path.sh
./clean.sh

go run builddictdata.go
go run training.go
go run testing.go
