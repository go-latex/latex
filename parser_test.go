// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package latex provides types and functions to work with LaTeX.
package latex // import "github.com/go-latex/latex"

import (
	"reflect"
	"strings"
	"testing"

	"github.com/go-latex/latex/ast"
)

func TestParser(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  ast.Node
	}{
		{
			input: `hello`,
			want:  ast.List{&ast.Word{Text: "hello"}},
		},
		{
			input: `hello world`,
			want: ast.List{
				&ast.Word{Text: "hello"},
				&ast.Symbol{Text: " "},
				&ast.Word{Text: "world"},
			},
		},
		{
			input: `empty equation $$`,
			want: ast.List{
				&ast.Word{Text: "empty"},
				&ast.Symbol{Text: " "},
				&ast.Word{Text: "equation"},
				&ast.Symbol{Text: " "},
				&ast.MathExpr{
					Delim: "$",
				},
			},
		},
		{
			input: `$+10x$`,
			want: ast.List{
				&ast.MathExpr{
					Delim: "$",
					List: ast.List{
						&ast.Symbol{Text: "+"},
						&ast.Literal{Text: "10"},
						&ast.Word{Text: "x"},
					},
				},
			},
		},
		{
			input: `${}+10x$`,
			want: ast.List{
				&ast.MathExpr{
					Delim: "$",
					List: ast.List{
						ast.List{}, // FIXME(sbinet): shouldn't this be a "group"?
						&ast.Symbol{Text: "+"},
						&ast.Literal{Text: "10"},
						&ast.Word{Text: "x"},
					},
				},
			},
		},
		{
			input: `$\cos$`,
			want: ast.List{
				&ast.MathExpr{
					Delim: "$",
					List: ast.List{
						&ast.Macro{
							Name: &ast.Ident{Name: `\cos`},
						},
					},
				},
			},
		},
		{
			input: `$\sqrt{2x\pi}$`,
			want: ast.List{
				&ast.MathExpr{
					Delim: "$",
					List: ast.List{
						&ast.Macro{
							Name: &ast.Ident{Name: `\sqrt`},
							Args: ast.List{
								&ast.Arg{
									List: ast.List{
										&ast.Literal{
											Text: "2",
										},
										&ast.Word{
											Text: "x",
										},
										&ast.Macro{Name: &ast.Ident{Name: `\pi`}},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: `$\sqrt[3]{2x\pi}$`,
			want: ast.List{
				&ast.MathExpr{
					Delim: "$",
					List: ast.List{
						&ast.Macro{
							Name: &ast.Ident{Name: `\sqrt`},
							Args: ast.List{
								&ast.OptArg{
									List: ast.List{
										&ast.Literal{
											Text: "3",
										},
									},
								},
								&ast.Arg{
									List: ast.List{
										&ast.Literal{
											Text: "2",
										},
										&ast.Word{
											Text: "x",
										},
										&ast.Macro{Name: &ast.Ident{Name: `\pi`}},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: `$\sqrt[n]{2x\pi}$`,
			want: ast.List{
				&ast.MathExpr{
					Delim: "$",
					List: ast.List{
						&ast.Macro{
							Name: &ast.Ident{Name: `\sqrt`},
							Args: ast.List{
								&ast.OptArg{
									List: ast.List{
										&ast.Word{
											Text: "n",
										},
									},
								},
								&ast.Arg{
									List: ast.List{
										&ast.Literal{
											Text: "2",
										},
										&ast.Word{
											Text: "x",
										},
										&ast.Macro{Name: &ast.Ident{Name: `\pi`}},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: `$\exp{2x\pi}$`,
			want: ast.List{
				&ast.MathExpr{
					Delim: "$",
					List: ast.List{
						&ast.Macro{
							Name: &ast.Ident{Name: `\exp`},
							Args: ast.List{
								&ast.Arg{
									List: ast.List{
										&ast.Literal{
											Text: "2",
										},
										&ast.Word{
											Text: "x",
										},
										&ast.Macro{Name: &ast.Ident{Name: `\pi`}},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: `$e^\pi$`,
			want: ast.List{
				&ast.MathExpr{
					Delim: "$",
					List: ast.List{
						&ast.Word{Text: "e"},
						&ast.Sup{Node: &ast.Macro{
							Name: &ast.Ident{Name: `\pi`},
						}},
					},
				},
			},
		},
		{
			input: `$\mathcal{L}$`,
			want: ast.List{
				&ast.MathExpr{
					Delim: "$",
					List: ast.List{
						&ast.Macro{
							Name: &ast.Ident{Name: `\mathcal`},
							Args: ast.List{
								&ast.Arg{
									List: ast.List{
										&ast.Word{Text: "L"}, // FIXME: or Ident?
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: `$\frac{num}{den}$`,
			want: ast.List{
				&ast.MathExpr{
					Delim: "$",
					List: ast.List{
						&ast.Macro{
							Name: &ast.Ident{Name: `\frac`},
							Args: ast.List{
								&ast.Arg{
									List: ast.List{
										&ast.Word{Text: "num"},
									},
								},
								&ast.Arg{
									List: ast.List{
										&ast.Word{Text: "den"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: `$\sqrt{\frac{e^{3i\pi}}{2\cos 3\pi}}$`,
			want: ast.List{
				&ast.MathExpr{
					List: ast.List{
						&ast.Macro{
							Name: &ast.Ident{Name: `\sqrt`},
							Args: ast.List{
								&ast.Arg{
									List: ast.List{
										&ast.Macro{
											Name: &ast.Ident{Name: `\frac`},
											Args: ast.List{
												&ast.Arg{
													List: ast.List{
														&ast.Word{Text: "e"},
														&ast.Sup{Node: ast.List{
															&ast.Literal{Text: "3"},
															&ast.Word{Text: "i"},
															&ast.Macro{Name: &ast.Ident{Name: "\\pi"}},
														}},
													},
												},
												&ast.Arg{
													List: ast.List{
														&ast.Literal{Text: "2"},
														&ast.Macro{Name: &ast.Ident{Name: `\cos`}},
														&ast.Literal{Text: "3"},
														&ast.Macro{Name: &ast.Ident{Name: `\pi`}},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: `$\sqrt{\frac{e^{3i\pi}}{2\cos 3\pi}}$ \textbf{APLAS} Dummy -- $\sqrt{s}=13\,$TeV $\mathcal{L}\,=\,3\,ab^{-1}$`,
			want: ast.List{
				&ast.MathExpr{
					List: ast.List{
						&ast.Macro{
							Name: &ast.Ident{Name: `\sqrt`},
							Args: ast.List{
								&ast.Arg{
									List: ast.List{
										&ast.Macro{
											Name: &ast.Ident{Name: `\frac`},
											Args: ast.List{
												&ast.Arg{
													List: ast.List{
														&ast.Word{Text: "e"},
														&ast.Sup{Node: ast.List{
															&ast.Literal{Text: "3"},
															&ast.Word{Text: "i"},
															&ast.Macro{Name: &ast.Ident{Name: "\\pi"}},
														}},
													},
												},
												&ast.Arg{
													List: ast.List{
														&ast.Literal{Text: "2"},
														&ast.Macro{Name: &ast.Ident{Name: `\cos`}},
														&ast.Literal{Text: "3"},
														&ast.Macro{Name: &ast.Ident{Name: `\pi`}},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				&ast.Symbol{Text: " "},
				&ast.Macro{
					Name: &ast.Ident{Name: `\textbf`},
					Args: ast.List{
						&ast.Arg{
							List: ast.List{
								&ast.Word{Text: "APLAS"},
							},
						},
					},
				},
				&ast.Symbol{Text: " "},
				&ast.Word{Text: "Dummy"},
				&ast.Symbol{Text: " "},
				&ast.Symbol{Text: "-"},
				&ast.Symbol{Text: "-"},
				&ast.Symbol{Text: " "},
				&ast.MathExpr{
					List: ast.List{
						&ast.Macro{
							Name: &ast.Ident{Name: "\\sqrt"},
							Args: ast.List{
								&ast.Arg{
									List: ast.List{
										&ast.Word{Text: "s"},
									},
								},
							},
						},
						&ast.Symbol{Text: "="},
						&ast.Literal{Text: "13"},
						&ast.Macro{Name: &ast.Ident{Name: "\\,"}},
					},
				},
				&ast.Word{Text: "TeV"},
				&ast.Symbol{Text: " "},
				&ast.MathExpr{
					List: ast.List{
						&ast.Macro{
							Name: &ast.Ident{Name: "\\mathcal"},
							Args: ast.List{
								&ast.Arg{
									List: ast.List{
										&ast.Word{Text: "L"},
									},
								},
							},
						},
						&ast.Macro{Name: &ast.Ident{Name: "\\,"}},
						&ast.Symbol{Text: "="},
						&ast.Macro{Name: &ast.Ident{Name: "\\,"}},
						&ast.Literal{Text: "3"},
						&ast.Macro{Name: &ast.Ident{Name: "\\,"}},
						&ast.Word{Text: "ab"},
						&ast.Sup{
							Node: ast.List{
								&ast.Symbol{Text: "-"},
								&ast.Literal{Text: "1"},
							},
						},
					},
				},
			},
		},
		//	{ // FIXME(sbinet): not ready
		//		input: `\[x =3\]`,
		//		want:  nil,
		//	},
		//	{ // FIXME(sbinet): not ready
		//		input: `\(x =3\)`,
		//		want:  nil,
		//	},
		//	{ // FIXME(sbinet): not ready
		//		input: `\begin{equation}x=3\end{equation}`,
		//		want: nil,
		//	},
		{
			input: `$x_i$`,
			want: ast.List{
				&ast.MathExpr{
					List: ast.List{
						&ast.Word{Text: "x"},
						&ast.Sub{
							Node: &ast.Word{Text: "i"},
						},
					},
				},
			},
		},
		{
			input: `$x^n$`,
			want: ast.List{
				&ast.MathExpr{
					List: ast.List{
						&ast.Word{Text: "x"},
						&ast.Sup{
							Node: &ast.Word{Text: "n"},
						},
					},
				},
			},
		},
		{
			input: `$\sum_{i=0}^{n}$`,
			want: ast.List{
				&ast.MathExpr{
					List: ast.List{
						&ast.Macro{
							Name: &ast.Ident{Name: `\sum`},
						},
						&ast.Sub{
							Node: ast.List{
								&ast.Word{Text: "i"},
								&ast.Symbol{Text: "="},
								&ast.Literal{Text: "0"},
							},
						},
						&ast.Sup{
							Node: ast.List{
								&ast.Word{Text: "n"},
							},
						},
					},
				},
			},
		},
	} {
		t.Run("", func(t *testing.T) {
			node, err := ParseExpr(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			got := new(strings.Builder)
			ast.Print(got, node)

			want := new(strings.Builder)
			ast.Print(want, tc.want)

			if got.String() != want.String() {
				t.Fatalf("invalid ast:\ngot: %v\nwant:%v", got, want)
			}
		})
	}

}

func TestTokenPos(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  ast.Node
	}{
		{
			input: `hello`,
			want:  ast.List{&ast.Word{WordPos: 0, Text: "hello"}},
		},
		{
			input: `hello world`,
			want: ast.List{
				&ast.Word{Text: "hello"},
				&ast.Symbol{Text: " ", SymPos: 5},
				&ast.Word{Text: "world", WordPos: 6},
			},
		},
		{
			input: `empty equation $$`,
			want: ast.List{
				&ast.Word{Text: "empty"},
				&ast.Symbol{Text: " ", SymPos: 5},
				&ast.Word{Text: "equation", WordPos: 6},
				&ast.Symbol{Text: " ", SymPos: 14},
				&ast.MathExpr{
					Delim: "$",
					Left:  15,
					Right: 16,
				},
			},
		},
		{
			input: `$+10x$`,
			want: ast.List{
				&ast.MathExpr{
					Delim: "$",
					Left:  0,
					List: ast.List{
						&ast.Symbol{Text: "+", SymPos: 1},
						&ast.Literal{Text: "10", LitPos: 2},
						&ast.Word{Text: "x", WordPos: 4},
					},
					Right: 5,
				},
			},
		},
		{
			input: `$\sqrt{2x\pi}$`,
			want: ast.List{
				&ast.MathExpr{
					Delim: "$",
					Left:  0,
					List: ast.List{
						&ast.Macro{
							Name: &ast.Ident{Name: `\sqrt`, NamePos: 1},
							Args: ast.List{
								&ast.Arg{
									Lbrace: 6,
									List: ast.List{
										&ast.Literal{
											Text:   "2",
											LitPos: 7,
										},
										&ast.Word{
											Text:    "x",
											WordPos: 8,
										},
										&ast.Macro{Name: &ast.Ident{
											Name:    `\pi`,
											NamePos: 9,
										}},
									},
									Rbrace: 12,
								},
							},
						},
					},
					Right: 13,
				},
			},
		},
		{
			input: `$e^\pi$`,
			want: ast.List{
				&ast.MathExpr{
					Delim: "$",
					Left:  0,
					List: ast.List{
						&ast.Word{Text: "e", WordPos: 1},
						&ast.Sup{
							HatPos: 2,
							Node: &ast.Macro{
								Name: &ast.Ident{Name: `\pi`, NamePos: 3},
							},
						},
					},
					Right: 6,
				},
			},
		},
		//	{ // FIXME(sbinet): not ready
		//		input: `\[x =3\]`,
		//		want:  nil,
		//	},
		//	{ // FIXME(sbinet): not ready
		//		input: `\(x =3\)`,
		//		want:  nil,
		//	},
		//	{ // FIXME(sbinet): not ready
		//		input: `\begin{equation}x=3\end{equation}`,
		//		want: nil,
		//	},
		{
			input: `$x_i$`,
			want: ast.List{
				&ast.MathExpr{
					Delim: "$",
					Left:  0,
					List: ast.List{
						&ast.Word{Text: "x", WordPos: 1},
						&ast.Sub{
							UnderPos: 2,
							Node:     &ast.Word{Text: "i", WordPos: 3},
						},
					},
					Right: 4,
				},
			},
		},
	} {
		t.Run("", func(t *testing.T) {
			node, err := ParseExpr(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			if got, want := node, tc.want; !reflect.DeepEqual(got, want) {
				t.Fatalf("invalid positions:\ngot= %v\nwant=%v", got, want)
			}
		})
	}

}
