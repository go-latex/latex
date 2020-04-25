// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tex

import "testing"

func TestDetermineOrder(t *testing.T) {
	for _, tc := range []struct {
		totals []float64
		want   int
	}{
		{
			totals: []float64{1, 2, 3, 0},
			want:   2,
		},
		{
			totals: []float64{1, 2, 3, 4},
			want:   3,
		},
		{
			totals: []float64{0, 2, 3, 0},
			want:   2,
		},
		{
			totals: []float64{0, 1, 0, 0},
			want:   1,
		},
	} {
		t.Run("", func(t *testing.T) {
			got := determineOrder(tc.totals)
			if got != tc.want {
				t.Fatalf("%v: got=%v, want=%v", tc.totals, got, tc.want)
			}
		})
	}
}
