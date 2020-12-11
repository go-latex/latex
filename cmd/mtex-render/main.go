// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command mtex-render renders a LaTeX math expression to a PNG, PDF, ... file.
//
// Example:
//
//  $> mtex-render "\$\\sqrt{x}\$"
package main

import (
	"flag"
	"log"
	"os"

	"github.com/go-latex/latex/drawtex/drawimg"
	"github.com/go-latex/latex/mtex"
)

func main() {

	log.SetPrefix("mtex: ")
	log.SetFlags(0)

	var (
		dpi  = flag.Float64("dpi", 72, "dots-per-inch to use")
		size = flag.Float64("font-size", 12, "font size to use")
		out  = flag.String("o", "out.png", "path to output file")
	)

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		log.Fatalf("missing math expression to render")
	}

	expr := flag.Arg(0)

	log.Printf("rendering math expression: %q", expr)

	f, err := os.Create(*out)
	if err != nil {
		log.Fatalf("could not create output file: %+v", err)
	}
	defer f.Close()

	dst := drawimg.NewRenderer(f)
	err = mtex.Render(dst, expr, *size, *dpi)
	if err != nil {
		log.Fatalf("could not render math expression %q: %+v", expr, err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("could not close output file: %+v", err)
	}
}
