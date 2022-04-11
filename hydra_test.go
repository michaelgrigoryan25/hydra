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
package hydra_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"polygon.am/hydra"
)

const TestConfigLookupPath = "testdata/"

type SampleConfig struct {
	Int       int     `yaml:"int" hydra:"env=INT_ENV_VAR"`
	Bool      bool    `yaml:"bool" hydra:"env=BOOL_ENV_VAR"`
	Float     float64 `yaml:"float" hydra:"env=FLOAT_ENV_VAR"`
	String    string  `yaml:"string" hydra:"env=STRING_ENV_VAR"`
	Subconfig struct {
		Nested string `yaml:"nested" hydra:"env=NESTED_ENV_VAR"`
	} `yaml:"subconfig"`
}

func TestLoadConfig(t *testing.T) {
	h := hydra.Hydra{Config: hydra.Config{
		Filename: "test-load.ok.yaml",
		Paths:    []string{TestConfigLookupPath},
	}}

	expected := SampleConfig{
		Bool:   true,
		Int:    12345,
		String: "test",
		Float:  1.2345,
		Subconfig: struct {
			Nested string "yaml:\"nested\" hydra:\"env=NESTED_ENV_VAR\""
		}{
			Nested: "test-nested",
		},
	}

	var config SampleConfig
	_, err := h.Load(&config)
	assert.NoError(t, err)
	assert.Equal(t, expected, config)
}

func TestParseEnv(t *testing.T) {
	h := hydra.Hydra{Config: hydra.Config{
		Filename: "test-load.ok.yaml",
		Paths:    []string{TestConfigLookupPath},
	}}

	type ExpectedValue struct {
		key   string
		value any
	}

	expectedValues := []ExpectedValue{
		{
			value: "test",
			key:   "STRING_ENV_VAR",
		},
		{
			value: 12345,
			key:   "INT_ENV_VAR",
		},
		{
			value: 1.2345,
			key:   "FLOAT_ENV_VAR",
		},
		{
			value: true,
			key:   "BOOL_ENV_VAR",
		},
		{
			value: "test-nested",
			key:   "NESTED_ENV_VAR",
		},
	}

	for _, v := range expectedValues {
		os.Setenv(v.key, fmt.Sprint(v.value))
	}

	expected := SampleConfig{
		String: expectedValues[0].value.(string),
		Int:    expectedValues[1].value.(int),
		Float:  expectedValues[2].value.(float64),
		Bool:   expectedValues[3].value.(bool),
		Subconfig: struct {
			Nested string "yaml:\"nested\" hydra:\"env=NESTED_ENV_VAR\""
		}{
			Nested: expectedValues[4].value.(string),
		},
	}

	var config SampleConfig
	_, err := h.Load(&config)

	assert.NoError(t, err)
	assert.Equal(t, expected, config)
}

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
	assert.NoError(t, err)

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
			assert.Error(t, err)
		}
	}
}
