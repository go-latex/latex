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

type dbXHs map[xhKey]float64
type dbFonts map[fontKey]font.Metrics
type dbKerns map[kernKey]float64

var (
	fontsDb dbFonts
	xhsDb   dbXHs
	kernsDb dbKerns
)

type Backend struct {
	fonts dbFonts
	xhs   dbXHs
	kerns dbKerns
}

func New() *Backend {
	return &Backend{
		fonts: fontsDb,
		xhs:   xhsDb,
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

// XHeight returns the xheight for the given font and dpi.
func (be *Backend) XHeight(font font.Font, dpi float64) float64 {
	key := xhKey{font.Name, font.Size, dpi}
	xh, ok := be.xhs[key]
	if !ok {
		panic(fmt.Errorf("no pre-generated xheight for %#v", key))
	}

	return xh
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
	if ft1 == ft2 {
		return 0
	}

	kern, ok := be.kerns[kernKey{ft1, sym1, sym2}]
	if !ok {
		return 0
		//	panic(fmt.Errorf(
		//		"no pre-generated kerning for ft1=%v, sym1=%q, ft2=%v, sym2=%q",
		//		ft1, sym1, ft2, sym2,
		//	))
	}
	return kern
}

type fontKey struct {
	symbol string
	font   font.Font
	math   bool
}

type xhKey struct {
	font string
	size float64
	dpi  float64
}

type kernKey struct {
	//f1, f2 tex.Font
	font   font.Font
	s1, s2 string
}

var (
	_ font.Backend = (*Backend)(nil)
)
