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

Special thanks to [github.com/koding/mulitconfig](https://github.com/koding/multiconfig) and [github.com/kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig) for serving as inspiration for this package.