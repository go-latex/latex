// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package symbols

import "testing"

func TestIsSpaced(t *testing.T) {
	for _, tc := range []struct {
		symbol string
		want   bool
	}{
		{`\leftarrow`, true},
		{`\dashv`, true},
		{`=`, true},
		{`<`, true},
		{`+`, true},
		{`\pm`, true},
		{`\sum`, false},
		{`\alpha`, false},
		{` `, false},
	} {
		t.Run(tc.symbol, func(t *testing.T) {
			got := IsSpaced(tc.symbol)
			if got != tc.want {
				t.Fatalf("got: %v, want: %v", got, tc.want)
			}
		})
	}
}
