package defaults_test

import (
	"testing"
	"time"

	"github.com/rotationalio/confire/assert"
	"github.com/rotationalio/confire/defaults"
)

type Specification struct {
	Embedded
	EmbeddedDefault     `default:"apple"`
	Debug               bool           `default:"true"`
	Port                int            `default:"8888"`
	Rate                float32        `default:"0.5"`
	User                string         `default:"admin"`
	TTL                 uint32         `default:"60"`
	Timeout             time.Duration  `default:"5m"`
	Epoch               time.Time      `default:"2020-02-20T14:22:22.000Z"`
	AdminUsers          []string       `default:"apples,pears,oranges"`
	MagicNumbers        []int          `default:"2,4,6,8"`
	ColorCodes          map[string]int `default:"red:1,blue:2,green:3"`
	NestedSpecification struct {
		PropertyA string `default:"inner"`
		PropertyB string `default:"fuzzybydefault"`
	}
	// NestedPtr       *Embedded
	// NestedDefault   *EmbeddedDefault  `default:"pear"`
	AfterNested     string            `default:"after"`
	MapField        map[string]string `default:"one:two,three:four"`
	NoDefaultString string
	NoDefaultNum    int64
	NoDefaulttime   time.Time
	NoDefaultNested struct {
		PropertyA string
		PropertyB string
	}
	NoDefaultDuration time.Duration
}

type Embedded struct {
	Enabled bool     `default:"true"`
	Port    int      `default:"443"`
	Langs   []string `default:"en,fr"`
}

type EmbeddedDefault struct {
	ValueA string
	ValueB int64
}

func (e *EmbeddedDefault) Decode(v string) error {
	if v == "apple" {
		e.ValueA = "apple"
		e.ValueB = 42
	}

	if v == "pear" {
		e.ValueA = "pear"
		e.ValueB = 109
	}

	return nil
}

func TestDefaults(t *testing.T) {
	spec := Specification{}
	err := defaults.Process(&spec)
	assert.Ok(t, err)
}
