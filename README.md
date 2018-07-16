# urlshorter

[![Build Status](https://travis-ci.org/brucemaclin/urlshorter.svg?branch=master)](https://travis-ci.org/brucemaclin/urlshorter) [![Go Report Card](https://goreportcard.com/badge/github.com/brucemaclin/urlshorter)](https://goreportcard.com/report/github.com/brucemaclin/urlshorter)
[![godoc](https://godoc.org/github.com/brucemaclin/urlshorter?status.svg)](https://godoc.org/github.com/brucemaclin/urlshorter)

golang url shorter like t.co,t.cn and so on.

Use uint64 id and convert it to 62 decimal string .

You should save origURL ,id and the 62 decimal string as shortURL by yourself.

# install

```
go get -u -v github.com/brucemaclin/urlshorter
```

# use

you can use without DB,but only use shortURLGene and GetID by short URL.

example:

```
    package main
    import (
        "github.com/brucemaclin/urlshorter"
        "math/rand"
        "fmt"
    )

    func main() {
        id := rand.Uint64()
        shortURL := urlshorter.ShorterURLGene(id)
        fmt.Println("shortURL for id:", shortURL)

    }
```


if use with db,you should implement the urlshorter.DB interface.

see DefaultDB in [db.go](db.go) .

you can use mysql/redis/mongo as your wish.

you can try with demo.

```
    go run demo.go
```

then you can type http://127.0.0.1:8080/?url=www.google.com,it will return
a shorter url.
