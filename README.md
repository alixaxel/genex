# genex [![GoDoc](https://godoc.org/github.com/alixaxel/genex?status.svg)](https://godoc.org/github.com/alixaxel/genex) [![GoCover](http://gocover.io/_badge/github.com/alixaxel/genex)](http://gocover.io/github.com/alixaxel/genex) [![Go Report Card](https://goreportcard.com/badge/github.com/alixaxel/genex)](https://goreportcard.com/report/github.com/alixaxel/genex)

Genex package for Go

Easy and efficient package to expand any given regex into all the possible strings that it can match.

This is the code that powers [namegrep](https://namegrep.com/).

## Usage

```go
package main

import (
    "fmt"
    "regexp/syntax"

    "github.com/alixaxel/genex"
)

func main() {
    charset, _ := syntax.Parse(`[0-9a-z]`, syntax.Perl)

    if input, err := syntax.Parse(`(foo|bar|baz){1,2}\d`, syntax.Perl); err == nil {
    	fmt.Println("Count:", genex.Count(input, charset, 3))

    	genex.Generate(input, charset, 3, func(output string) {
    		fmt.Println("[*]", output)
    	})
    }
}
```

## Output

```
Count: 120

[*] foo0
[*] ...
[*] foo9
[*] foofoo0
[*] ...
[*] foofoo9
[*] foobar0
[*] ...
[*] foobar9
[*] foobaz0
[*] ...
[*] foobaz9
[*] bar0
[*] ...
[*] bar9
[*] barfoo0
[*] ...
[*] barfoo9
[*] barbar0
[*] ...
[*] barbar9
[*] barbaz0
[*] ...
[*] barbaz9
[*] baz0
[*] ...
[*] baz9
[*] bazfoo0
[*] ...
[*] bazfoo9
[*] bazbar0
[*] ...
[*] bazbar9
[*] bazbaz0
[*] ...
[*] bazbaz9
```

## Install

	go get github.com/alixaxel/genex

## License

MIT
