// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package drawpdf implements a canvas for PDF.
package drawpdf // import "github.com/go-latex/latex/drawtex/drawpdf"

import (
	"log"

	"github.com/go-latex/latex/drawtex"
	pdf "github.com/go-pdf/fpdf"
)

func Write(fname string, w, h float64, c *drawtex.Canvas) error {
	doc := pdf.NewCustom(&pdf.InitType{
		UnitStr: "pt",
		Size:    pdf.SizeType{Wd: w, Ht: h},
	})
	doc.AddPage()

	for _, op := range c.Ops() {
		switch op := op.(type) {
		case drawtex.GlyphOp:
			log.Printf(">>> %T: %#v", op, op)
			drawGlyph(doc, op)
		case drawtex.RectOp:
			log.Printf(">>> %T: %#v", op, op)
			drawRect(doc, op)
		default:
			log.Panicf("unknown drawtex op %T", op)
		}
	}
	return doc.OutputFileAndClose(fname)
}

func drawGlyph(doc *pdf.Fpdf, op drawtex.GlyphOp) {}
func drawRect(doc *pdf.Fpdf, op drawtex.RectOp)   {}
