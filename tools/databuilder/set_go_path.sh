#!/bin/bash
# usage: . set_go_path.sh

PARENT_PATH=`dirname $PWD`
ROOT_PATH=`dirname $PARENT_PATH`
GOPATH="$ROOT_PATH/vender:$ROOT_PATH/pkg:$ROOT_PATH/jworld:$ROOT_PATH/server" 
