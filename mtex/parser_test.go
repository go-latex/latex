// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mtex

import (
	"math"
	"testing"

	"github.com/go-latex/latex/internal/fakebackend"
)

func TestParse(t *testing.T) {
	const (
		dpi    = 72
		ftsize = 10
	)
	var (
		be = fakebackend.New()
	)
	for _, tc := range []struct {
		expr    string
		w, h, d float64
	}{
		{
			expr: "hello",
			w:    24.1650390625,
			h:    7.59375,
			d:    0.140625,
		},
		{
			expr: "$hello$",
			w:    24.1650390625,
			h:    7.59375,
			d:    0.140625,
		},
		{
			expr: `$\sigma$`,
			w:    6.337890625,
			h:    5.46875,
			d:    0.140625,
		},
		{
			expr: `$\sigma$ is $12$`,
			w:    33.408203125,
			h:    7.59375,
			d:    0.140625,
		},
		{
			expr: `$1.1$`,
			w:    15.9033203125,
			h:    7.296875,
			d:    0.0,
		},
		{
			expr: `$1.$`,
			w:    11.4892578125,
			h:    7.296875,
			d:    0.0,
		},
		{
			expr: `$.2$`,
			w:    11.4892578125,
			h:    7.421875,
			d:    0.0,
		},
		{
			expr: `$.$`,
			w:    5.126953125,
			h:    1.234375,
			d:    0.0,
		},
		{
			expr: `$x.x$`,
			w:    16.962890625,
			h:    5.46875,
			d:    0.0,
		},
		//{ // FIXME(sbinet): handler for '('
		//	expr: `$\sigma = f(x)$`,
		//	w:    35.8544921875,
		//	h:    7.59375,
		//	d:    1.3125,
		//},
		{
			expr: `$\sigma \rightarrow \infty$`,
			w:    26.943359375,
			h:    5.46875,
			d:    0.140625,
		},
		{
			expr: `$\sigma\,=\infty$`,
			w:    28.566927001953125,
			h:    5.46875,
			d:    0.140625,
		},
		{
			expr: `$\sigma\hspace{2}=\infty$`,
			w:    46.42578125,
			h:    5.46875,
			d:    0.140625,
		},
		{
			expr: `$\cos\theta$`,
			w:    22.9443359375,
			h:    7.671875,
			d:    0.140625,
		},
		{
			expr: `$\pi$`,
			w:    6.0205078125,
			h:    5.46875,
			d:    0.1875,
		},
		{
			expr: `$\frac{x}{y}$`,
			w:    5.392578125,
			h:    8.2109375,
			d:    4.0249999999999995,
		},
		{
			expr: `$\frac{1}{2}$`,
			w:    5.70361328125,
			h:    9.490625,
			d:    3.9375,
		},
		{
			expr: `$\frac{1}{2\pi}$`,
			w:    9.91796875,
			h:    9.490625,
			d:    4.06875,
		},
		{
			expr: `$\dfrac{1}{2}$`,
			w:    7.6123046875,
			h:    11.6796875,
			d:    6.1640625,
		},
		{
			expr: `$\tfrac{1}{2}$`,
			w:    5.70361328125,
			h:    9.490625,
			d:    3.9375,
		},
		{
			expr: `$\binom{1}{x}$`,
			w:    15.713043212890625, // w/o cm-fallback
			h:    8.865625,
			d:    2.5703124999999996,
		},
		{
			expr: `$\sqrt{2x}$`,
			w:    22.864837646484375,
			h:    11.146875,
			d:    0,
		},
		{
			// FIXME(sbinet): check values somehow.
			// but... matplotlib.mathtex doesn't handle `$\sqrt[3]{x}$`
			expr: `$\sqrt[3]{2x}$`,
			w:    21.940084838867186,
			h:    11.146875,
			d:    0,
		},
		{
			expr: `$\overline{ab}$`,
			w:    12.4755859375,
			h:    10.06875,
			d:    0.140625,
		},
	} {
		t.Run("", func(t *testing.T) {
			defer func() {
				err := recover()
				if err != nil {
					t.Errorf("%q: panic: %+v", tc.expr, err)
					panic(err)
				}
			}()
			got, err := Parse(tc.expr, ftsize, dpi, be)
			if err != nil {
				t.Fatalf("could not parse %q: %+v", tc.expr, err)
			}

			var (
				w = got.Width()
				h = got.Height()
				d = got.Depth()
			)

			if got, want := w, tc.w; got != want {
				t.Fatalf("%q: invalid width: got=%g, want=%g", tc.expr, got, want)
			}

			if got, want := h, tc.h; !cmpEq(got, want) {
				t.Fatalf("%q: invalid height: got=%g, want=%g", tc.expr, got, want)
			}

			if got, want := d, tc.d; !cmpEq(got, want) {
				t.Fatalf("%q: invalid depth: got=%g, want=%g", tc.expr, got, want)
			}
		})
	}
}

func cmpEq(a, b float64) bool {
	switch {
	case math.IsInf(a, -1):
		return math.IsInf(b, -1)
	case math.IsInf(a, +1):
		return math.IsInf(b, +1)
	default:
		return a == b
	}
}
