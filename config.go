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
	"os"
	"path/filepath"
)

// hydra.Config is the configuration struct which is going to be
// used by Hydra for initialization.
type Config struct {
	Paths []string
}

// This function will attempt to get the first usable and correct
// path where the configuration has been specified. Will return an
// error if something goes wrong.
func (c *Config) findConfigPath() (string, error) {
	// In this part, we are looping through all the provided
	// configuration Paths, in search of the first valid path.
	for _, path := range c.Paths {
		path = filepath.Clean(path) // cleaning the path from various redundancies.
		dir, file := filepath.Split(path)

		// Scanning all the files in all the directories provided by
		// the user, and checking file/folder entries.
		if entries, err := os.ReadDir(dir); err == nil {
			for _, i := range entries {
				// Since we only need files and not folders, we can break the loop
				// from here and move forward.
				if i.Type().IsRegular() && i.Name() == file {
					return path, nil
				}
			}
		}
	}

	// If nothing was matched, or the user did not provide any
	// search paths, defaulting to an empty string and moving on
	// to parsing the environment variables.
	return "", nil
}
