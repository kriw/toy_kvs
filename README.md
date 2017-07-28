# toy kvs

Simple key value store for homework of Security-Camp2017.

# Usage

## Server
```bash
./run.sh server
```

## Client
```shell
./run.sh client
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
