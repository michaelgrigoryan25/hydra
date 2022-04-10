package hydra

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	StringVal string       `yaml:"string_val" hydra:"env=STRING_VAL"`
	IntVal    int          `yaml:"int_val" hydra:"env=INT_VAL"`
	FloatVar  float64      `yaml:"float_val" hydra:"env=FLOAT_VAL"`
	BoolVar   bool         `yaml:"bool_val" hydra:"env=BOOL_VAL"`
	Nested    NestedStruct `yaml:"nested"`
}

type NestedStruct struct {
	NestedStringVal string `yaml:"nested_string_val" hydra:"env=NESTED_STRING_VAL"`
}

func TestLoadConfig(t *testing.T) {
	h := Hydra{Config: Config{
		Filename: "test.yaml",
		Paths:    []string{"."},
	}}

	expected := TestConfig{
		StringVal: "TESTSTRING",
		IntVal:    666,
		FloatVar:  2.72,
		BoolVar:   true,
		Nested:    NestedStruct{NestedStringVal: "NESTEDSTRING"},
	}

	cfg := TestConfig{}
	_, err := h.Load(&cfg)
	assert.NoError(t, err)

	assert.Equal(t, expected, cfg)
}

func TestParseEnv(t *testing.T) {
	h := Hydra{Config: Config{
		Filename: "test.yaml",
		Paths:    []string{"."},
	}}

	expectedString := "parsed_from_environment"
	expectedInt := 99999
	expectedFloat := 3.14
	expectedBool := true
	expectedNestedString := "NESTED"

	expectedCfg := TestConfig{
		StringVal: expectedString,
		IntVal:    expectedInt,
		FloatVar:  expectedFloat,
		BoolVar:   expectedBool,
		Nested:    NestedStruct{NestedStringVal: expectedNestedString},
	}

	cfg := TestConfig{}

	os.Setenv("STRING_VAL", expectedString)
	os.Setenv("INT_VAL", strconv.Itoa(expectedInt))
	os.Setenv("FLOAT_VAL", strconv.FormatFloat(expectedFloat, 'f', 2, 64))
	os.Setenv("BOOL_VAL", strconv.FormatBool(expectedBool))
	os.Setenv("NESTED_STRING_VAL", expectedNestedString)

	err := h.parseEnv(&cfg)

	assert.NoError(t, err)
	assert.Equal(t, expectedCfg, cfg)
}
