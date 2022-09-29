/*
MIT License

(C) Copyright 2022 Hewlett Packard Enterprise Development LP

Permission is hereby granted, free of charge, to any person obtaining a
copy of this software and associated documentation files (the "Software"),
to deal in the Software without restriction, including without limitation
the rights to use, copy, modify, merge, publish, distribute, sublicense,
and/or sell copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included
in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

*/

package version

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/Cray-HPE/gru/pkg/buildinfo"
	"github.com/stretchr/testify/assert"
)

func TestPrintVersion(t *testing.T) {
	// set up some non-empty buildinfo values, but put them back to their
	// defaults at the end of the test
	var (
		origVersion      = buildinfo.Version
		origGitSHA       = buildinfo.GitSHA
		origGitTreeState = buildinfo.GitTreeState
	)
	defer func() {
		buildinfo.Version = origVersion
		buildinfo.GitSHA = origGitSHA
		buildinfo.GitTreeState = origGitTreeState
	}()
	buildinfo.Version = "v1.0.0"
	buildinfo.GitSHA = "somegitsha"
	buildinfo.GitTreeState = "dirty"

	clientVersion := fmt.Sprintf("Client:\n\tVersion: %s\n\tGit commit: %s\n", buildinfo.Version, buildinfo.FormattedGitSHA())

	tests := []struct {
		name string
		want string
	}{
		{
			name: "version returns expected format",
			want: clientVersion,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf = new(bytes.Buffer)
			printVersion(buf)
			assert.Equal(t, tc.want, buf.String())
		})
	}
}
