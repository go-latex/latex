// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ttf

import (
	"fmt"
	"testing"

	"github.com/go-fonts/dejavu/dejavusans"
	"github.com/go-fonts/dejavu/dejavusansoblique"
	"github.com/go-latex/latex/font"
	"github.com/go-latex/latex/internal/fakebackend"
	"golang.org/x/image/font/sfnt"
)

func TestDejaVuBackend(t *testing.T) {
	var (
		be  = newBackend()
		ref = fakebackend.New()
	)
	for _, sym := range []string{
		"A",
		// "B",
		// "a",
		// "g",
		"z",
		"Z",
		// "I",
		"T",
		"i",
		"t",
		// `\sum`,
		// `\sigma`,
	} {
		for _, math := range []bool{
			true,
			false,
		} {
			for _, descr := range []font.Font{
				{Name: "default", Size: 12, Type: "rm"},
				//{Name: "default", Size: 10, Type: "rm"},
				//{Name: "it", Size: 12, Type: "it"},
				//{Name: "it", Size: 10, Type: "it"},
			} {
				t.Run(fmt.Sprintf("Metrics/%s-math=%v-%s-%g-%s", sym, math, descr.Name, descr.Size, descr.Type), func(t *testing.T) {
					got := be.Metrics(sym, descr, 72, math)
					if got, want := got, ref.Metrics(sym, descr, 72, math); got != want {
						t.Fatalf("invalid metrics.\ngot= %#v\nwant=%#v\n", got, want)
					}
				})
				t.Run(fmt.Sprintf("XHeight/%s-math=%v-%s-%g-%s", sym, math, descr.Name, descr.Size, descr.Type), func(t *testing.T) {
					got := be.XHeight(descr, 72)
					if got, want := got, ref.XHeight(descr, 72); got != want {
						t.Fatalf("invalid xheight.\ngot= %#v\nwant=%#v\n", got, want)
					}
				})
			}
		}
	}
}

func newBackend() *Backend {
	be := &Backend{
		glyphs: make(map[ttfKey]ttfVal),
		fonts:  make(map[string]*sfnt.Font),
	}

	ftmap := map[string][]byte{
		"default": dejavusans.TTF,
		"regular": dejavusans.TTF,
		"rm":      dejavusans.TTF,
		"it":      dejavusansoblique.TTF,
	}
	for k, raw := range ftmap {
		ft, err := sfnt.Parse(raw)
		if err != nil {
			panic(fmt.Errorf("could not parse %q: %+v", k, err))
		}
		be.fonts[k] = ft
	}

	return be
}

func TestGofontBackend(t *testing.T) {
	be := New(nil)
	{
		fnt := font.Font{Name: "default", Size: 12, Type: "regular"}
		got := be.Metrics("A", fnt, 72, true)
		want := font.Metrics{Advance: 8.00390625, Height: 8.671875, Width: 7.75, XMin: 0.109375, XMax: 7.859375, YMin: 0, YMax: 8.671875, Iceberg: 8.671875, Slanted: false}
		if got != want {
			t.Fatalf("got=%#v\nwant=%#v", got, want)
		}
	}
	{
		fnt := font.Font{Name: "it", Size: 12, Type: "it"}
		got := be.Metrics("A", fnt, 72, true)
		want := font.Metrics{Advance: 8.1328125, Height: 8.671875, Width: 7.75, XMin: 0.171875, XMax: 7.921875, YMin: 0, YMax: 8.671875, Iceberg: 8.671875, Slanted: true}
		if got != want {
			t.Fatalf("got=%#v\nwant=%#v", got, want)
		}
	}
}
