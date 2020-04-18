// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

import "testing"

func TestToken(t *testing.T) {
	want := "token text"
	tok := Token{Text: want}

	got := tok.String()
	if got != want {
		t.Fatalf("invalid stringer: got=%q, want=%q", got, want)
	}
}

func TestKind(t *testing.T) {
	want := "Macro"
	kind := Macro
	got := kind.String()
	if got != want {
		t.Fatalf("invalid stringer: got=%q, want=%q", got, want)
	}
}
