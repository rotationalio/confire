package usage_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"testing"
	"text/tabwriter"
	"time"

	"github.com/rotationalio/confire/assert"
	"github.com/rotationalio/confire/usage"
)

func TestUsageDefault(t *testing.T) {
	orig := os.Stdout
	t.Cleanup(func() {
		os.Stdout = orig
	})

	r, w, _ := os.Pipe()
	os.Stdout = w

	var s Specification
	err := usage.Usage("confire", &s)

	// copy the output in separate go routine so printing can't block
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	out := <-outC

	assert.Ok(t, err)
	compareUsage(t, "testdata/default_table.txt", out)
}

func TestUsageTable(t *testing.T) {
	buf := &bytes.Buffer{}
	tabs := tabwriter.NewWriter(buf, 1, 0, 4, ' ', 0)

	var s Specification
	err := usage.Usagef("confire", &s, tabs, usage.DefaultTableFormat)
	assert.Ok(t, err)

	tabs.Flush()
	compareUsage(t, "testdata/default_table.txt", buf.String())
}

func TestUsageList(t *testing.T) {
	buf := &bytes.Buffer{}
	tabs := tabwriter.NewWriter(buf, 1, 0, 4, ' ', 0)

	var s Specification
	err := usage.Usagef("confire", &s, tabs, usage.DefaultListFormat)
	assert.Ok(t, err)

	tabs.Flush()
	compareUsage(t, "testdata/default_list.txt", buf.String())
}

func TestUsageCustomFormat(t *testing.T) {
	buf := &bytes.Buffer{}
	tabs := tabwriter.NewWriter(buf, 1, 0, 4, ' ', 0)

	var s Specification
	err := usage.Usagef("confire", &s, tabs, "{{range .}}{{usage_key .}}={{usage_description .}}\n{{end}}")
	assert.Ok(t, err)

	tabs.Flush()
	compareUsage(t, "testdata/custom.txt", buf.String())
}

func TestUnknownKey(t *testing.T) {
	buf := &bytes.Buffer{}
	tabs := tabwriter.NewWriter(buf, 1, 0, 4, ' ', 0)

	var s Specification
	err := usage.Usagef("confire", &s, tabs, "{{.UnknownKey}}")
	assert.Assert(t, err != nil, "expected an error")
	assert.Equals(t, "template: confire:1:2: executing \"confire\" at <.UnknownKey>: can't evaluate field UnknownKey in type []env.Info", err.Error())
}

func TestUsageBadFormat(t *testing.T) {
	buf := &bytes.Buffer{}
	var s Specification
	err := usage.Usagef("env_config", &s, buf, "{{range .}}{.Key}\n{{end}}")
	assert.Ok(t, err)
	compareUsage(t, "testdata/fault.txt", buf.String())
}

type Specification struct {
	UserSpecification
	Debug    bool          `default:"false"`
	Host     string        `required:"true"`
	Port     int           `required:"true"`
	Rate     float64       `default:"0.25" desc:"value between 0 and 1"`
	Cancel   uint16        `default:"16"`
	LogLevel LogLevel      `default:"info" split_words:"true"`
	Timeout  time.Duration `default:"30s" desc:"amount of time to wait for a respone"`
	Colors   struct {
		Primary   Color `default:"#cc6699"`
		Secondary Color `default:"#eeffee"`
	}
	SendGridAPIKey string `env:"SENDGRID_API_KEY"`
}

type UserSpecification struct {
	Admins  map[string]string `required:"true" desc:"map username to password"`
	Roles   []string          `default:"observer,admin"`
	AuthURL *CustomURL        `required:"true"`
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

type CustomURL struct {
	Value *url.URL
}

func (e *CustomURL) Decode(value string) (err error) {
	e.Value, err = url.Parse(value)
	return err
}

func compareUsage(t *testing.T, path, actual string) {
	data, err := os.ReadFile(path)
	assert.Ok(t, err)
	expected := string(data)

	actual = strings.ReplaceAll(actual, " ", "·")
	if actual != expected {
		shortest := len(expected)
		if len(actual) < shortest {
			shortest = len(actual)
		}

		if len(expected) != len(actual) {
			t.Errorf("expected result length %d, found %d", len(expected), len(actual))
		}

		for i := 0; i < shortest; i++ {
			if expected[i] != actual[i] {
				t.Errorf("difference at index %d, expected '%c' (%v), found '%c' (%v)\n", i, expected[i], expected[i], actual[i], actual[i])
				break
			}
		}
		t.Errorf("Complete Expected:\n%s\nComplete Actual:\n%s\n", expected, actual)
	}
}
