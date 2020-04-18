// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate stringer -type Kind

package token // import "github.com/go-latex/latex/token"

import (
	"go/token"
)

// Kind is a kind of LaTeX token.
type Kind int

const (
	Invalid Kind = iota
	Macro
	EmptyLine
	Comment
	Space
	Word
	Number
	Dollar
	Lbrace
	Rbrace
	Lbrack
	Rbrack
	Equal
	Underscore
	Lparen
	Rparen
	Lt
	Gt
	Hat
	Div
	Mul
	Sub
	Add
	Not
	Colon
	Backslash
	Other
	Verbatim
	EOF
)

// Token holds informations about a token.
type Token struct {
	Kind Kind // Kind is the kind of token.
	Pos  Pos  // Pos is the position of a token.
	Text string
}

func (t Token) String() string { return t.Text }

// Pos is a compact encoding of a source position within a file set.
//
// Aliased from go/token.Pos
type Pos = token.Pos
