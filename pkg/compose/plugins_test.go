/*
   Copyright 2020 Docker Compose CLI authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package compose

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/cli/cli-plugins/metadata"
	"gotest.tools/v3/assert"
)

func TestFindPluginBinary(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, metadata.NamePrefix+executable("acme"))
	assert.NilError(t, os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755))

	// Found in a configured extra directory.
	assert.Equal(t, bin, findPluginBinary("acme", []string{dir}))

	// Unknown provider is not resolved.
	assert.Equal(t, "", findPluginBinary("missing", []string{dir}))

	// A directory matching the expected name is not a valid binary.
	assert.NilError(t, os.Mkdir(filepath.Join(dir, metadata.NamePrefix+executable("dir")), 0o755))
	assert.Equal(t, "", findPluginBinary("dir", []string{dir}))
}

// TestSystemPluginDirsExcludeProgramData is a regression test for
// CVE-2025-15558: the world-writable %PROGRAMDATA%\Docker\cli-plugins path must
// never appear in the system plugin search paths.
func TestSystemPluginDirsExcludeProgramData(t *testing.T) {
	for _, dir := range systemPluginDirs {
		assert.Assert(t, !strings.Contains(dir, "ProgramData"),
			"system plugin dir %q must not include the insecure ProgramData path", dir)
	}
}
