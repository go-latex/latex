// Copyright ©2021 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package lm provides a ttf.Fonts value populated with latin-modern,
// a LaTeX-looking font.
package lm // import "github.com/go-latex/latex/font/lm"

import (
	"log"
	"sync"

	lmromanbold "github.com/go-fonts/latin-modern/lmroman10bold"
	lmromanbolditalic "github.com/go-fonts/latin-modern/lmroman10bolditalic"
	lmromanitalic "github.com/go-fonts/latin-modern/lmroman10italic"
	lmromanregular "github.com/go-fonts/latin-modern/lmroman10regular"
	"golang.org/x/image/font/sfnt"

	"github.com/go-latex/latex/font/ttf"
)

var (
	once sync.Once
	fnts *ttf.Fonts
)

// Fonts returns a ttf.Fonts value populated with latin-modern fonts.
func Fonts() *ttf.Fonts {
	once.Do(func() {
		rm, err := sfnt.Parse(lmromanregular.TTF)
		if err != nil {
			log.Panicf("could not parse fonts: %+v", err)
		}

		it, err := sfnt.Parse(lmromanitalic.TTF)
		if err != nil {
			log.Panicf("could not parse fonts: %+v", err)
		}

		bf, err := sfnt.Parse(lmromanbold.TTF)
		if err != nil {
			log.Panicf("could not parse fonts: %+v", err)
		}

		bfit, err := sfnt.Parse(lmromanbolditalic.TTF)
		if err != nil {
			log.Panicf("could not parse fonts: %+v", err)
		}

		fnts = &ttf.Fonts{
			Default: rm,
			Rm:      rm,
			It:      it,
			Bf:      bf,
			BfIt:    bfit,
		}
	})

	return fnts
}
