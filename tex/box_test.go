// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tex

import (
	"math"
	"testing"

	"github.com/go-latex/latex/font"
	"github.com/go-latex/latex/internal/fakebackend"
)

func TestBox(t *testing.T) {
	const dpi = 72
	be := fakebackend.New()
	state := NewState(be, font.Font{
		Name: "default",
		Size: 12,
		Type: "rm",
	}, dpi)
	for _, tc := range []struct {
		node    Node
		w, h, d float64
	}{
		{
			node: HBox(10),
			w:    10,
			h:    0,
			d:    0,
		},
		{
			node: VBox(10, 20),
			w:    0,
			h:    10,
			d:    20,
		},
		{
			node: HListOf([]Node{VBox(10, 20), HBox(30)}, false),
			w:    30,
			h:    10,
			d:    20,
		},
		{
			node: HListOf([]Node{VBox(10, 20), HBox(30)}, true),
			w:    30,
			h:    10,
			d:    20,
		},
		{
			node: VListOf([]Node{VBox(10, 20), HBox(30)}),
			w:    30,
			h:    30,
			d:    0,
		},
		{
			node: HListOf([]Node{
				VBox(10, 20), HBox(30),
				HListOf([]Node{HBox(11), HBox(22)}, false),
			}, false),
			w: 63,
			h: 10,
			d: 20,
		},
		{
			node: HListOf([]Node{
				VBox(10, 20), HBox(30),
				HListOf([]Node{HBox(11), HBox(22)}, false),
				VListOf([]Node{HBox(15), VBox(11, 22)}),
			}, false),
			w: 78,
			h: 11,
			d: 22,
		},
		{
			node: HListOf([]Node{VBox(10, 20), NewKern(15), HBox(30)}, true),
			w:    45,
			h:    10,
			d:    20,
		},
		{
			node: VListOf([]Node{
				VBox(10, 20),
				VListOf([]Node{
					VBox(11, 22),
					NewKern(10),
					HBox(40),
				}),
				HListOf([]Node{VBox(10, 20), NewKern(15), HBox(30)}, true),
				HBox(30),
			}),
			w: 45,
			h: 103,
			d: 0,
		},
		{
			node: NewKern(10),
			w:    10,
			h:    0,
			d:    0,
		},
		{
			node: NewGlue("fil"),
		},
		{
			node: NewGlue("fill"),
		},
		{
			node: NewGlue("filll"),
		},
		{
			node: NewGlue("neg_fil"),
		},
		{
			node: NewGlue("neg_fill"),
		},
		{
			node: NewGlue("neg_filll"),
		},
		{
			node: NewGlue("empty"),
		},
		{
			node: NewGlue("ss"),
		},
		{
			node: VListOf([]Node{
				NewKern(10),
				VBox(10, 20),
				NewKern(10),
				VListOf([]Node{
					NewKern(10),
					VBox(11, 22),
					NewKern(10),
					HBox(40),
					NewKern(10),
				}),
				NewKern(10),
				HListOf([]Node{
					NewKern(10), VBox(10, 20),
					NewKern(15), HBox(30),
					NewKern(10),
				}, true),
				NewKern(10),
				HBox(30),
				NewKern(10),
			}),
			w: 65,
			h: 173,
			d: 0,
		},
		{
			node: VListOf([]Node{
				NewKern(10),
				VBox(10, 20),
				NewGlue("fill"),
				NewKern(10),
				VListOf([]Node{
					NewKern(10),
					VBox(11, 22),
					NewKern(10),
					NewGlue("neg_fill"),
					HBox(40),
					NewKern(10),
				}),
				NewKern(10),
				HListOf([]Node{
					NewKern(10), VBox(10, 20),
					NewGlue("empty"),
					NewKern(15), HBox(30),
					NewKern(10),
				}, true),
				NewKern(10),
				NewGlue("ss"),
				HBox(30),
				NewKern(10),
			}),
			w: 65,
			h: 173,
			d: 0,
		},
		{
			node: HListOf([]Node{
				NewKern(10),
				NewGlue("fil"),
				VBox(10, 20), HBox(30),
				NewGlue("fil"),
				HListOf([]Node{HBox(11), NewGlue("filll"), HBox(22)}, true),
				VListOf([]Node{HBox(15), NewGlue("neg_filll"), VBox(11, 22)}),
			}, true),
			w: 88,
			h: 11,
			d: 22,
		},
		{
			node: HCentered([]Node{
				VBox(10, 20),
				HBox(30),
				NewKern(15),
				HBox(40),
				VBox(20, 10),
			}),
			w: 85,
			h: 20,
			d: 20,
		},
		{
			node: VCentered([]Node{
				VBox(10, 20),
				HBox(30),
				NewKern(15),
				HBox(40),
				VBox(20, 10),
			}),
			w: 40,
			h: 75,
			d: 0,
		},
		{
			node: NewChar("a", state, false),
			w:    5.53125,
			h:    6.71875,
			d:    0.171875,
		},
		{
			node: NewChar("a", state, true),
			w:    5.53125,
			h:    6.71875,
			d:    0.171875,
		},
		{
			node: NewChar(" ", state, false),
			w:    3.814453125,
			h:    0,
			d:    0,
		},
		{
			node: NewChar(" ", state, true),
			w:    3.814453125,
			h:    0,
			d:    0,
		},
		{
			node: NewChar(`\sigma`, state, true),
			w:    6.578125,
			h:    6.5625,
			d:    0.171875,
		},
		{
			node: NewChar(`\sum`, state, true),
			w:    10.890625,
			h:    12.203125,
			d:    3.265625,
		},
		{
			node: NewChar(`\oint`, state, true),
			w:    6.90625,
			h:    12.84375,
			d:    3.609375,
		},
		{
			node: NewAccent(`é`, state, false),
			w:    6.09375,
			h:    9.765625,
			d:    0,
		},
		{
			node: HRule(state, -1),
			w:    math.Inf(+1),
			h:    0.375,
			d:    0.375,
		},
		{
			node: HRule(state, 10),
			w:    math.Inf(+1),
			h:    5,
			d:    5,
		},
		{
			node: VRule(state),
			w:    0.75,
			h:    math.Inf(+1),
			d:    math.Inf(+1),
		},
		{
			node: HListOf([]Node{
				NewKern(10),
				NewChar("A", state, false), NewChar("V", state, false),
				NewChar("A", state, false), NewChar("V", state, false),
				NewAccent("é", state, false), NewAccent("é", state, false),
				NewChar("A", state, false), NewChar(`\sigma`, state, true),
				NewChar(`\sum`, state, true), NewChar(`\sigma`, state, true),
				NewKern(10),
			}, true),
			w: 102.453125,
			h: 12.203125,
			d: 3.265625,
		},
		{
			node: HListOf([]Node{
				NewKern(10),
				NewChar("A", state, false), NewChar("V", state, false),
				NewChar("A", state, false), NewChar("V", state, false),
				NewAccent("é", state, false), NewAccent("é", state, false),
				NewChar("A", state, false), NewChar(`\sigma`, state, true),
				NewChar(`\sum`, state, true), NewChar(`\sigma`, state, true),
				NewKern(10),
			}, false),
			w: 96.3125,
			h: 12.203125,
			d: 3.265625,
		},
		{
			node: HListOf([]Node{
				NewKern(10),
				NewChar(`\sigma`, state, true),
				HRule(state, -1),
				NewChar(`\sum`, state, true),
				NewKern(10),
			}, true),
			w: math.Inf(+1),
			h: 12.203125,
			d: 3.265625,
		},
		{
			node: HListOf([]Node{
				NewKern(10),
				NewChar(`\sigma`, state, true),
				HRule(state, -1),
				NewChar(`\sum`, state, true),
				NewKern(10),
			}, false),
			w: math.Inf(+1),
			h: 12.203125,
			d: 3.265625,
		},
		{
			node: HListOf([]Node{
				NewKern(10),
				NewChar(`\sigma`, state, true),
				VRule(state),
				NewChar(`\sum`, state, true),
				NewKern(10),
			}, true),
			w: 39.787109375,
			h: 12.203125,
			d: 3.265625,
		},
		{
			node: HListOf([]Node{
				NewKern(10),
				NewChar(`\sigma`, state, true),
				VRule(state),
				NewChar(`\sum`, state, true),
				NewKern(10),
				HListOf(nil, true),
			}, false),
			w: 38.21875,
			h: 12.203125,
			d: 3.265625,
		},
		{
			node: HListOf([]Node{
				NewKern(10),
				NewAccent(`é`, state, false),
				VRule(state),
				NewChar(`\sum`, state, true),
				NewKern(10),
				HListOf(nil, true),
				VListOf(nil),
			}, false),
			w: 37.734375,
			h: 12.203125,
			d: 3.265625,
		},
		{
			node: VListOf([]Node{
				NewKern(10),
				VRule(state),
				HRule(state, 10),
				NewKern(10),
				HListOf(nil, true),
				VListOf(nil),
			}),
			w: 0.75,
			h: math.Inf(+1),
			d: 0,
		},
		{
			node: AutoHeightChar(`(`, 8.865625, 2.5703124999999996, state, 0),
			w:    5.0047149658203125,
			h:    9.734375,
			d:    1.6875,
		},
	} {
		t.Run("", func(t *testing.T) {
			var (
				w = tc.node.Width()
				h = tc.node.Height()
				d = tc.node.Depth()
			)

			if got, want := w, tc.w; !cmpEq(got, want) {
				t.Fatalf("invalid width: got=%g, want=%g", got, want)
			}

			if got, want := h, tc.h; !cmpEq(got, want) {
				t.Fatalf("invalid height: got=%g, want=%g", got, want)
			}

			if got, want := d, tc.d; !cmpEq(got, want) {
				t.Fatalf("invalid depth: got=%g, want=%g", got, want)
			}
		})
	}
}

func TestShip(t *testing.T) {
	const dpi = 72
	be := fakebackend.New()
	state := NewState(be, font.Font{
		Name: "default",
		Size: 12,
		Type: "rm",
	}, dpi)
	for _, tc := range []struct {
		node Tree
		s    int
		v, h float64
	}{
		{
			node: HListOf([]Node{
				NewKern(10),
				VRule(state),
				HListOf(nil, true),
				VListOf(nil),
			}, false),
			s: 0,
			v: 0,
			h: 10.75,
		},
		{
			node: HListOf([]Node{
				NewKern(10),
				NewAccent(`é`, state, false),
				VRule(state),
				NewChar(`\sum`, state, true),
				NewKern(10),
				HListOf(nil, true),
				VListOf(nil),
			}, false),
			s: 0,
			v: 0,
			h: 37.734375,
		},
		{
			node: VListOf([]Node{
				HRule(state, 10),
			}),
			s: 0,
			v: 0,
			h: math.Inf(+1),
		},
		{
			node: VListOf([]Node{
				NewKern(10),
				VRule(state),
				NewKern(10),
				HListOf(nil, true),
				VListOf(nil),
				NewGlue("filll"),
			}),
			s: 0,
			v: 0,
			h: 20.75,
		},
		{
			node: VListOf([]Node{
				NewKern(10),
				VRule(state),
				NewKern(10),
				HListOf(nil, true),
				NewGlue("fil"),
				HCentered([]Node{
					NewChar("a", state, false),
					NewChar("a", state, false),
				}),
				NewGlue("fill"),
				NewGlue("filll"),
				VListOf([]Node{
					NewKern(10),
					NewGlue("neg_fil"),
					NewGlue("neg_fill"),
					NewGlue("neg_filll"),
					HListOf(nil, true),
					HListOf([]Node{
						NewKern(10),
						NewKern(10),
						NewAccent("é", state, false),
						NewChar(`\sigma`, state, true),
						NewChar(`a`, state, false),
						HRule(state, 10),
						VRule(state),
						HListOf([]Node{NewChar("a", state, false)}, true),
						VListOf([]Node{
							VBox(1, 2),
							VBox(2, 3),
							HBox(10),
							VListOf([]Node{
								VBox(1, 2),
								VBox(1, 2),
								HBox(10),
								VListOf(nil),
								HRule(state, 10),
								VRule(state),
							}),
						}),
					}, true),
					NewGlue("filll"),
				}),
			}),
			s: 0,
			v: 0,
			h: 31.8125,
		},
	} {
		t.Run("", func(t *testing.T) {
			var ship Ship
			ship.Call(0, 0, tc.node)
			var (
				s = ship.cur.s
				v = ship.cur.v
				h = ship.cur.h
			)

			if got, want := s, tc.s; got != want {
				t.Fatalf("invalid shift: got=%d, want=%d", got, want)
			}

			if got, want := v, tc.v; !cmpEq(got, want) {
				t.Fatalf("invalid v: got=%g, want=%g", got, want)
			}

			if got, want := h, tc.h; !cmpEq(got, want) {
				t.Fatalf("invalid h: got=%g, want=%g", got, want)
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
