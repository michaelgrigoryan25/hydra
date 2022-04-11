package hydra_test

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"polygon.am/hydra"
)

// BSD 3-Clause License

// Copyright (c) 2021, Michael Grigoryan
// All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:

// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.

// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.

// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

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
	h := hydra.Hydra{Config: hydra.Config{
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
	h := hydra.Hydra{Config: hydra.Config{
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

	_, err := h.Load(&cfg)

	assert.NoError(t, err)
	assert.Equal(t, expectedCfg, cfg)
}

const TestConfigLookupPath = "testdata/"

func TestLoadAndParseConfigs(t *testing.T) {
	type EntryMetadata struct {
		MustFail bool
		Filename string
	}

	// This function will determine which configuration files are
	// expected to be valid or invalid, and automatically creates
	// config file metadata from the filenames.
	//
	// This will check whether the filename contains `err` keyword
	// and based on it will mark the configuration as a failing one.
	scanConfigs := func(path string) ([]EntryMetadata, error) {
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}

		var configs []EntryMetadata
		// Scanning all file entries found in the `testdata/` directory
		for _, entry := range entries {
			configs = append(configs, EntryMetadata{
				Filename: entry.Name(),

				// If the filename contains the `err` keyword, then it must
				// return an error during the process of parsing the YAML.
				MustFail: strings.Contains(entry.Name(), "err"),
			})
		}

		return configs, nil
	}

	c, err := scanConfigs(TestConfigLookupPath)
	if err != nil {
		t.Fatal(err)
	}

	for _, config := range c {
		hydraConfig := hydra.Config{
			Filename: config.Filename,
			Paths:    []string{TestConfigLookupPath}, // `testdata/`
		}

		hydra := hydra.Hydra{
			Config: hydraConfig,
		}

		_, err := hydra.Load(new(interface{}))
		if config.MustFail {
			// `.Load` will return an error when the parsing process fails
			assert.NotNil(t, err)
		}
	}
}
