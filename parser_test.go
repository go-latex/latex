// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package latex provides types and functions to work with LaTeX.
package latex // import "github.com/go-latex/latex"

import (
	"strings"
	"testing"

	"github.com/go-latex/latex/ast"
)

func TestParser(t *testing.T) {
	for _, tc := range []struct {
		input string
	}{
		{input: `hello`},
		{input: `hello world`},
		{input: `empty equation $$`},
		{input: `$+10x$`},
		{input: `${}+10x$`},
		{input: `$\cos$`},
		{input: `$\sqrt{2x\pi}$`},
		{input: `$\sqrt[3]{2x\pi}$`},
		{input: `$\sqrt[n]{2x\pi}$`},
		{input: `$\exp{2x\pi}$`},
		{input: `$e^\pi$`},
		{input: `$\mathcal{L}$`},
		{input: `$\frac{num}{den}$`},
		{input: `$\sqrt{\frac{e^{3i\pi}}{2\cos 3\pi}}$ \textbf{APLAS} Dummy -- $\sqrt{s}=13\,$TeV $\mathcal{L}\,=\,3\,ab^{-1}$`},
	} {
		t.Run("", func(t *testing.T) {
			node, err := ParseExpr(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			o := new(strings.Builder)
			ast.Print(o, node)
			t.Logf("node: %v", o)
		})
	}

}
