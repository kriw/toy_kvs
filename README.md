# toy kvs

Simple key value store for a homework of Security-Camp2017.

# Usage

## Server
```shell
$ [[ -e /tmp/echo.sock ]] && rm /tmp/echo.sock
$ go run src/main.go
```

## Client
```shell
socat - unix-connect:/tmp/echo.sock
```


## Example

```shell
$ socat - unix-connect:/tmp/echo.sock
set hoge piyo
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
