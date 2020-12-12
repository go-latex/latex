// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/app/headless"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"

	"github.com/go-latex/latex/drawtex"
	"github.com/go-latex/latex/drawtex/drawimg"
	"github.com/go-latex/latex/font/ttf"
	"github.com/go-latex/latex/mtex"
)

const useLiberation = true

func runGio() {
	ui := NewUI()
	go func() {
		win := app.NewWindow(
			app.Title("mtex"),
			app.Size(unit.Dp(ui.width), unit.Dp(ui.height)),
		)
		err := ui.Run(win)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

type UI struct {
	Theme *material.Theme

	Button widget.Clickable
	Editor widget.Editor

	expr string

	width  float32
	height float32
	screen bool

	Image image.Image
}

func NewUI() *UI {
	ui := &UI{
		Theme:  material.NewTheme(gofont.Collection()),
		Image:  image.NewRGBA(image.Rect(0, 0, 1, 1)),
		width:  800,
		height: 900,
	}
	ui.expr = `$\sqrt{x + y}$`
	ui.expr = `$f(x) = \frac{\sqrt{x +20}}{2\pi} +\hbar \sum y\partial y$`
	ui.Editor.SetText(ui.expr)
	return ui
}

func (ui *UI) Run(win *app.Window) error {
	defer win.Close()

	var ops op.Ops
	for e := range win.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case key.Event:
			switch e.Name {
			case key.NameEscape:
				return nil
			case key.NameEnter, key.NameReturn:
				if e.Modifiers.Contain(key.ModCtrl) && e.State == key.Press {
					ui.expr = ui.Editor.Text()
					win.Invalidate()
				}
			case "F11":
				if e.State == key.Press {
					ui.screen = true
					win.Invalidate()
				}
			}

		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			ui.Layout(gtx)
			e.Frame(gtx.Ops)
			if ui.screen {
				ui.screen = false
				ui.screenshot(gtx.Ops)
			}
		}
	}

	return nil
}

var (
	margin = unit.Dp(10)
	list   = &layout.List{
		Axis: layout.Vertical,
	}
)

type (
	D = layout.Dimensions
	C = layout.Context
)

func (ui *UI) screenshot(ops *op.Ops) {
	win, err := headless.NewWindow(int(ui.width), int(ui.height))
	if err != nil {
		return
	}

	err = win.Frame(ops)
	if err != nil {
		return
	}

	img, err := win.Screenshot()
	if err != nil {
		return
	}

	f, err := os.Create("ooo.png")
	if err != nil {
		return
	}
	defer f.Close()

	_ = png.Encode(f, img)
}

func (ui *UI) Layout(gtx C) D {
	widgets := []layout.Widget{
		material.H3(ui.Theme, "Math-TeX renderer").Layout,
		func(gtx C) D {
			gtx.Constraints.Max.Y = gtx.Px(unit.Dp(200))
			ed := material.Editor(ui.Theme, &ui.Editor, "")
			ed.TextSize = ed.TextSize.Scale(1.5)
			return widget.Border{
				Color:        color.NRGBA{A: 107},
				CornerRadius: unit.Dp(4),
				Width:        unit.Dp(2),
			}.Layout(gtx, ed.Layout)
		},
		func(gtx C) D {
			return layout.UniformInset(margin).Layout(gtx, func(gtx C) D {
				for range ui.Button.Clicks() {
					ui.expr = ui.Editor.Text()
				}
				return material.Button(ui.Theme, &ui.Button, "Render").Layout(gtx)
			})
		},
		material.H5(ui.Theme, "Img renderer").Layout,
		func(gtx C) D {
			return layout.UniformInset(margin).Layout(gtx, func(gtx C) D {
				_ = ui.imgRender()
				return widget.Border{
					Color:        color.NRGBA{A: 107},
					CornerRadius: unit.Dp(4),
					Width:        unit.Dp(2),
				}.Layout(
					gtx,
					widget.Image{
						Src: paint.NewImageOp(ui.Image),
					}.Layout,
				)
			})
		},
		material.H5(ui.Theme, "Gio renderer").Layout,
		func(gtx C) D {
			return layout.UniformInset(margin).Layout(gtx, func(gtx C) D {
				_ = ui.render(gtx)
				max := gtx.Constraints.Constrain(gtx.Constraints.Max)
				return D{
					Size: max,
				}
			})
		},
	}

	return list.Layout(gtx, len(widgets), func(gtx C, i int) D {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, widgets[i])
	})
}

func (ui *UI) render(gtx C) error {
	return ui.gioRender(gtx)
}

func (ui *UI) imgRender() error {
	const dpi = 256
	o := new(bytes.Buffer)
	dst := drawimg.NewRenderer(o)
	fnt := lmromanFonts()
	if useLiberation {
		fnt = liberationFonts()
	}
	err := mtex.Render(
		dst,
		ui.expr,
		float64(ui.Theme.TextSize.V),
		dpi,
		fnt,
	)
	if err != nil {
		return err
	}
	ui.Image, err = png.Decode(o)
	if err != nil {
		return err
	}
	return nil
}

func (ui *UI) gioRender(gtx C) error {
	const dpi = 256
	dst := newGioRenderer(gtx, ui.Theme, ui.expr)
	err := mtex.Render(
		dst,
		ui.expr,
		float64(ui.Theme.TextSize.V),
		dpi,
		dst.ft,
	)
	if err != nil {
		return err
	}
	return nil
}

type gioRenderer struct {
	gtx layout.Context
	col color.NRGBA
	th  *material.Theme
	ft  *ttf.Fonts

	offset f32.Point
}

func newGioRenderer(gtx C, th *material.Theme, txt string) *gioRenderer {
	fonts := latinmodernCollection()
	r := gioRenderer{
		gtx: gtx,
		col: color.NRGBA{A: 255},
		th:  material.NewTheme(fonts),
		ft:  lmromanFonts(),
	}
	if useLiberation {
		fonts = liberationCollection()
		r.th = material.NewTheme(fonts)
		r.ft = liberationFonts()
	}
	r.th.TextSize = th.TextSize

	ppem := fixed.Int26_6(r.ft.Rm.UnitsPerEm())

	met, err := r.ft.Rm.Metrics(new(sfnt.Buffer), ppem, font.HintingNone)
	if err != nil {
		panic(fmt.Errorf("could not extract font extents: %+v", err))
	}
	scale := float32(th.TextSize.V) / float32(ppem)

	r.offset = f32.Pt(0, -scale*float32(met.Height-met.Ascent))
	return &r
}

func (r *gioRenderer) Render(width, height, dpi float64, c *drawtex.Canvas) error {

	if false {
		stk := op.Save(r.gtx.Ops)
		var p clip.Path
		p.Begin(r.gtx.Ops)
		p.MoveTo(f32.Point{})
		p.LineTo(r.pt(width*dpi, 0))
		p.LineTo(r.pt(width*dpi, height*dpi))
		p.LineTo(r.pt(0, height*dpi))
		p.Close()
		clip.Stroke{
			Path: p.End(),
			Style: clip.StrokeStyle{
				Width: 2,
			},
		}.Op().Add(r.gtx.Ops)
		paint.Fill(r.gtx.Ops, color.NRGBA{R: 255, A: 255})
		stk.Load()
	}

	dpi /= 72
	for _, opTex := range c.Ops() {
		switch opTex := opTex.(type) {
		case drawtex.GlyphOp:
			r.drawGlyph(dpi, opTex)
		case drawtex.RectOp:
			r.drawRect(dpi, opTex)
		default:
			panic(fmt.Errorf("unknown drawtex op %T", opTex))
		}
	}
	return nil
}

func (r *gioRenderer) drawGlyph(dpi float64, tex drawtex.GlyphOp) {
	defer op.Save(r.gtx.Ops).Load()
	x := tex.X * dpi
	y := (tex.Y - tex.Glyph.Size) * dpi
	op.Offset(r.pt(x, y)).Add(r.gtx.Ops)
	lbl := material.Label(
		r.th,
		unit.Px(float32(tex.Glyph.Size*dpi)),
		tex.Glyph.Symbol,
	)
	if tex.Glyph.Metrics.Slanted {
		lbl.Font.Style = text.Italic
	}
	lbl.Color = r.col
	lbl.Alignment = text.Start
	lbl.Layout(r.gtx)
}

func (r *gioRenderer) drawRect(dpi float64, tex drawtex.RectOp) {
	defer op.Save(r.gtx.Ops).Load()

	var p clip.Path
	p.Begin(r.gtx.Ops)
	p.MoveTo(r.pt(tex.X1*dpi, tex.Y1*dpi).Add(r.offset))
	p.LineTo(r.pt(tex.X2*dpi, tex.Y1*dpi).Add(r.offset))
	p.LineTo(r.pt(tex.X2*dpi, tex.Y2*dpi).Add(r.offset))
	p.LineTo(r.pt(tex.X1*dpi, tex.Y2*dpi).Add(r.offset))
	p.Close()
	clip.Outline{
		Path: p.End(),
	}.Op().Add(r.gtx.Ops)

	paint.Fill(r.gtx.Ops, r.col)
}

func (*gioRenderer) pt(x, y float64) f32.Point {
	return f32.Point{
		X: float32(x),
		Y: float32(y),
	}
}

var (
	_ mtex.Renderer = (*gioRenderer)(nil)
)
