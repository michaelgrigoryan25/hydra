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
	"fmt"
	"os"
	"path/filepath"
)

// This is the configuration struct that is going to be
// used by Hydra for initialization.
type Config struct {
	Filename string
	Paths    []string
}

// This function will attempt to get the first usable and correct
// path where the configuration has been specified. Will return an
// error if something goes wrong.
func (c *Config) findConfigPath() (string, error) {
	// Hydra configuration must contain at least 1 configuration path.
	if len(c.Paths) == 0 {
		err := fmt.Sprintf("must specify at least 1 config search path. found: %v", len(c.Paths))
		return "", errors.New(err)
	} else {
		var value string

		// In this part, we are looping through all the provided
		// configuration Paths, in search of the first valid path.
		for _, path := range c.Paths {
			// Getting the absolute path of the configuration and
			// returning an error if something fails.
			absolute, err := filepath.Abs(path)
			if err != nil {
				return "", err
			}

			// Scanning the contents of the matched directory
			// and returning an error if something fails.
			if entries, err := os.ReadDir(absolute); err != nil {
				return "", err
			} else {
				for _, entry := range entries {
					// We only need files, not directories, which will contain
					// the filename, specified by the user. Skipping the ones
					// that do not match this criteria.
					if entry.Type().IsRegular() && entry.Name() == c.Filename {
						// Getting the full path of the configuration file
						path, err := filepath.Abs(filepath.Join(absolute, entry.Name()))
						if err != nil {
							return "", err
						}

						// Value will be a valid string, since the error is
						// handled before return.
						value = path
						break
					}
				}
			}
		}

		return value, nil
	}
}
