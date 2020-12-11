// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mtex

import (
	"fmt"
	"testing"

	"github.com/go-latex/latex/drawtex"
)

type dummyRenderer struct{}

func (dummyRenderer) Render(width, height, dpi float64, c *drawtex.Canvas) error {
	for _, op := range c.Ops() {
		switch op.(type) {
		case drawtex.GlyphOp:
		case drawtex.RectOp:
		default:
			panic(fmt.Errorf("unknown drawtex operation %T", op))
		}
	}
	return nil
}

func TestRender(t *testing.T) {
	const (
		dpi    = 72
		ftsize = 10
	)
	for _, tc := range []struct {
		expr string
		want error
	}{
		{
			expr: `math $x= 42$`,
		},
		{
			expr: `math $\sum\sqrt{\frac{a+b}{2\pi}}\cos\Phi$`,
		},
		{
			expr: `math: $\sum\sqrt{\frac{a+b}{2\pi}}\cos\omega\binom{a+b}{\beta}\prod \alpha x$`,
		},
		{
			expr: `$\int\frac{\partial x}{x}$`,
		},
	} {
		t.Run(tc.expr, func(t *testing.T) {
			err := Render(dummyRenderer{}, tc.expr, ftsize, dpi, nil)

			switch {
			case err != nil && tc.want != nil:
				if got, want := err.Error(), tc.want.Error(); got != want {
					t.Fatalf("invalid error:\ngot= %v\nwant=%v", got, want)
				}
			case err == nil && tc.want != nil:
				t.Fatalf("expected an error: %v", tc.want)
			case err != nil && tc.want == nil:
				t.Fatalf("could not render: %v", err)
			case err == nil && tc.want == nil:
				// ok.
			}
		})
	}
}
