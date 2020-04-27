// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestWalk(t *testing.T) {
	for _, tc := range []struct {
		node Node
		want string
	}{
		{
			node: List{
				&Word{Text: "hello"},
				&Word{Text: "world"},
			},
			want: "ast.List *ast.Word <nil> *ast.Word <nil> <nil>",
		},
		{
			node: &Macro{
				Name: &Ident{Name: "sqrt"},
				Args: []Node{
					&OptArg{
						List: []Node{
							&Word{Text: "n"},
						},
					},
					&Arg{
						List: []Node{
							&Literal{Text: "2"},
							&Word{Text: "x"},
						},
					},
				},
			},
			want: "*ast.Macro *ast.Ident <nil> *ast.OptArg *ast.Word <nil> <nil> *ast.Arg *ast.Literal <nil> *ast.Word <nil> <nil> <nil>",
		},
		{
			node: &MathExpr{
				Delim: "$",
				List: []Node{
					&Literal{Text: "2"},
					&Symbol{Text: "+"},
					&Word{Text: "x"},
				},
			},
			want: "*ast.MathExpr *ast.Literal <nil> *ast.Symbol <nil> *ast.Word <nil> <nil>",
		},
		{
			node: &Sub{
				Node: &Word{Text: "i"},
			},
			want: "*ast.Sub *ast.Word <nil> <nil>",
		},
		{
			node: &Sup{
				Node: &Literal{Text: "2"},
			},
			want: "*ast.Sup *ast.Literal <nil> <nil>",
		},
	} {
		t.Run("", func(t *testing.T) {
			o := new(strings.Builder)
			v := &sprinter{o}
			Walk(v, tc.node)
			got := strings.TrimSpace(o.String())
			if got != tc.want {
				t.Fatalf("invalid walk:\ngot= %v\nwant=%v", got, tc.want)
			}
		})
	}
}

type sprinter struct {
	w io.Writer
}

func (p *sprinter) Visit(n Node) Visitor {
	fmt.Fprintf(p.w, "%T ", n)
	return p
}

func TestInspect(t *testing.T) {
	for _, tc := range []struct {
		node Node
		want string
	}{
		{
			node: List{
				&Word{Text: "hello"},
				&Word{Text: "world"},
			},
			want: "ast.List *ast.Word <nil> *ast.Word <nil> <nil>",
		},
		{
			node: &Macro{
				Name: &Ident{Name: "sqrt"},
				Args: []Node{
					&OptArg{
						List: []Node{
							&Word{Text: "n"},
						},
					},
					&Arg{
						List: []Node{
							&Literal{Text: "2"},
							&Word{Text: "x"},
						},
					},
				},
			},
			want: "*ast.Macro *ast.Ident <nil> *ast.OptArg *ast.Word <nil> <nil> *ast.Arg *ast.Literal <nil> *ast.Word <nil> <nil> <nil>",
		},
		{
			node: &MathExpr{
				Delim: "$",
				List: []Node{
					&Literal{Text: "2"},
					&Symbol{Text: "+"},
					&Word{Text: "x"},
				},
			},
			want: "*ast.MathExpr *ast.Literal <nil> *ast.Symbol <nil> *ast.Word <nil> <nil>",
		},
		{
			node: &Sub{
				Node: &Word{Text: "i"},
			},
			want: "*ast.Sub *ast.Word <nil> <nil>",
		},
		{
			node: &Sup{
				Node: &Literal{Text: "2"},
			},
			want: "*ast.Sup *ast.Literal <nil> <nil>",
		},
	} {
		t.Run("", func(t *testing.T) {
			o := new(strings.Builder)
			Inspect(tc.node, func(n Node) bool {
				fmt.Fprintf(o, "%T ", n)
				return true
			})
			got := strings.TrimSpace(o.String())
			if got != tc.want {
				t.Fatalf("invalid inspect:\ngot= %v\nwant=%v", got, tc.want)
			}
		})
	}
}
