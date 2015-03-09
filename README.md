# genex

Genex package for Go

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
    		fmt.Println(output)
    	})
    }
}
```

## Install

	go get github.com/alixaxel/genex

## License

MIT
