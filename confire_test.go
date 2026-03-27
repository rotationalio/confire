package confire_test

import (
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"go.rtnl.ai/confire"
	"go.rtnl.ai/confire/assert"
	"go.rtnl.ai/confire/contest"
)

//============================================================================
// Configuration Types
//============================================================================

// Test Configuration
type Config struct {
	Debug       bool          `default:"false" desc:"enable debug mode"`
	ServiceName string        `required:"true" split_words:"true" desc:"name of the service"`
	Host        string        `default:"0.0.0.0" desc:"host to listen on"`
	Port        int           `default:"8080" desc:"port to listen on"`
	Rate        float64       `default:"1.0" desc:"percentage of requests to process"`
	Timeout     time.Duration `default:"10s" desc:"timeout for requests"`
	LogLevel    LevelDecoder  `env:"OTEL_LOG_LEVEL" default:"info" desc:"log level"`
	Database    DatabaseConfig
	UI          UIConfig
}

type DatabaseConfig struct {
	URL      string `env:"DATABASE_URL" required:"true" desc:"database connection URL"`
	ReadOnly bool   `default:"false" desc:"read only mode"`
}

type UIConfig struct {
	Enabled   bool             `default:"true" desc:"enable the customized user interface"`
	Primary   Color            `default:"#cc6699" desc:"primary color"`
	Secondary Color            `default:"#eeffee" desc:"secondary color"`
	Palette   map[string]Color `desc:"palette of colors"`
}

func (c Config) Validate() (err error) {
	if c.Port < 1024 || c.Port > 65535 {
		err = confire.Join(err, confire.Invalid("", "port", "must be in the integer range [1024, 65535]"))
	}

	if c.Rate < 0.0 || c.Rate > 1.0 {
		err = confire.Join(err, confire.Invalid("", "rate", "must be in the float range [0.0, 1.0]"))
	}

	return err
}

func (c *UIConfig) Validate() (err error) {
	if c.Enabled {
		if len(c.Palette) == 0 {
			err = confire.Join(err, confire.Required("ui", "palette"))
		}

		if len(c.Palette) > 8 {
			err = confire.Join(err, confire.Invalid("ui", "palette", "palette must be less than 8 colors"))
		}

		hasPrimary, hasSecondary := false, false
		for _, color := range c.Palette {
			if color == c.Primary {
				hasPrimary = true
			}
			if color == c.Secondary {
				hasSecondary = true
			}
		}

		if !hasPrimary {
			err = confire.Join(err, confire.Invalid("ui", "palette", "primary color must be included in the palette"))
		}

		if !hasSecondary {
			err = confire.Join(err, confire.Invalid("ui", "palette", "secondary color must be included in the palette"))
		}
	}
	return err
}

type LevelDecoder uint8

const (
	LevelDebug LevelDecoder = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
	LevelPanic
)

func (l *LevelDecoder) Decode(value string) error {
	value = strings.TrimSpace(strings.ToLower(value))
	switch value {
	case "debug":
		*l = LevelDebug
	case "info":
		*l = LevelInfo
	case "warning":
		*l = LevelWarning
	case "error":
		*l = LevelError
	case "fatal":
		*l = LevelFatal
	case "panic":
		*l = LevelPanic
	default:
		return confire.Invalid("", "log level", "invalid log level: %q", value)
	}
	return nil
}

type Color [3]uint8

func (c *Color) Decode(value string) error {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.TrimPrefix(value, "#")

	if len(value) != 6 {
		return confire.Invalid("", "color", "invalid color: %q", value)
	}

	if n, err := hex.Decode(c[:], []byte(value)); n != 3 || err != nil {
		return confire.Invalid("", "color", "could not decode hex color")
	}
	return nil
}

//============================================================================
// Configuration Tests
//============================================================================

var testEnv = contest.Env{
	"CONFIRE_DEBUG":             "true",
	"CONFIRE_SERVICE_NAME":      "myapp",
	"CONFIRE_HOST":              "127.0.0.1",
	"CONFIRE_PORT":              "8000",
	"CONFIRE_RATE":              "0.75",
	"CONFIRE_TIMEOUT":           "1m30s",
	"OTEL_LOG_LEVEL":            "warning",
	"DATABASE_URL":              "sqlite://myapp.db",
	"CONFIRE_DATABASE_READONLY": "true",
	"CONFIRE_UI_ENABLED":        "true",
	"CONFIRE_UI_PRIMARY":        "#ff0000",
	"CONFIRE_UI_SECONDARY":      "#00ff00",
	"CONFIRE_UI_PALETTE":        "primary:#ff0000,secondary:#00ff00",
}

var validConfig = Config{
	Debug:       true,
	ServiceName: "myapp",
	Host:        "127.0.0.1",
	Port:        8000,
	Rate:        0.75,
	Timeout:     time.Minute + (30 * time.Second),
	LogLevel:    LevelWarning,
	Database: DatabaseConfig{
		URL:      "sqlite://myapp.db",
		ReadOnly: true,
	},
	UI: UIConfig{
		Enabled:   true,
		Primary:   Color{255, 0, 0},
		Secondary: Color{0, 255, 0},
		Palette: map[string]Color{
			"primary":   {255, 0, 0},
			"secondary": {0, 255, 0},
		},
	},
}

func TestConfig(t *testing.T) {
	t.Cleanup(testEnv.Set())

	var conf Config
	err := confire.Process("confire", &conf)
	assert.Ok(t, err)

	assert.Equals(t, validConfig, conf)
}

func TestValidation(t *testing.T) {

	t.Run("Defaults", func(t *testing.T) {
		t.Cleanup(testEnv.Clear())

		var conf Config
		err := confire.Process("confire", &conf)
		assert.NotOk(t, err)
		assert.True(t, confire.IsValidationErrors(err))

		errs, ok := confire.ValidationErrors(err)
		assert.True(t, ok)
		assert.Assert(t, len(errs) == 5, "expected 5 validation errors got %d", len(errs))
		assert.Equals(t, "invalid configuration: ServiceName is required but not set", errs[0].Error())
		assert.Equals(t, "invalid configuration: URL is required but not set", errs[1].Error())
		assert.Equals(t, "invalid configuration: ui.palette is required but not set", errs[2].Error())
		assert.Equals(t, "invalid configuration: ui.palette primary color must be included in the palette", errs[3].Error())
		assert.Equals(t, "invalid configuration: ui.palette secondary color must be included in the palette", errs[4].Error())
	})

	t.Run("SoWrong", func(t *testing.T) {
		env := contest.Env{
			"CONFIRE_PORT":       "22",
			"CONFIRE_RATE":       "1.2",
			"DATABASE_URL":       "",
			"CONFIRE_UI_ENABLED": "true",
			"CONFIRE_UI_PALETTE": "primary:#000000,secondary:#000000,warning:#000000,error:#000000,fatal:#000000,panic:#000000",
		}

		t.Cleanup(env.Set())

		var conf Config
		err := confire.Process("confire", &conf)
		assert.NotOk(t, err)
		assert.Assert(t, confire.IsValidationErrors(err), "expected validation errors got %s", err.Error())

		errs, ok := confire.ValidationErrors(err)
		assert.True(t, ok)
		assert.Assert(t, len(errs) == 6, "expected 6 validation errors got %d", len(errs))
		assert.Equals(t, "invalid configuration: port must be in the integer range [1024, 65535]", errs[0].Error())
		assert.Equals(t, "invalid configuration: rate must be in the float range [0.0, 1.0]", errs[1].Error())
		assert.Equals(t, "invalid configuration: ServiceName is required but not set", errs[2].Error())
		assert.Equals(t, "invalid configuration: URL is required but not set", errs[3].Error())
		assert.Equals(t, "invalid configuration: ui.palette primary color must be included in the palette", errs[4].Error())
		assert.Equals(t, "invalid configuration: ui.palette secondary color must be included in the palette", errs[5].Error())
	})
}
