package env_test

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	. "github.com/rotationalio/confire/env"
)

const testPrefix = "confire"

type Specification struct {
	Embedded                     `desc:"can we document a struct"`
	EmbeddedButIgnored           `ignored:"true"`
	Debug                        bool
	Port                         int
	Rate                         float32
	User                         string
	TTL                          uint32
	Timeout                      time.Duration
	AdminUsers                   []string
	MagicNumbers                 []int
	EmptyNumbers                 []int
	ByteSlice                    []byte
	ColorCodes                   map[string]int
	MultiWordVar                 string
	MultiWordVarWithAutoSplit    uint32 `split_words:"true"`
	MultiWordACRWithAutoSplit    uint32 `split_words:"true"`
	SomePointer                  *string
	SomePointerWithDefault       *string `default:"foo2baz" desc:"foorbar is the word"`
	MultiWordVarWithAlt          string  `envconfig:"MULTI_WORD_VAR_WITH_ALT" desc:"what alt"`
	MultiWordVarWithLowerCaseAlt string  `envconfig:"multi_word_var_with_lower_case_alt"`
	NoPrefixWithAlt              string  `envconfig:"SERVICE_HOST"`
	DefaultVar                   string  `default:"foobar"`
	RequiredVar                  string  `required:"True"`
	NoPrefixDefault              string  `envconfig:"BROKER" default:"127.0.0.1"`
	RequiredDefault              string  `required:"true" default:"foo2bar"`
	Ignored                      string  `ignored:"true"`
	NestedSpecification          struct {
		Property            string `envconfig:"inner"`
		PropertyWithDefault string `default:"fuzzybydefault"`
	} `envconfig:"outer"`
	AfterNested  string
	DecodeStruct HonorDecodeInStruct `envconfig:"honor"`
	Datetime     time.Time
	MapField     map[string]string `default:"one:two,three:four"`
	UrlValue     CustomURL
	UrlPointer   *CustomURL
}

type Embedded struct {
	Enabled             bool `desc:"some embedded value"`
	EmbeddedPort        int
	MultiWordVar        string
	MultiWordVarWithAlt string `envconfig:"MULTI_WITH_DIFFERENT_ALT"`
	EmbeddedAlt         string `envconfig:"EMBEDDED_WITH_ALT"`
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
	"CONFIRE_BYTE_SLICE":                     "theeaglefliesatmidnight",
	"CONFIRE_COLOR_CODES":                    "red:1,green:2,blue:3",
	"SERVICE_HOST":                           "127.0.0.7",
	"CONFIRE_TTL":                            "30",
	"CONFIRE_REQUIREDVAR":                    "foo",
	"CONFIRE_IGNORED":                        "was-not-ignored",
	"CONFIRE_OUTER_INNER":                    "iamnested",
	"CONFIRE_AFTERNESTED":                    "after",
	"CONFIRE_HONOR":                          "honor",
	"CONFIRE_DATETIME":                       "2023-07-19T14:53:06Z",
	"CONFIRE_MULTI_WORD_VAR_WITH_AUTO_SPLIT": "24",
	"CONFIRE_MULTI_WORD_ACR_WITH_AUTO_SPLIT": "25",
	"CONFIRE_URLVALUE":                       "https://rotational.io/blog/",
	"CONFIRE_URLPOINTER":                     "https://rotational.io/blog/",
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
