// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"strings"

	"gioui.org/font/opentype"
	"gioui.org/text"
	lmromanbold "github.com/go-fonts/latin-modern/lmroman10bold"
	lmromanbolditalic "github.com/go-fonts/latin-modern/lmroman10bolditalic"
	lmromanitalic "github.com/go-fonts/latin-modern/lmroman10italic"
	lmromanregular "github.com/go-fonts/latin-modern/lmroman10regular"
	"github.com/go-fonts/liberation/liberationserifbold"
	"github.com/go-fonts/liberation/liberationserifbolditalic"
	"github.com/go-fonts/liberation/liberationserifitalic"
	"github.com/go-fonts/liberation/liberationserifregular"

	"github.com/go-latex/latex/font/liberation"
	"github.com/go-latex/latex/font/lm"
	"github.com/go-latex/latex/font/ttf"
)

func liberationFonts() *ttf.Fonts {
	return liberation.Fonts()
}

func lmromanFonts() *ttf.Fonts {
	return lm.Fonts()
}

func registerFont(fnt text.Font, name string, raw []byte) text.FontFace {
	face, err := opentype.Parse(raw)
	if err != nil {
		log.Fatalf("could not parse fonts: %+v", err)
	}

	if strings.Contains(name, "-") {
		i := strings.Index(name, "-")
		name = name[:i]
	}
	fnt.Typeface = text.Typeface(name)
	return text.FontFace{
		Font: fnt,
		Face: face,
	}
}

func liberationCollection() []text.FontFace {
	var coll []text.FontFace

	coll = append(coll,
		registerFont(
			text.Font{},
			"Liberation",
			liberationserifregular.TTF,
		),
		registerFont(
			text.Font{Weight: text.Bold},
			"Liberation",
			liberationserifbold.TTF,
		),
		registerFont(
			text.Font{Style: text.Italic},
			"Liberation",
			liberationserifitalic.TTF,
		),
		registerFont(
			text.Font{Weight: text.Bold, Style: text.Italic},
			"Liberation",
			liberationserifbolditalic.TTF,
		),
	)
	return coll
}

func latinmodernCollection() []text.FontFace {
	var coll []text.FontFace

	coll = append(coll,
		registerFont(
			text.Font{},
			"LatinModern-Regular",
			lmromanregular.TTF,
		),
		registerFont(
			text.Font{Weight: text.Bold},
			"LatinModern-Bold",
			lmromanbold.TTF,
		),
		registerFont(
			text.Font{Style: text.Italic},
			"LatinModern-Italic",
			lmromanitalic.TTF,
		),
		registerFont(
			text.Font{Weight: text.Bold, Style: text.Italic},
			"LatinModern-BoldItalic",
			lmromanbolditalic.TTF,
		),
	)
	return coll
}
