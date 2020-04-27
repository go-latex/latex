// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fakebackend

import (
	"testing"

	"github.com/go-latex/latex/font"
)

func TestBackend(t *testing.T) {
	be := New()
	{
		fnt := font.Font{Name: "default", Size: 12, Type: "regular"}
		got := be.Metrics("A", fnt, 72, true)
		want := font.Metrics{Advance: 8.208984375, Height: 8.75, Width: 8.015625, XMin: 0.09375, XMax: 8.109375, YMin: 0, YMax: 8.75, Iceberg: 8.75, Slanted: false}
		if got != want {
			t.Fatalf("got=%#v\nwant=%#v", got, want)
		}
	}
	{
		fnt := font.Font{Name: "it", Size: 12, Type: "it"}
		got := be.Metrics("A", fnt, 72, true)
		want := font.Metrics{Advance: 8.208984375, Height: 8.75, Width: 8.015625, XMin: -0.640625, XMax: 7.390625, YMin: 0, YMax: 8.75, Iceberg: 8.75, Slanted: true}
		if got != want {
			t.Fatalf("got=%#v\nwant=%#v", got, want)
		}
	}
	{
		fnt := font.Font{Name: "default", Size: 12, Type: "regular"}
		got := be.Metrics(`\oint`, fnt, 72, true)
		want := font.Metrics{Advance: 8.8359375, Height: 16.453125, Width: 6.90625, XMin: 0.96875, XMax: 7.875, YMin: -3.609375, YMax: 12.84375, Iceberg: 12.84375, Slanted: true}
		if got != want {
			t.Fatalf("got=%#v\nwant=%#v", got, want)
		}
	}

	{
		fnt := font.Font{Name: "default", Size: 12, Type: "regular"}
		got := be.XHeight(fnt, 72)
		want := 6.5625
		if got != want {
			t.Fatalf("got=%#v\nwant=%#v", got, want)
		}
	}
}
