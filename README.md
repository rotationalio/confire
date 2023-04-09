# Confire

**Configuration management for services and distributed systems**

## Install

```
$ go get github.com/rotationalio/confire
```

## Usage

Define a configuration struct in your code and load it with confire:

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rotationalio/confire"
)

type Config struct {
	Debug   bool
	Port    int
	Timeout time.Duration
	Rate    float64
}

func main() {
	conf := &Config{}
	if err = confire.Process("myapp", conf); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", conf)
}
```

## Credits

Special thanks to the following libraries for providing inspiration and code snippets using their open sources licenses:

- [github.com/koding/mulitconfig](https://github.com/koding/multiconfig)
- [github.com/kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig)
- [github.com/aglyzov/go-patch](https://github.com/aglyzov/go-patch)
- [github.com/fatih/structs](https://github.com/fatih/structs/)

This package makes detailed use of the `reflect` package in Go and a lot of reference code was necessary to make this happen easily!