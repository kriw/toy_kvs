#!/bin/sh
if [ $# -le 0 ] ; then
    cmd='none'
else
    cmd=$1
fi

if test $cmd == "server"; then
    arg=$2
    if [ "$arg" == "" ]; then
        arg='127.0.0.1:8000'
    fi
    go run server/main.go $arg

elif test $cmd == "client"; then
    arg=$2
    if [ "$arg" == "" ]; then
        arg='127.0.0.1:8000'
    fi
    go run client/main.go $arg
elif test $cmd == "benchmark"; then
    go run server/main.go '127.0.0.1:8000' &
    srvId=$!
    go run benchmark/main.go ${@:2}
    kill -9 $srvId
else
    echo "Usage $0 [client|server|benchmark]"
fi
