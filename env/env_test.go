package env_test

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/rotationalio/confire/assert"
	. "github.com/rotationalio/confire/env"
)

const testPrefix = "confire"

type Specification struct {
	Embedded                        `desc:"can we document a struct"`
	EmbeddedButIgnored              `ignored:"true"`
	Debug                           bool
	Port                            int
	Rate                            float32
	User                            string
	TTL                             uint32
	Timeout                         time.Duration
	AdminUsers                      []string
	MagicNumbers                    []int
	EmptyNumbers                    []int
	ByteSlice                       []byte
	ColorCodes                      map[string]int
	MultiWordVar                    string
	MultiWordVarWithAutoSplit       uint32 `split_words:"true"`
	MultiWordACRWithAutoSplit       uint32 `split_words:"true"`
	SomePointer                     *string
	MultiWordVarWithAlt             string `envconfig:"MULTI_WORD_VAR_WITH_ALT" desc:"what alt"`
	MultiWordVarWithLowerCaseAlt    string `envconfig:"multi_word_var_with_lower_case_alt"`
	NoPrefixWithAlt                 string `envconfig:"SERVICE_HOST"`
	NoPrefixWithAltEnv              string `env:"SERVICE_HOSTNAME"`
	MultiWordVarWithAltEnv          string `env:"MULTI_WORD_VAR_WITH_ALT_ENV" desc:"what alt"`
	MultiWordVarWithLowerCaseAltEnv string `env:"multi_word_var_with_lower_case_alt_env"`
	Ignored                         string `ignored:"true"`
	NestedSpecification             struct {
		Property         string `envconfig:"inner"`
		PropertyNoPrefix string `env:"OUTER_INNER"`
	} `envconfig:"outer"`
	AfterNested  string
	DecodeStruct HonorDecodeInStruct `envconfig:"honor"`
	Datetime     time.Time
	MapField     map[string]string
	UrlValue     CustomURL
	UrlPointer   *CustomURL
}

type Embedded struct {
	Enabled             bool `desc:"some embedded value"`
	EmbeddedPort        int
	MultiWordVar        string
	MultiWordVarWithAlt string `envconfig:"MULTI_WITH_DIFFERENT_ALT"`
	EmbeddedAlt         string `env:"EMBEDDED_WITH_ALT"`
	EmbeddedIgnored     string `ignored:"true"`
}

type EmbeddedButIgnored struct {
	FirstEmbeddedButIgnored  string
	SecondEmbeddedButIgnored string
}

type HonorDecodeInStruct struct {
	Value string
}

func (h *HonorDecodeInStruct) Decode(env string) error {
	h.Value = "decoded"
	return nil
}

type CustomURL struct {
	Value *url.URL
}

func (cu *CustomURL) UnmarshalBinary(data []byte) error {
	u, err := url.Parse(string(data))
	cu.Value = u
	return err
}

type bracketed string

func (b *bracketed) Set(value string) error {
	*b = bracketed("[" + value + "]")
	return nil
}

func (b bracketed) String() string {
	return string(b)
}

// quoted is used to test the precedence of Decode over Set.
// The sole field is a flag.Value rather than a setter to validate that
// all flag.Value implementations are also Setter implementations.
type quoted struct{ flag.Value }

func (d quoted) Decode(value string) error {
	return d.Set(`"` + value + `"`)
}

type setterStruct struct {
	Inner string
}

func (ss *setterStruct) Set(value string) error {
	ss.Inner = fmt.Sprintf("setterstruct{%q}", value)
	return nil
}

var testEnv = map[string]string{
	"CONFIRE_DEBUG":                          "true",
	"CONFIRE_PORT":                           "8888",
	"CONFIRE_RATE":                           "0.25",
	"CONFIRE_USER":                           "werebear",
	"CONFIRE_TIMEOUT":                        "5m",
	"CONFIRE_ADMINUSERS":                     "werewolf,vampire,ghast",
	"CONFIRE_MAGICNUMBERS":                   "3,7,12",
	"CONFIRE_EMPTYNUMBERS":                   "",
	"CONFIRE_BYTESLICE":                      "theeaglefliesatmidnight",
	"CONFIRE_COLORCODES":                     "red:1,green:2,blue:3",
	"SERVICE_HOST":                           "127.0.0.7",
	"CONFIRE_TTL":                            "30",
	"CONFIRE_REQUIREDVAR":                    "foo",
	"CONFIRE_IGNORED":                        "was-not-ignored",
	"CONFIRE_OUTER_INNER":                    "iamnested",
	"OUTER_INNER":                            "iamouterinner",
	"CONFIRE_AFTERNESTED":                    "after",
	"CONFIRE_HONOR":                          "honor",
	"CONFIRE_DATETIME":                       "2023-07-19T14:53:06Z",
	"CONFIRE_MULTI_WORD_VAR_WITH_AUTO_SPLIT": "24",
	"CONFIRE_MULTI_WORD_ACR_WITH_AUTO_SPLIT": "25",
	"CONFIRE_URLVALUE":                       "https://rotational.io/blog/",
	"CONFIRE_URLPOINTER":                     "https://rotational.io/blog/",
}

func TestProcess(t *testing.T) {
	t.Cleanup(cleanupEnv())
	setEnv()

	var s Specification
	err := Process(testPrefix, &s)
	assert.Ok(t, err)

	assert.Equals(t, "127.0.0.7", s.NoPrefixWithAlt)
	assert.True(t, s.Debug)
	assert.Equals(t, 8888, s.Port)
	assert.Equals(t, float32(0.25), s.Rate)
	assert.Equals(t, uint32(30), s.TTL)
	assert.Equals(t, "werebear", s.User)
	assert.Equals(t, 5*time.Minute, s.Timeout)
	assert.Equals(t, []string{"werewolf", "vampire", "ghast"}, s.AdminUsers)
	assert.Equals(t, []int{3, 7, 12}, s.MagicNumbers)
	assert.True(t, len(s.EmptyNumbers) == 0)
	assert.Equals(t, []byte("theeaglefliesatmidnight"), s.ByteSlice)
	assert.Equals(t, "", s.Ignored)
	assert.Equals(t, map[string]int{"red": 1, "green": 2, "blue": 3}, s.ColorCodes)
	assert.Equals(t, "iamnested", s.NestedSpecification.Property)
	assert.Equals(t, "iamouterinner", s.NestedSpecification.PropertyNoPrefix)
	assert.Equals(t, "after", s.AfterNested)
	assert.Equals(t, "decoded", s.DecodeStruct.Value)
	assert.Equals(t, time.Date(2023, 7, 19, 14, 53, 6, 0, time.UTC), s.Datetime)
	assert.Equals(t, uint32(24), s.MultiWordVarWithAutoSplit)
	assert.Equals(t, uint32(25), s.MultiWordACRWithAutoSplit)

	u, err := url.Parse("https://rotational.io/blog/")
	assert.Ok(t, err)

	assert.Equals(t, *u, *s.UrlValue.Value)
	assert.Equals(t, *u, *s.UrlPointer.Value)
}

func TestCustomValueFields(t *testing.T) {
	t.Cleanup(cleanupEnv())

	var s struct {
		Foo    string
		Bar    bracketed
		Baz    quoted
		Struct setterStruct
	}

	// Set would panic when the receiver is nil,
	// so make sure it has an initial value to replace.
	s.Baz = quoted{new(bracketed)}

	os.Setenv("CONFIRE_FOO", "foo")
	os.Setenv("CONFIRE_BAR", "bar")
	os.Setenv("CONFIRE_BAZ", "baz")
	os.Setenv("CONFIRE_STRUCT", "inner")

	err := Process(testPrefix, &s)
	assert.Ok(t, err)

	assert.Equals(t, "foo", s.Foo)
	assert.Equals(t, "[bar]", s.Bar.String())
	assert.Equals(t, `["baz"]`, s.Baz.String())
	assert.Equals(t, `setterstruct{"inner"}`, s.Struct.Inner)
}

func TestCustomPointerFields(t *testing.T) {
	t.Cleanup(cleanupEnv())

	var s struct {
		Foo    string
		Bar    *bracketed
		Baz    *quoted
		Struct *setterStruct
	}

	// Set would panic when the receiver is nil,
	// so make sure they have initial values to replace.
	s.Bar = new(bracketed)
	s.Baz = &quoted{new(bracketed)}

	os.Setenv("CONFIRE_FOO", "foo")
	os.Setenv("CONFIRE_BAR", "bar")
	os.Setenv("CONFIRE_BAZ", "baz")
	os.Setenv("CONFIRE_STRUCT", "inner")

	err := Process(testPrefix, &s)
	assert.Ok(t, err)

	assert.Equals(t, "foo", s.Foo)
	assert.Equals(t, "[bar]", s.Bar.String())
	assert.Equals(t, `["baz"]`, s.Baz.String())
	assert.Equals(t, `setterstruct{"inner"}`, s.Struct.Inner)
}

func BenchmarkGather(b *testing.B) {
	b.Cleanup(cleanupEnv())
	setEnv()

	for i := 0; i < b.N; i++ {
		var s Specification
		Gather(testPrefix, &s)
	}
}

// Returns the current environment for the specified keys, or if no keys are specified
// then it returns the current environment for all keys in the testEnv variable.
func curEnv(keys ...string) map[string]string {
	env := make(map[string]string)
	if len(keys) > 0 {
		for _, key := range keys {
			if val, ok := os.LookupEnv(key); ok {
				env[key] = val
			}
		}
	} else {
		for key := range testEnv {
			env[key] = os.Getenv(key)
		}
	}

	return env
}

// Sets the environment variables from the testEnv variable. If no keys are specified,
// then this function sets all environment variables from the testEnv.
func setEnv(keys ...string) {
	if len(keys) > 0 {
		for _, key := range keys {
			if val, ok := testEnv[key]; ok {
				os.Setenv(key, val)
			}
		}
	} else {
		for key, val := range testEnv {
			os.Setenv(key, val)
		}
	}
}

// Cleanup helper function that can be run when the tests are complete to reset the
// environment back to its previous state before the test was run.
//
// Usage: t.Cleanup(cleanupEnv())
func cleanupEnv(keys ...string) func() {
	prevEnv := curEnv(keys...)
	return func() {
		for key, val := range prevEnv {
			if val != "" {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	}
}
