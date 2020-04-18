// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tex

type State struct {
	be   Backend
	Font Font
	DPI  float64
}

func NewState(be Backend, font Font, dpi float64) State {
	return State{
		be:   be,
		Font: font,
		DPI:  dpi,
	}
}

func (state State) Backend() Backend { return state.be }

type Font struct {
	Name string
	Type string
	Size float64
}
