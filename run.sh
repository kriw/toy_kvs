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
    arg=$2
    if [ "$arg" == "" ]; then
        arg='/tmp/tmp.sock'
    fi
    go run client/main.go $arg
elif test $cmd == "benchmark"; then
    go run benchmark/main.go ${@:2}
    [[ -e '/tmp/tmp.sock' ]] && rm /tmp/tmp.sock
else
    echo "Usage $0 [client|server]"
fi
