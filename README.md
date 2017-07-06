# toy kvs

Simple key value store for a homework of Security-Camp2017.

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
$ ./run.sh client
> set hoge piyo
OK
> get hoge
piyo
> get aaaaaa

> aaaaaa
Unknown query.
> set hoge po
OK
> get hoge
po
```
