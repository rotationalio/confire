# Confire

[![Go Reference](https://pkg.go.dev/badge/github.com/rotationalio/confire.svg)](https://pkg.go.dev/github.com/rotationalio/confire)
[![Tests](https://github.com/rotationalio/confire/actions/workflows/test.yaml/badge.svg)](https://github.com/rotationalio/confire/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rotationalio/confire)](https://goreportcard.com/report/github.com/rotationalio/confire)

**Configuration management for services and distributed systems**

## Install

```
$ go get github.com/rotationalio/confire
```

## Usage

Confire uses struct tags to understand how to load a configuration from specified default values and the environment (soon also from configuration files) and then validates the configuration on behalf of your application.

Basic usage is as follows. Define a configuration struct in your code and load it with confire:

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
	Port    int     `required:"true"`
	Level   string  `default:"info"`
	Rate    float64 `default:"1.0"`
	Timeout time.Duration
	Colors  map[string]int
	Peers   []string
}

func main() {
	conf := &Config{}
	if err = confire.Process("myapp", conf); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", conf)
}
```

Set some environment variables for configuration:

```
export MYAPP_DEBUG=true
export MYAPP_PORT=8888
export MYAPP_TIMEOUT="5s"
export MYAPP_RATE="0.25"
export MYAPP_COLORS="red:1,green:2,blue:3"
export MYAPP_PEERS="alpha,bravo,charlie"
```

Note that the environment variable is uppercase and is prefixed by the specified prefix passed to the `confire.Process` function and an underscore.

Output (slightly cleaned up for multiple lines and readability):

```
Debug:true
Port:8888
Level:info
Rate:0.25
Timeout:5s
Colors:map[blue:3 green:2 red:1]
Peers:[alpha bravo charlie]
```

[Try it yourself!](https://go.dev/play/p/PJ-Gw5C5Lrp)

### Advanced Usage

Confire uses struct tags to specify the environment variable to, fields to ignore, default values, and how to validate a field.

Consider the following struct:

```go
type Config struct {
	ManualOverride   string `env:"manual_override"`
	EnvconfigCompat  bool   `envconfig:"MY_ENVCONFIG_VAR"`
	DefaultVar       string `default:"foo"`
	RequiredVar      string `required:"true"`
	IgnoredVar       string `ignored:"true"`
	AutoSplitVar     string `split_words:"true"`
	ValidatedVar     string `validate:"required"`
	IgnoreValidation string `validate:"ignore"`
}
```

Generally speaking, confire will look for an environment variable in the form of `PREFIX_VARNAME`, where `prefix` is specified to the `confire.Process` function. The environment variable can be specified in two ways:

1. Specifying an alternate using the `env` or `envconfig` struct tags.
2. Specifying `split_words:"true"` to convert CamelCase to UPPER_UNDERSCORE case.

The `default` struct tag will process the given string as the default value before loading it from the environment.

The `required` and `validate` struct tags allow users to specify validation mechanism for the field. And the `ignored` tag will ensure the value is not processed from the environment or validated. If you want the value to be processed by the environment but not validated, then use the `validate:"ignore"` struct tag.

### Supported Field Types

Currently confire supports parsing these struct field types:

- string
- int8, int16, int32, int64
- uint8, uint16, uint32, uint64
- bool
- float32, float64
- [time.Duration](https://golang.org/pkg/time/#Duration)
- slices of any supported type
- maps (keys and values of any supported type)
- [encoding.TextUnmarshaler](https://golang.org/pkg/encoding/#TextUnmarshaler)
- [encoding.BinaryUnmarshaler](https://golang.org/pkg/encoding/#BinaryUnmarshaler)
- [parse.Decoder](https://pkg.go.dev/github.com/rotationalio/confire/parse#Decoder)
- [parse.Setter](https://pkg.go.dev/github.com/rotationalio/confire/parse#Setter)

Note that `time.Time` is also supported because it implements `encoding.TextUnmarshaler`.

## Defaults

The confire `defaults` package can be used to populate a struct from default values.

```go
import (
	"fmt"
	"log"

	"github.com/rotationalio/confire/defaults"
)

type Config struct {
	Enabled bool     `default:"true"`
	Port    int      `default:"443"`
	Langs   []string `default:"en,fr"`
}

func main() {
	var conf Config
	if err := defaults.Process(&conf); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", conf)
}
```

The `defaults` package uses the same parsing mechanism as the environment variables to convert struct tag strings into the specified type. The output of the above code will be a non-zero config that is populated with the specified values.

If you do not want confire to automatically process defaults, use the `NoDefaults` option as follows:

```go
confire.Process(&conf, confire.NoDefaults)
```

## Configuration Files

Coming soon!

## Environment Variables

## Usage

## Validation

## Parsing

## Structs

## Merging and Patching

Coming soon!

## Credits

Special thanks to the following libraries for providing inspiration and code snippets using their open sources licenses:

- [github.com/koding/mulitconfig](https://github.com/koding/multiconfig)
- [github.com/kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig)
- [github.com/aglyzov/go-patch](https://github.com/aglyzov/go-patch)
- [github.com/fatih/structs](https://github.com/fatih/structs/)

This package makes detailed use of the `reflect` package in Go and a lot of reference code was necessary to make this happen easily!