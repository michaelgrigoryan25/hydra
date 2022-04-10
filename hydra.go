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
package hydra

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// This struct will initialize a new Hydra instance which will
// take care of handling your configuration.
type Hydra struct {
	Config Config
}

// This function will attempt to open, read and parse the file at the
// provided path into YAML.
func (h *Hydra) readAndParse(path string, dst any) error {
	// Opening the file in readonly mode, to not cause accidental damage
	// to the configuration.
	f, err := os.OpenFile(path, os.O_RDONLY, fs.ModeTemporary)
	if err != nil {
		return err
	}

	c, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(c, dst)
	if err != nil {
		return err
	}

	return nil
}

// This function will parse environment variables for the struct fields
// which present 'env=VALUE' key-value pair in their hydra struct tag
func (h *Hydra) parseEnv(dst any) error {
	typeOf := reflect.TypeOf(dst)
	if typeOf.Kind() != reflect.Pointer || typeOf.Elem().Kind() != reflect.Struct {
		return errors.New("destination argument must be a pointer to struct, but it is a " + typeOf.Kind().String())
	}

	val := reflect.ValueOf(dst).Elem()

	fields := reflect.VisibleFields(val.Type())

	for _, f := range fields {
		//Parse recursively in case of nested structs
		if f.Type.Kind() == reflect.Struct {
			err := h.parseEnv(val.FieldByIndex(f.Index).Addr().Interface())
			if err != nil {
				return err
			}
			continue
		}
		if envVarName, ok := getHydraTag(f.Tag, "env"); ok {
			envVar := os.Getenv(envVarName)
			//skip if environment variable is empty
			if envVar == "" {
				continue
			}
			//convert the string value to corresponding type
			switch f.Type.Kind() {
			case reflect.String:
				val.FieldByIndex(f.Index).SetString(envVar)
			case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
				value, err := strconv.ParseInt(envVar, 10, 64)
				if err != nil {
					return err
				} else {
					val.FieldByIndex(f.Index).SetInt(value)
				}
			case reflect.Float64, reflect.Float32:
				value, err := strconv.ParseFloat(envVar, 64)
				if err != nil {
					return err
				} else {
					val.FieldByIndex(f.Index).SetFloat(value)

				}
			case reflect.Bool:
				value, err := strconv.ParseBool(envVar)
				if err != nil {
					return err
				} else {
					val.FieldByIndex(f.Index).SetBool(value)

				}
			}

		}
	}
	return nil
}

// This function checks if the given key exists, and extracts the value if there is any
// from the struct field's hydra struct tag.
func getHydraTag(field reflect.StructTag, tag string) (string, bool) {
	keyValues := strings.Split(field.Get("hydra"), ";")
	for _, t := range keyValues {
		split := strings.Split(t, "=")
		//check if it is only key or a key-value pair
		if len(split) > 1 {
			if split[0] == tag {
				return split[1], true
			}
		} else {
			if t == tag {
				return "", true
			}
		}
	}
	return "", false
}

// This function will attempt to load all configuration variables
// both, from the environment and the YAML configuration file,
// which must be specified when initializing a `Hydra` struct.
func (h *Hydra) Load(dst any) (any, error) {
	p, err := h.Config.findConfigPath()
	if err != nil {
		return nil, err
	}

	err = h.readAndParse(p, dst)
	if err != nil {
		return nil, err
	}

	err = h.parseEnv(dst)
	if err != nil {
		return nil, err
	}

	return dst, nil
}
