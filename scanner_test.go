// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package latex

import (
	"log"
	"strings"
	"testing"
)

func TestScanner(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
	}{
		{
			name:  "math",
			input: `$\sigma_1 = 22x$ ? ok`,
		},
		{
			name:  "",
			input: `$\sqrt{\frac{e^{3i\pi}}{2\cos 3\pi}}$`,
		},
		{
			name:  "",
			input: `\textbf{APLAS} Dummy -- $\sqrt{s}=13\,$TeV $\mathcal{L}\,=\,3\,ab^{-1}$`,
		},
		{
			name:  "comment",
			input: "% boo is 42\r\n%% bar\tis not boo",
		},
		{
			name:  "",
			input: "hello\n\\\\world!\\ boo",
		},
		{
			name:  "numbers",
			input: "x=23.4\ny=42.\nz=43.x\nw=0x32\nu='c'\nv=\"hello\"",
		},
		{
			name:  "chars",
			input: `x='cos'`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sc := newScanner(strings.NewReader(tc.input))
			for sc.Next() {
				tok := sc.Token()
				log.Printf("tok: %#v", tok)
			}
		})
	}
}
