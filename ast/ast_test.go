// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import (
	"strings"
	"testing"

	"github.com/go-latex/latex/token"
)

func TestPrint(t *testing.T) {
	for _, tc := range []struct {
		node Node
		want string
		pos  token.Pos
	}{
		{
			node: nil,
			want: "<nil>",
		},
		{
			node: &Macro{
				Name: &Ident{
					NamePos: 42,
					Name:    `\cos`,
				},
				Args: nil,
			},
			pos:  42,
			want: `ast.Macro{"\\cos"}`,
		},
		{
			node: &Macro{
				Name: &Ident{
					NamePos: 42,
					Name:    `\sqrt`,
				},
				Args: List{
					&Arg{List: List{&Word{Text: "x"}}},
				},
			},
			pos:  42,
			want: `ast.Macro{"\\sqrt", Args:{ast.Word{"x"}}}`,
		},
		{
			node: &Macro{
				Name: &Ident{
					NamePos: 42,
					Name:    `\sqrt`,
				},
				Args: List{
					&OptArg{List: List{&Word{Text: "n"}}},
					&Arg{List: List{&Word{Text: "x"}}},
				},
			},
			pos:  42,
			want: `ast.Macro{"\\sqrt", Args:[ast.Word{"n"}], {ast.Word{"x"}}}`,
		},
		{
			node: &Word{Text: "hello"},
			want: `ast.Word{"hello"}`,
		},
		{
			node: &Symbol{Text: "$"},
			want: `ast.Symbol{"$"}`,
		},
		{
			node: &Literal{Text: "10"},
			want: `ast.Lit{"10"}`,
		},
		{
			node: &Sup{Node: &Literal{Text: "10"}},
			want: `ast.Sup{ast.Lit{"10"}}`,
		},
		{
			node: &Sub{Node: &Literal{Text: "10"}},
			want: `ast.Sub{ast.Lit{"10"}}`,
		},
		{
			node: List{&Literal{Text: "1"}, &Literal{Text: "2"}},
			want: `ast.List{ast.Lit{"1"}, ast.Lit{"2"}}`,
		},
		{
			node: &MathExpr{
				List: List{
					&Literal{Text: "1"},
					&Word{Text: "x"},
				},
			},
			want: `ast.MathExpr{List:ast.Lit{"1"}, ast.Word{"x"}}`,
		},
		{
			node: &Ident{Name: `\cos`},
			want: `ast.Ident{"\\cos"}`,
		},
	} {
		t.Run("", func(t *testing.T) {
			o := new(strings.Builder)
			Print(o, tc.node)

			if got, want := o.String(), tc.want; got != want {
				t.Fatalf("error:\ngot=%v\nwant=%v", got, want)
			}

			if tc.node == nil {
				return
			}

			tc.node.isNode()
			if got, want := tc.node.Pos(), tc.pos; got != want {
				t.Fatalf("invalid node position: got=%v, want=%v", got, want)
			}
			_ = tc.node.End()
		})
	}
}
