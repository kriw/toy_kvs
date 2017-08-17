# toy kvs

Simple key value store for homework of Security-Camp2017.

# Requirement

## yara
See https://github.com/hillu/go-yara .

## directIO
See https://github.com/ncw/directio

# Usage

## Server
```bash
./run.sh server [endpoint]
```
* The default endpoint is `/tmp/tmp.sock`.

With output, 
```bash
DEBUG=true ./run.sh server [endpoint]
```

## Client
```shell
./run.sh client
```

## Benchmark
```bash
./run.sh benchmark
client-num: 2, repeats: 5, elapsed: 17.904102ms

./run.sh benchmark -client-num=1 -repeats=3
client-num: 1, repeats: 3, elapsed: 8.942194ms
```

## Example

```shell
$ echo "This is a test." > po.txt
$ ./run.sh client
> set po.txt
> key: 11586d2eb43b73e539caa3d158c883336c0e2c904b309c0c5ffe2c9b83d562a1
> OK
> get 11586d2eb43b73e539caa3d158c883336c0e2c904b309c0c5ffe2c9b83d562a1
> This is a test.

> save poyo
> OK
>
$ rm po.txt
$ ls | grep poyo
poyo*
```
