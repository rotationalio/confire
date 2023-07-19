package parse_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/rotationalio/confire/assert"
	"github.com/rotationalio/confire/parse"
	"github.com/rotationalio/confire/structs"
)

func TestParse(t *testing.T) {
	var s string
	err := parse.Parse("super", reflect.ValueOf(&s))
	assert.Ok(t, err)
	assert.Equals(t, "super", s)

	var dur time.Duration
	err = parse.Parse("1h3m", reflect.ValueOf(&dur))
	assert.Ok(t, err)
	assert.Equals(t, time.Hour+(3*time.Minute), dur)
}

func TestParseField(t *testing.T) {
	spec := &Specification{}
	s, err := structs.New(spec)
	assert.Ok(t, err)

	debug := true
	name := "jonsey"
	ttl := uint32(132)
	port := 8888
	rate := float64(0.25)
	names := []string{"lad", "mack", "em", "bri"}

	testCases := map[string]struct {
		value    string
		expected any
	}{
		"Level":      {"warn", LevelWarning},
		"URI":        {"https://rotational.io/", Escape("https%3A%2F%2Frotational.io%2F")},
		"Timestamp":  {"2023-07-19T18:13:45.000Z", time.Date(2023, 7, 19, 18, 13, 45, 0, time.UTC)},
		"Color":      {"cc6699", Color{0xcc, 0x66, 0x99}},
		"Timeout":    {"1.892s", 1892 * time.Millisecond},
		"Debug":      {"true", debug},
		"DebugPtr":   {"true", &debug},
		"Name":       {"jonsey", name},
		"NamePtr":    {"jonsey", &name},
		"Names":      {"lad,mack,em,bri", names},
		"NamePtrs":   {"lad,mack,em,bri", []*string{&names[0], &names[1], &names[2], &names[3]}},
		"TTL":        {"132", ttl},
		"TTLPtr":     {"132", &ttl},
		"Port":       {"8888", port},
		"PortPtr":    {"8888", &port},
		"Rate":       {"0.25", rate},
		"RatePtr":    {"0.25", &rate},
		"Ages":       {"lad:4,mack:16,em:8,bri:14", map[string]int{"lad": 4, "mack": 16, "em": 8, "bri": 14}},
		"ColorNames": {"crimson:#dc143c,peach:#ffdab9,cadet:#5f9ea0", map[string]Color{"cadet": {0x5f, 0x9e, 0xa0}, "crimson": {0xdc, 0x14, 0x3c}, "peach": {0xff, 0xda, 0xb9}}},
		"Multi":      {"https://rotational.io,https://ensign.world", []Escape{"https%3A%2F%2Frotational.io", "https%3A%2F%2Fensign.world"}},
		"EmptyMap":   {"", map[string]int{}},
		"EmptySlice": {"", []string{}},
	}

	for _, field := range s.Fields() {
		tc := testCases[field.Name()]
		err := parse.ParseField(tc.value, field)
		assert.Ok(t, err)
		assert.Equals(t, tc.expected, field.Value())
	}
}

type Specification struct {
	Level      LogLevel      `desc:"log level implements Decoder"`
	URI        Escape        `desc:"escape implements Setter"`
	Timestamp  time.Time     `desc:"time.Time implements TextUnmarshaler"`
	Color      Color         `desc:"color implements BinaryUnmarshaler"`
	Timeout    time.Duration `desc:"duration is handled specially"`
	Debug      bool
	DebugPtr   *bool
	Name       string
	NamePtr    *string
	Names      []string
	NamePtrs   []*string
	TTL        uint32
	TTLPtr     *uint32
	Port       int
	PortPtr    *int
	Rate       float64
	RatePtr    *float64
	Ages       map[string]int
	ColorNames map[string]Color
	Multi      []Escape
	EmptyMap   map[string]int
	EmptySlice []string
}

type Color [3]uint8

func (c *Color) UnmarshalBinary(v []byte) error {
	// Strip a leading #
	if bytes.HasPrefix(v, []byte("#")) {
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

func (ll *LogLevel) Decode(v string) error {
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

type Escape string

func (e *Escape) Set(value string) error {
	*e = Escape(url.QueryEscape(value))
	return nil
}
