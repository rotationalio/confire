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
	Port    int            `required:"true"`
	Level   string         `default:"info"`
	Rate    float64        `default:"1.0"`
	Timeout time.Duration  `desc:"read timeout"`
	Colors  map[string]int `desc:"at least three colors required"`
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
	ManualOverride   string `env:"manual_override" desc:"only set if you're sure"`
	EnvconfigCompat  bool   `envconfig:"MY_ENVCONFIG_VAR"`
	DefaultVar       string `default:"foo"`
	RequiredVar      string `required:"true" desc:"set anything here"`
	IgnoredVar       string `ignored:"true"`
	AutoSplitVar     string `split_words:"true" default:"bar"`
	ValidatedVar     string `validate:"required"`
	IgnoreValidation string `validate:"ignore"`
}
```

Generally speaking, confire will look for an environment variable in the form of `PREFIX_VARNAME`, where `prefix` is specified to the `confire.Process` function. The environment variable can be specified in two ways:

1. Specifying an alternate using the `env` or `envconfig` struct tags.
2. Specifying `split_words:"true"` to convert CamelCase to UPPER_UNDERSCORE case.

The `default` struct tag will process the given string as the default value before loading it from the environment.

The `required` and `validate` struct tags allow users to specify validation mechanism for the field. And the `ignored` tag will ensure the value is not processed from the environment or validated. If you want the value to be processed by the environment but not validated, then use the `validate:"ignore"` struct tag.

Finally the `desc` tag is used for documentation purposes and helps users understand what the variable is for. This is also printed out using the `usage.Usage` function described below.

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

Confire automatically looks for an environment variable to set on your configuration struct based on the name of the struct variable. Consider the following go code:

```go
import "github.com/rotationalio/confire/env"

type Config struct {
	Enabled    bool
	BindAddr   string
	ValidLangs []string
}

func main() {
	var conf Config
	env.Process("myapp", &conf)
}
```

The `env.Process` and `confire.Process` functions will go through all of the fields in the struct and lookup environment variables based on the name of the field, upper cased and prefixed with the string passed into the `Process` method. For example, the environment variables used will be:

- `$MYAPP_ENABLED`
- `$MYAPP_BINDADDR`
- `$MYAPP_VALIDLANGS`

You can modify the name of the environment variable in two ways. First, if you want to convert a CamelCase variable name to a UPPER_SNAKE case environment variable, you can specify the `split_words` struct tag:

```go
type Config struct {
	BindAddr string `split_words:"true"`
	TCPHosts string `split_words:"true"`
}
```

This will cause the environment variable to become `$MYAPP_BIND_ADDR` for the `BindAddr` variable. The library does it's best to preserve acryonyms so the `TCPHosts` variable will be looked up using `$MYAPP_TCP_HOSTS`.

You can also specify manual overrides for the environment variable which will provide an alternate lookup:

```go
type Config struct {
	AWSClientID string `env:"aws_client_id"`
	AWSSecret   string `envconfig:"aws_client_secret"`
}
```

In this case, confire will first lookup `$MYAPP_AWS_CLIENT_ID` then `$AWS_CLIENT_ID` and `$MYAPP_AWS_CLIENT_SECRET` and `$AWS_CLIENT_SECRET` in that order. Note that the `envconfig` tag is specified for compatibility with the `github.com/kelseyhightower/envconfig` library.

If you would like a single field in your config to not be processed by the `env` library then set the `ignored` tag as follows:

```go
type Config struct {
	SuperSecret string `ignored:"true"`
}
```

This field will be neither loaded from the environment nor validated.

If you do not want confire to process environment variables, use the `NoEnv` option as follows:

```go
confire.Process(&conf, confire.NoEnv)
```

## Usage

Configuration structs get big fast, and it can be a real pain to manage them. To provide some assistance, confire provides a method for printing out the environment variables, types, required validation, and default values from your struct tags:

```go
import "github.com/rotationalio/confire/usage"

type Config struct {
	Debug   bool
	Port    int            `required:"true"`
	Level   string         `default:"info"`
	Rate    float64        `default:"1.0"`
	Timeout time.Duration  `desc:"read timeout"`
	Colors  map[string]int `desc:"at least three colors required"`
	Peers   []string
}

func main() {
	var conf Config
	usage.Usage("myapp", &conf)
}
```

This will print out:

```
This application is configured via the environment. The following environment
variables can be used:

KEY              TYPE                                            DEFAULT    REQUIRED    DESCRIPTION
MYAPP_DEBUG      True or False
MYAPP_PORT       Integer                                                    true
MYAPP_LEVEL      String                                          info
MYAPP_RATE       Float                                           1.0
MYAPP_TIMEOUT    Duration                                                               read timeout
MYAPP_COLORS     Comma-separated list of String:Integer pairs                           at least three colors required
MYAPP_PEERS      Comma-separated list of String
```

The `usage.Usage` command does its best to determine the environment variable, but will always use the priority variable rather than the alternate variable.

Use the `desc` tag to provide a description and help document your code!

You can also print out a list format instead of the table format using:

```go
usage.Usagef("myapp", &conf, os.Stdout, usage.DefaultListFormat)
```

Which outputs:

```
This application is configured via the environment. The following environment
variables can be used:

MYAPP_DEBUG
  [description]
  [type]        True or False
  [default]
  [required]
MYAPP_PORT
  [description]
  [type]        Integer
  [default]
  [required]    true
MYAPP_LEVEL
  [description]
  [type]        String
  [default]     info
  [required]
MYAPP_RATE
  [description]
  [type]        Float
  [default]     1.0
  [required]
MYAPP_TIMEOUT
  [description] read timeout
  [type]        Duration
  [default]
  [required]
MYAPP_COLORS
  [description] at least three colors required
  [type]        Comma-separated list of String:Integer pairs
  [default]
  [required]
MYAPP_PEERS
  [description]
  [type]        Comma-separated list of String
  [default]
  [required]
```

You can pass your own custom format string in using `Usagef` or a template using `Usaget`. See the documentation for more information about what variables are available.

## Validation

Fields and structs can be automatically validated after processing by confire or by using the `validate.Validate` command. Validation occurs three ways:

1. Checking that the field isn't zero-valued using the `required` tag
2. Calling the `Validate` method of a field that implements the `Validator` interface
3. Validating the field using a built-in validator specified by the `validate` tag

All three methods can be used in the above order to perform validation and all methods specified by the struct tag must pass in order for the validation to pass.

The required tag is pretty straight forward:

```go
type Config struct {
	BindAddr string `required:"true"`
}
```

This ensures that `conf.BindAddr` cannot be an empty string (`""`).

The `Validator` interface is:

```go
type Validator interface {
	Validate() error
}
```

If the field implements this interface, the `Validate()` method is called and any error that is returned is converted into an `errors.ValidationError` from the confire error package.

Finally built-in validators can be used using the `validate` tag:

```go
type Config struct {
	BindAddr string `validate:"required"`
}
```

This will ensure that the "required" built-in validator is used. Current built-in validators are:

- `required`: ensure the field isn't zero-valued
- `ignore`: skip validation
- More coming soon!

You can ignore validation on any field by specifying the `validate:"ignore"` tag, this will prevent validation but still load the variable from the environment. You can also use the `ignored:"true"` tag, which will skip both environment loading and validation.

If you do not want confire to perform any validation at all, use the `NoValidate` option as follows:

```go
confire.Process(&conf, confire.NoValidate)
```

## Parsing

Environment variables and default values in struct tags are all strings that must be parsed into more complex types such as `bool`, `uint64`, `[]string`, `map[int]string` and others, therefore some parsing is required.

Default types such as `bool`, `int`, `uint`, `float`, and their bit-variants are parsed using the `strconv` library. Therefore you should use `true` and `false` for bools, and decimal integer representations without separators for numbers.

The `time.Duration` type is specifically handled using `time.ParseDuration` so you should pass in a duration string such as `"5s"` for 5 seconds or `3h2m10ms` for 3 hours, 2 minutes, 10 milliseconds.

Slices are parsed as comma-separated values of handled types. For example, a `[]time.Duration` type needs to be `"5s,10s,1m,1m30s"` which will result in a duration slice of length 4. There is no escaping or advanced handling for these values, so care is needed, particularly for `[]string`.

Byte slices, `[]byte`, must be represented by base64 encoded strings and are decoded as base64 arrays.

Maps are parsed by comma-separated key value pairs where the keys and values should be handled types. For example, a `map[string]uint64` should be represented as `alpha:32,bravo:41,charlie:51` to create a map with length 3. Again, there is no escaping or complex validation of these strings.

Finally, the `encoding.TextUnmarshaler` and `encoding.BinaryUnmarshaler` are also respected for parsing, which means other built-in types such as `time.Time` work using its `time.TextUnmarshal` method.

For more advanced parsing, use the `Decoder` or `Setter` interfaces as described below.

### Decoder Interface

The `Decoder` interface takes precedence over all other parsing methods and is defined as:

```go
type Decoder interface {
	Decode(value string) error
}
```

An example of using the `Decoder` interface is as follows:

```go
type Color [3]uint8

// Decode converts a hex color string such as #cc6699 into an RGB byte array.
func (c *Color) Decode(v string) error {
	// Strip a leading #
	if strings.HasPrefix(v, "#") {
		v = v[1:]
	}

	n, err := hex.Decode(c[:], v)
	if err != nil {
		return err
	}

	if n != 3 {
		return bytes.ErrTooLarge
	}
	return nil
}
```

### Setter Interface

Confire will use the `Setter` interface like from the `flag.Value` interface if implemented, though `Decoder` will take precedence.

```go
type Setter interface {
	Set(value string) error
}
```

An example of using the `Setter` interface with an enumeration is as follows:

```go
type LogLevel uint8

const (
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
	LevelPanic
)

func (ll *LogLevel) Set(v string) error {
	v = strings.TrimSpace(strings.ToLower(v))
	switch v {
	case "trace":
		*ll = LevelTrace
	case "debug":
		*ll = LevelDebug
	case "info":
		*ll = LevelInfo
	case "warning", "warn":
		*ll = LevelWarning
	case "error":
		*ll = LevelError
	case "fatal":
		*ll = LevelFatal
	case "panic":
		*ll = LevelPanic
	default:
		return fmt.Errorf("unknown level %q", v)
	}
	return nil
}
```

## Structs

This package makes use of reflection and you might want to use it's reflection in your code as well. We've ported and adapted the `github.com/fatih/structs` package into the confire library to make this a bit simpler. Please see the code documentation for more detail about the available methods. The basic way to loop through all the fields of a struct is as follows:

```go
func main() {
	var s *structs.Struct
	if s, err = structs.New(&conf); err != nil {
		return errors.ErrInvalidSpecification
	}

	if !s.IsPointer() {
		return errors.ErrInvalidSpecification
	}

	for _, field := range s.Fields() {
		// Use field.Kind() to recurse into nested structs.
	}
}
```

Obviously this example is missing a lot of detail, but you can refer to the code in the `defaults`, `validate`, and `env` package to see how they iterate through the fields in a `struct` and fetch tags and perform both read-only and modifying operations.

## Merging and Patching

Coming soon!

## Credits

Special thanks to the following libraries for providing inspiration and code snippets using their open sources licenses:

- [github.com/koding/mulitconfig](https://github.com/koding/multiconfig)
- [github.com/kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig)
- [github.com/aglyzov/go-patch](https://github.com/aglyzov/go-patch)
- [github.com/fatih/structs](https://github.com/fatih/structs/)

This package makes detailed use of the `reflect` package in Go and a lot of reference code was necessary to make this happen easily!
