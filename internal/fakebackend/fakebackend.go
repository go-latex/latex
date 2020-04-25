// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fakebackend provides a fake tex Backend for testing purposes.
package fakebackend // import "github.com/go-latex/latex/internal/fakebackend"

//go:generate go run ./gen-fakebackend.go

import (
	"fmt"

	"github.com/go-latex/latex/font"
)

type dbFonts map[fontKey]font.Metrics
type dbKerns map[kernKey]float64

var (
	fontsDb dbFonts
	kernsDb dbKerns
)

type Backend struct {
	fonts dbFonts
	kerns dbKerns
}

func New() *Backend {
	return &Backend{
		fonts: fontsDb,
		kerns: kernsDb,
	}
}

// RenderGlyphs renders the glyph g at the reference point (x,y).
func (be *Backend) RenderGlyph(x, y float64, font font.Font, symbol string, dpi float64) {
	//panic("not implemented")
}

// RenderRectFilled draws a filled black rectangle from (x1,y1) to (x2,y2).
func (be *Backend) RenderRectFilled(x1, y1, x2, y2 float64) {
	//panic("not implemented")
}

// Metrics returns the metrics.
func (be *Backend) Metrics(symbol string, font font.Font, dpi float64, math bool) font.Metrics {
	if dpi != 72 {
		panic(fmt.Errorf("no pre-generated metrics for dpi=%v", dpi))
	}

	key := fontKey{symbol, font, math}
	metrics, ok := be.fonts[key]
	if !ok {
		panic(fmt.Errorf("no pre-generated metrics for %#v", key))
	}

	return metrics
}

// UnderlineThickness returns the line thickness that matches the given font.
// It is used as a base unit for drawing lines such as in a fraction or radical.
func (be *Backend) UnderlineThickness(font font.Font, dpi float64) float64 {
	// theoretically, we could grab the underline thickness from the font
	// metrics.
	// but that information is just too un-reliable.
	// so, it is hardcoded.
	return (0.75 / 12 * font.Size * dpi) / 72
}

// Kern returns the kerning distance between two symbols.
func (be *Backend) Kern(ft1 font.Font, sym1 string, ft2 font.Font, sym2 string, dpi float64) float64 {
	if ft1 != ft2 {
		panic(fmt.Errorf("kern w/ different fonts not supported (ft1=%#v, ft2=%#v, sym1=%q, sym2=%q)", ft1, ft2, sym1, sym2))
	}

	kern, ok := be.kerns[kernKey{ft1, sym1, sym2}]
	if !ok {
		panic(fmt.Errorf(
			"no pre-generated kerning for ft1=%v, sym1=%q, ft2=%v, sym2=%q",
			ft1, sym1, ft2, sym2,
		))
	}
	return kern
}

type fontKey struct {
	symbol string
	font   font.Font
	math   bool
}

type kernKey struct {
	//f1, f2 tex.Font
	font   font.Font
	s1, s2 string
}

var (
	_ font.Backend = (*Backend)(nil)
)
