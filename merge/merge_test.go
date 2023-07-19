package merge_test

import (
	"testing"
	"time"

	"github.com/rotationalio/confire/assert"
	"github.com/rotationalio/confire/errors"
	"github.com/rotationalio/confire/merge"
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
	assert.Ok(t, err)
	assert.Assert(t, changed, "expected the destination to be changed")
	assert.Equals(t, uint8(0x7b), dst.Red)
	assert.Equals(t, uint8(0xa0), dst.Green)
	assert.Equals(t, uint8(0x5b), dst.Blue)
}

func TestMergeEmbedded(t *testing.T) {
	t.Skip("embedded structs are not settable and cannot be merged?")
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
	assert.Ok(t, err)
	assert.Assert(t, changed, "expected the destination to be changed")
	assert.Equals(t, "Thistle", dst.Name)
	assert.Equals(t, uint8(0xd8), dst.Red)
	assert.Equals(t, uint8(0xbf), dst.Green)
	assert.Equals(t, uint8(0xd8), dst.Blue)
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
	assert.Ok(t, err)
	assert.Assert(t, changed, "expected the destination to be changed")
	assert.Equals(t, "bravo", dst.unexported)
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
	assert.Ok(t, err)
	assert.Assert(t, !changed, "the crayon should not have been changed")
	assert.Equals(t, "Asparagus", c.Name)
	assert.Equals(t, uint8(0x7b), c.Red)
	assert.Equals(t, uint8(0xa0), c.Green)
	assert.Equals(t, uint8(0x5b), c.Blue)

	changed, err = merge.Merge(c, struct{}{})
	assert.Ok(t, err)
	assert.Assert(t, !changed, "the crayon should not have been changed")
}

func TestOnlyStructs(t *testing.T) {
	_, err := merge.Merge(42, Crayon{})
	assert.ErrorIs(t, err, errors.ErrNotAStruct)

	_, err = merge.Merge(Crayon{}, 41)
	assert.ErrorIs(t, err, errors.ErrNotAStruct)
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
