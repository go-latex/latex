// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package drawimg implements a canvas for img.
package drawimg // import "github.com/go-latex/latex/drawtex/drawimg"

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math"

	"github.com/fogleman/gg"
	"github.com/go-latex/latex/drawtex"
	"github.com/go-latex/latex/mtex"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Renderer struct {
	w io.Writer
}

func NewRenderer(w io.Writer) *Renderer {
	return &Renderer{w: w}
}

func (r *Renderer) Render(width, height, dpi float64, c *drawtex.Canvas) error {
	var (
		w   = width * dpi
		h   = height * dpi
		ctx = gg.NewContext(int(math.Ceil(w)), int(math.Ceil(h)))
	)
	// log.Printf("write: w=%g, h=%g", w, h)

	if false {
		draw.Draw(ctx.Image().(draw.Image), ctx.Image().Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	}

	ctx.SetColor(color.Black)

	for _, op := range c.Ops() {
		switch op := op.(type) {
		case drawtex.GlyphOp:
			drawGlyph(ctx, dpi, op)
		case drawtex.RectOp:
			drawRect(ctx, dpi, op)
		default:
			panic(fmt.Errorf("unknown drawtex op %T", op))
		}
	}

	return png.Encode(r.w, ctx.Image())
}

func drawGlyph(ctx *gg.Context, dpi float64, op drawtex.GlyphOp) {
	face, err := opentype.NewFace(op.Glyph.Font, &opentype.FaceOptions{
		DPI:     dpi,
		Size:    op.Glyph.Size,
		Hinting: font.HintingNone,
	})
	if err != nil {
		panic(fmt.Errorf("could not open font face for glyph %q: %+v",
			op.Glyph.Symbol, err,
		))
	}
	defer face.Close()
	ctx.SetFontFace(face)

	dpi /= 72

	x := op.X * dpi
	y := op.Y * dpi
	//	log.Printf("draw-glyph: %q w=%g, h=%g x=%g, y=%g, size=%v",
	//		op.Glyph.Symbol,
	//		w, h, x, y, op.Glyph.Size,
	//	)
	ctx.DrawString(op.Glyph.Symbol, x, y)
}

func drawRect(ctx *gg.Context, dpi float64, op drawtex.RectOp) {
	dpi /= 72
	ctx.NewSubPath()
	ctx.MoveTo(op.X1*dpi, op.Y1*dpi)
	ctx.LineTo(op.X2*dpi, op.Y1*dpi)
	ctx.LineTo(op.X2*dpi, op.Y2*dpi)
	ctx.LineTo(op.X1*dpi, op.Y2*dpi)
	ctx.LineTo(op.X1*dpi, op.Y1*dpi)
	ctx.ClosePath()
	ctx.Fill()
	//	log.Printf("draw-rect: pt1=(%g, %g) -> (%g, %g)", op.X1, op.Y1, op.X2, op.Y2)
}

var (
	_ mtex.Renderer = (*Renderer)(nil)
)
