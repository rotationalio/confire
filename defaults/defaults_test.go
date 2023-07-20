package defaults_test

import (
	"errors"
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
	NestedPtr       *Embedded
	NestedDefault   *EmbeddedDefault `default:"pear"`
	AfterNested     string           `default:"after"`
	NoDefaultString string
	NoDefaultNum    int64
	NoDefaultTime   time.Time
	NoDefaultNested struct {
		PropertyA string
		PropertyB string
	}
	NoDefaultDuration time.Duration
}

type Embedded struct {
	Enabled bool     `default:"true"`
	EmPort  int      `default:"443"`
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
		return nil
	}

	if v == "pear" {
		e.ValueA = "pear"
		e.ValueB = 109
		return nil
	}

	return errors.New("unknown embedded default")
}

func TestDefaults(t *testing.T) {
	spec := Specification{}
	err := defaults.Process(&spec)
	assert.Ok(t, err)

	assert.True(t, spec.Enabled)
	assert.Equals(t, 443, spec.EmPort)
	assert.Equals(t, []string{"en", "fr"}, spec.Langs)
	assert.Equals(t, "apple", spec.ValueA)
	assert.Equals(t, int64(42), spec.ValueB)
	assert.True(t, spec.Debug)
	assert.Equals(t, 8888, spec.Port)
	assert.Equals(t, float32(0.5), spec.Rate)
	assert.Equals(t, "admin", spec.User)
	assert.Equals(t, uint32(60), spec.TTL)
	assert.Equals(t, 5*time.Minute, spec.Timeout)
	assert.Equals(t, time.Date(2020, 2, 20, 14, 22, 22, 0, time.UTC), spec.Epoch)
	assert.Equals(t, []string{"apples", "pears", "oranges"}, spec.AdminUsers)
	assert.Equals(t, []int{2, 4, 6, 8}, spec.MagicNumbers)
	assert.Equals(t, map[string]int{"red": 1, "blue": 2, "green": 3}, spec.ColorCodes)
	assert.Equals(t, "inner", spec.NestedSpecification.PropertyA)
	assert.Equals(t, "fuzzybydefault", spec.NestedSpecification.PropertyB)
	assert.True(t, spec.NestedPtr.Enabled)
	assert.Equals(t, 443, spec.NestedPtr.EmPort)
	assert.Equals(t, []string{"en", "fr"}, spec.NestedPtr.Langs)
	assert.Equals(t, "pear", spec.NestedDefault.ValueA)
	assert.Equals(t, int64(109), spec.NestedDefault.ValueB)
	assert.Equals(t, "after", spec.AfterNested)
	assert.Equals(t, "", spec.NoDefaultString)
	assert.Equals(t, int64(0), spec.NoDefaultNum)
	assert.Equals(t, time.Time{}, spec.NoDefaultTime)
	assert.Equals(t, "", spec.NoDefaultNested.PropertyA)
	assert.Equals(t, "", spec.NoDefaultNested.PropertyB)
	assert.Equals(t, time.Duration(0), spec.NoDefaultDuration)

}
