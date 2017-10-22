#!/bin/bash

resultDir="$(pwd)/benchResult"
tmpDir="/tmp/toy_kvs"
mkdir $resultDir
git clone https://github.com/kriw/toy_kvs $tmpDir
dir=("harddisk01" "nvme" "spdk_fuse")
for d in "${dir[@]}"
do
    cd "/mnt/$d"
    sudo rm -rf toy_kvs
    sudo cp -r $tmpDir ./
    sudo chown kriw toy_kvs -R
    cd toy_kvs
    for i in $(seq 1 9)
    do
        rm files/* 2>/dev/null
        bash run.sh server &
        sleep 15
        n=$(echo "2 ^ $i" | bc)
        echo $n
        bash run.sh benchmark -client-num=1024 -client-parallel=$n >> $resultDir"/"$d"NonDirectIO"
        pkill -9 go
        pkill -9 main
    done
    for i in $(seq 1 9)
    do
        rm files/* 2>/dev/null
        DIRECT=true bash run.sh server &
        sleep 15
        n=$(echo "2 ^ $i" | bc)
        echo $n
        bash run.sh benchmark -client-num=1024 -client-parallel=$n >> $resultDir"/"$d"DirectIO"
        pkill -9 go
        pkill -9 main
    done
done
