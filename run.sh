#!/bin/bash
if [ $# -le 0 ] ; then
    cmd='none'
else
    cmd=$1
fi

if [ ! -e logFiles ] ; then
    mkdir logFiles
fi

if test $cmd == "server"; then
    rm files/*
    arg=$2
    if [ "$arg" == "" ]; then
        arg='127.0.0.1:8000'
    fi
    go run -race server/main.go $arg 

elif test $cmd == "client"; then
    arg=$2
    if [ "$arg" == "" ]; then
        arg='127.0.0.1:8000'
    fi
    go run client/main.go $arg
elif test $cmd == "benchmark"; then
    go run benchmark/main.go ${@:2}
else
    echo "Usage $0 [client|server|benchmark]"
fi
