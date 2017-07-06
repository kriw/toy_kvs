#!/bin/sh
if [ $# -le 0 ] ; then
    cmd='none'
else
    cmd=$1
fi

if test $cmd == "server" ; then
    [[ -e /tmp/echo.sock ]] && rm /tmp/echo.sock
    go run server/main.go
elif test $cmd == "client" ; then
    ./client/client.sh
else
    echo "Usage $0 [client|server]"
fi
