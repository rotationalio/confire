package merge_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/rotationalio/confire/merge"
	"github.com/rotationalio/confire/structs"
)

func TestMerge(t *testing.T) {
	src := &Color{
		Red:  0x7b,
		Blue: 0x5b,
	}

	dst := &Color{
		Red:   0xc9,
		Green: 0xa0,
		Blue:  0xdc,
	}

	changed, err := merge.Merge(dst, src)
	ok(t, err)
	assert(t, changed, "expected the destination to be changed")
	equals(t, uint8(0x7b), dst.Red)
	equals(t, uint8(0xa0), dst.Green)
	equals(t, uint8(0x5b), dst.Blue)
}

func TestMergeEmbedded(t *testing.T) {
	src := &Crayon{
		Name: "Thistle",
		Color: Color{
			Red:   0xd8,
			Green: 0xbf,
			Blue:  0xd8,
		},
	}

	dst := &Crayon{}

	changed, err := merge.Merge(dst, src)
	ok(t, err)
	assert(t, changed, "expected the destination to be changed")
	equals(t, "Thistle", dst.Name)
	equals(t, uint8(0xd8), dst.Red)
	equals(t, uint8(0xbf), dst.Green)
	equals(t, uint8(0xd8), dst.Blue)
}

func TestMergeNested(t *testing.T) {
	t.Skip("nested patching isn't quite working yet")
	src := &CrayonBox{
		Title:   "box alpha",
		Created: time.Now(),
		Highlight: &Crayon{
			Name: "Thistle",
			Color: Color{
				Red:  0xd8,
				Blue: 0xd8,
			},
		},
		unexported: "alpha",
	}

	dst := &CrayonBox{
		Created: time.Date(2023, 4, 7, 12, 12, 12, 12, time.UTC),
		Crayons: []*Crayon{
			{Color: Color{Red: 0x00, Green: 0x00, Blue: 0x00}, Name: ""},
		},
		Highlight: &Crayon{
			Color: Color{
				Green: 0xbf,
			},
		},
		unexported: "bravo",
	}

	changed, err := merge.Merge(dst, src)
	ok(t, err)
	assert(t, changed, "expected the destination to be changed")
	equals(t, "bravo", dst.unexported)
}

func TestIgnoreEmptyPatch(t *testing.T) {
	c := Crayon{
		Name: "Asparagus",
		Color: Color{
			Red:   0x7b,
			Green: 0xa0,
			Blue:  0x5b,
		},
	}

	changed, err := merge.Merge(&c, struct{}{})
	ok(t, err)
	assert(t, !changed, "the crayon should not have been changed")
	equals(t, "Asparagus", c.Name)
	equals(t, uint8(0x7b), c.Red)
	equals(t, uint8(0xa0), c.Green)
	equals(t, uint8(0x5b), c.Blue)

	changed, err = merge.Merge(c, struct{}{})
	ok(t, err)
	assert(t, !changed, "the crayon should not have been changed")
}

func TestOnlyStructs(t *testing.T) {
	_, err := merge.Merge(42, Crayon{})
	assert(t, errors.Is(err, structs.ErrNotAStruct), "expected error when dst is not a struct")

	_, err = merge.Merge(Crayon{}, 41)
	assert(t, errors.Is(err, structs.ErrNotAStruct), "expected error when src is not a struct")
}

/*
Fixture Helpers
*/

type Color struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

type Crayon struct {
	Color
	Name string
}

type CrayonBox struct {
	Title      string
	Created    time.Time
	Crayons    []*Crayon
	Highlight  *Crayon
	unexported string
}

/*
Assertion Helpers

Because this is a library, we prefer to have no dependencies including our usual test
dependencies (e.g. testify require). So we have some basic assertion helpers for tests.

See: https://github.com/benbjohnson/testing
*/

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	tb.Helper()
	if !condition {
		tb.Logf("\n"+msg+"\n", v...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Logf("\nunexpected error: %q\n", err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	tb.Helper()
	if !reflect.DeepEqual(exp, act) {
		tb.Logf("\nactual value did not match expected:\n\n\t- exp: %#v\n\t- got: %#v\n", exp, act)
		tb.FailNow()
	}
}
