package structs_test

import (
	"testing"
	"time"

	"github.com/rotationalio/confire/assert"
	"github.com/rotationalio/confire/errors"
	"github.com/rotationalio/confire/structs"
)

func TestNewIsPointer(t *testing.T) {
	// Test struct initialization and IsPointer method
	spec := &Specification{}
	s, err := structs.New(spec)
	assert.Ok(t, err)
	assert.Assert(t, s.IsPointer(), "expected specification to be a pointer")

	val := Specification{}
	s, err = structs.New(val)
	assert.Ok(t, err)
	assert.Assert(t, !s.IsPointer(), "expected specification value to not be a pointer")
}

func TestNewError(t *testing.T) {
	var (
		foos  string   = ""
		foosp *string  = &foos
		fooi  int      = 42
		fooip *int     = &fooi
		foof  float64  = 3.14
		foofp *float64 = &foof
	)

	testCases := []interface{}{
		foos, foosp,
		fooi, fooip,
		foof, foofp,
		nil,
	}

	for _, tc := range testCases {
		_, err := structs.New(tc)
		assert.ErrorIs(t, err, errors.ErrNotAStruct)
	}
}

func TestName(t *testing.T) {
	testCases := []struct {
		spec interface{}
		name string
	}{
		{&Specification{}, "Specification"},
		{Specification{}, "Specification"},
		{&Nested{}, "Nested"},
		{Color{}, "Color"},
	}

	for _, tc := range testCases {
		s, err := structs.New(tc.spec)
		assert.Ok(t, err)
		assert.Equals(t, tc.name, s.Name())
	}
}

func TestNames(t *testing.T) {
	s, err := structs.New(&Specification{})
	assert.Ok(t, err)

	fields := s.Names()
	expected := []string{"Debug", "Name", "Port", "Rate", "Users", "Authn", "Timeout", "Started", "Colors", "OptionalColors", "Nested", "OptionalNested", "Level", "unexported"}
	assert.Equals(t, 14, len(fields))
	assert.Equals(t, expected, fields)
}

func TestFields(t *testing.T) {
	s, err := structs.New(&Specification{})
	assert.Ok(t, err)

	fields := s.Fields()
	assert.Equals(t, 14, len(fields))
}

func TestField(t *testing.T) {
	s, err := structs.New(&Specification{})
	assert.Ok(t, err)

	// Test a field that does exist
	_, err = s.Field("Debug")
	assert.Ok(t, err)

	// Test a field that does not exist
	_, err = s.Field("Foo")
	assert.Assert(t, err.Error() == "no field named \"Foo\" on Specification", "expected field Foo to return an error")
}

func TestIsZero(t *testing.T) {
	spec := &Specification{}
	s, err := structs.New(spec)
	assert.Ok(t, err)
	assert.Assert(t, s.IsZero(), "expected specification pointer to be zero-valued")

	val := Specification{}
	s, err = structs.New(val)
	assert.Ok(t, err)
	assert.Assert(t, s.IsZero(), "expected specification value to be zero-valued")

	// unexported fields should not be evaluated for zero-valuedness
	spec.unexported = "foo"
	assert.Assert(t, s.IsZero(), "expected specification with non-zero unexported field to be zero-valued")
}

func TestNotIsZero(t *testing.T) {
	testCases := []SpecMod{
		func(s *Specification) { s.Debug = true },
		func(s *Specification) { s.Name = "sultan" },
		func(s *Specification) { s.Port = 2121 },
		func(s *Specification) { s.Users = []string{"a", "b", "c"} },
		func(s *Specification) { s.Authn = map[string]string{"a": "b"} },
		func(s *Specification) { s.Timeout = 1500 * time.Millisecond },
		func(s *Specification) { s.Started = time.Now() },
		func(s *Specification) { s.Colors = []Color{{"Pink", 148, 29, 19}} },
		func(s *Specification) { s.OptionalColors = []*Color{nil, nil} },
		func(s *Specification) { s.Nested = Nested{Name: "foo"} },
		func(s *Specification) { s.OptionalNested = &Nested{Peers: 42} },
		func(s *Specification) { s.Level = LogLevel(16) },
	}

	for i, tc := range testCases {
		spec := &Specification{}
		tc(spec)

		s, err := structs.New(spec)
		assert.Ok(t, err)
		assert.Assert(t, !s.IsZero(), "test case %d failed: expected struct to not be zero-valued", i)
	}
}

func TestNotHasZero(t *testing.T) {
	spec := NewCompleteSpec()
	s, err := structs.New(spec)
	assert.Ok(t, err)
	assert.Assert(t, !s.HasZero(), "expected specification pointer to not have any zero-valued fields")

	s, err = structs.New(*spec)
	assert.Ok(t, err)
	assert.Assert(t, !s.HasZero(), "expected specification value to not have any zero-valued fields")

	// unexported fields should not be evaluated for zero-valuedness
	spec.unexported = ""
	assert.Assert(t, !s.HasZero(), "expected specification with non-zero unexported field to not have any zero-valued fields")
}

func TestHasZero(t *testing.T) {
	testCases := []SpecMod{
		func(s *Specification) { s.Debug = false },
		func(s *Specification) { s.Name = "" },
		func(s *Specification) { s.Port = 0 },
		func(s *Specification) { s.Users = nil },
		func(s *Specification) { s.Authn = nil },
		func(s *Specification) { s.Timeout = 0 },
		func(s *Specification) { s.Started = time.Time{} },
		func(s *Specification) { s.Colors = nil },
		func(s *Specification) { s.OptionalColors = nil },
		func(s *Specification) { s.Nested = Nested{} },
		func(s *Specification) { s.OptionalNested = nil },
		func(s *Specification) { s.Level = LogLevel(254) },
	}

	for i, tc := range testCases {
		spec := NewCompleteSpec()
		tc(spec)

		s, err := structs.New(spec)
		assert.Ok(t, err)
		assert.Assert(t, s.HasZero(), "test case %d failed: expected struct to have one zero-valued field", i)
	}
}

/*
Fixture Helpers
*/

type Specification struct {
	Debug          bool
	Name           string
	Port           int
	Rate           float64
	Users          []string
	Authn          map[string]string
	Timeout        time.Duration
	Started        time.Time
	Colors         []Color
	OptionalColors []*Color
	Nested         Nested
	OptionalNested *Nested
	Level          LogLevel
	unexported     string
}

type Color struct {
	Name  string
	Red   uint8
	Blue  uint8
	Green uint8
}

type Nested struct {
	Enabled    bool
	Name       string
	Peers      int
	Chance     float64
	Interval   time.Duration
	Deadline   time.Time
	Acks       []string
	Counts     map[string]int64
	unexported string
}

type LogLevel uint8

func (l LogLevel) IsZero() bool {
	if uint8(l) == 0 || uint8(l) >= 192 {
		return true
	}
	return false
}

func NewCompleteSpec() *Specification {
	return &Specification{
		Debug: true,
		Name:  "apples",
		Port:  5356,
		Rate:  0.44567,
		Users: []string{"Alice", "Bob", "Charlie", "Judy"},
		Authn: map[string]string{
			"Alice":   "admin",
			"Bob":     "member",
			"Charlie": "observer",
			"Judy":    "member",
		},
		Timeout: 30 * time.Second,
		Started: time.Date(2023, 4, 7, 19, 32, 21, 539212, time.UTC),
		Colors: []Color{
			{},
		},
		OptionalColors: []*Color{},
		Nested: Nested{
			Enabled:  true,
			Name:     "bananas",
			Peers:    14,
			Chance:   0.85,
			Interval: 92 * time.Minute,
			Deadline: time.Date(2023, 8, 31, 8, 14, 21, 92912, time.UTC),
			Acks:     []string{"baker", "lima", "foxtrot"},
			Counts: map[string]int64{
				"baker":   12,
				"lima":    93,
				"foxtrot": 49,
			},
			unexported: "mechanism",
		},
		OptionalNested: &Nested{
			Enabled:  true,
			Name:     "bananas",
			Peers:    14,
			Chance:   0.85,
			Interval: 92 * time.Minute,
			Deadline: time.Date(2023, 8, 31, 8, 14, 21, 92912, time.UTC),
			Acks:     []string{"baker", "lima", "foxtrot"},
			Counts: map[string]int64{
				"baker":   12,
				"lima":    93,
				"foxtrot": 49,
			},
			unexported: "mechanism",
		},
		Level:      LogLevel(8),
		unexported: "hello world!",
	}
}

func (s *Specification) Clone() *Specification {
	spec := &Specification{
		Debug:          s.Debug,
		Name:           s.Name,
		Port:           s.Port,
		Rate:           s.Rate,
		Users:          make([]string, 0, len(s.Users)),
		Authn:          make(map[string]string),
		Timeout:        s.Timeout,
		Started:        s.Started,
		Colors:         make([]Color, 0, len(s.Colors)),
		OptionalColors: s.OptionalColors,
		Nested: Nested{
			Enabled:    s.Nested.Enabled,
			Name:       s.Nested.Name,
			Peers:      s.Nested.Peers,
			Chance:     s.Nested.Chance,
			Interval:   s.Nested.Interval,
			Deadline:   s.Nested.Deadline,
			Acks:       make([]string, 0, len(s.Nested.Acks)),
			Counts:     make(map[string]int64),
			unexported: s.Nested.unexported,
		},
		OptionalNested: s.OptionalNested,
		Level:          s.Level,
		unexported:     s.unexported,
	}

	spec.Users = append(spec.Users, s.Users...)
	spec.Colors = append(spec.Colors, s.Colors...)
	spec.Nested.Acks = append(spec.Nested.Acks, s.Nested.Acks...)

	for key, val := range s.Authn {
		spec.Authn[key] = val
	}

	for key, val := range s.Nested.Counts {
		spec.Nested.Counts[key] = val
	}

	return spec
}

type SpecMod func(s *Specification)
