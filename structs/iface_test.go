package structs_test

import (
	"testing"
	"time"

	"github.com/rotationalio/confire/structs"
)

func TestTimeIsZero(t *testing.T) {
	ts := time.Time{}
	var i interface{} = ts

	_, ok := i.(structs.Zero)
	assert(t, ok, "time must implement the zero interface")

	var p interface{} = &ts
	_, ok = p.(structs.Zero)
	assert(t, ok, "time pointer must implement the zero interface")
}
