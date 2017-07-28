#!/bin/sh
if [ $# -le 0 ] ; then
    cmd='none'
else
    cmd=$1
fi

if test $cmd == "server"; then
    arg=$2
    if [ "$arg" == "" ]; then
        arg='/tmp/tmp.sock'
    fi
    [[ -e $arg ]] && rm $arg
    go run server/main.go $arg
    [[ -e $arg ]] && rm $arg

elif test $cmd == "client"; then
    go run client/main.go
else
    echo "Usage $0 [client|server]"
fi
