// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package drawimg_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/go-latex/latex/drawtex/drawimg"
	"github.com/go-latex/latex/mtex"
)

func TestRenderer(t *testing.T) {
	const (
		size = 12
		dpi  = 256
	)

	load := func(name string) []byte {
		name = "testdata/" + name + "_golden.png"
		raw, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatalf("could not read file %q: %+v", name, err)
		}
		return raw
	}

	for _, tc := range []struct {
		name string
		expr string
	}{
		{
			name: "sqrt",
			expr: `$\sqrt{x}$`,
		},
		{
			name: "sqrt_over_2pi",
			expr: `$\frac{\sqrt{x+20}}{2\pi}$`,
		},
		{
			name: "delta",
			expr: `$\delta x \neq \frac{\sqrt{x+20}}{2\Delta}$`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			dst := drawimg.NewRenderer(out)
			err := mtex.Render(dst, tc.expr, size, dpi)
			if err != nil {
				t.Fatalf("could not render expression %q: %+v", tc.expr, err)
			}

			if got, want := out.Bytes(), load(tc.name); !bytes.Equal(got, want) {
				err := ioutil.WriteFile("testdata/"+tc.name+".png", got, 0644)
				if err != nil {
					t.Fatalf("could not create output file: %+v", err)
				}
				t.Fatal("files differ")
			}
		})
	}
}
