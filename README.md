TSV parser for Go
====

[![Build Status](http://img.shields.io/travis/dogenzaka/tsv.svg?style=flat)](https://travis-ci.org/dogenzaka/tsv)
[![Coverage](http://img.shields.io/codecov/c/github/dogenzaka/tsv.svg?style=flat)](https://codecov.io/github/dogenzaka/tsv)
[![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/dogenzaka/rotator/blob/master/LICENSE)

tsv is tab-separated values parser for GO. It will parse and data into any type of struct. tsv supports both simple struct and tagged fields.

```
go get github.com/dogenzaka/tsv
```

Quickstart
--

Struct is indexed by field order.

```go

import (
    "fmt"
    "os"
    "testing"
    )

type TestRow struct {
  Name   string
  Age    int
  Gender string
  Active bool
}

func main() {

  file, _ := os.Open("example.tsv")
  defer file.Close()

  data := TestRow{}
  parser, _ := NewParser(file, &data)

  for {
    eof, err := parser.Next()
    if eof {
      return
    }
    if err != nil {
      panic(err)
    }
    fmt.Println(data)
  }

}

```

You can define tags of fields in struct to map tsv values without ordering.

```go
type TestRow struct {
  Name   string `tsv:"name"`
  Age    int    `tsv:"age"`
  Gender string `tsv:"gender"`
  Active bool   `tsv:"bool"`
}
```

Supported field types
--

Currently this library supports

- int
- string
- bool

