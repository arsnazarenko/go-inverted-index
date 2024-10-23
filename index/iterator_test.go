package index

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntersectSlices(t *testing.T) {
	cases := []struct {
		s1  []int
		s2  []int
		res []int
	}{
		{
			s1:  []int{1, 2, 3, 4},
			s2:  []int{1, 4},
			res: []int{1, 4},
		},
		{
			s1:  []int{1, 4},
			s2:  []int{1, 2, 3, 4},
			res: []int{1, 4},
		},
		{
			s1:  []int{1, 2, 4, 5, 6},
			s2:  []int{0, 9, 21},
			res: []int{},
		},
		{
			s1:  []int{1, 2, 4, 6, 7, 12, 123, 1222},
			s2:  []int{},
			res: []int{},
		},
		{
			s1:  []int{100, 200, 300, 400},
			s2:  []int{100, 200, 300, 400},
			res: []int{100, 200, 300, 400},
		},
	}
	for i, tt := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i), func(t *testing.T) {
			actual := intersectSlices(tt.s1, tt.s2)
			require.Equal(t, tt.res, actual)

		})
	}
}

func TestMergeSlices(t *testing.T) {
	cases := []struct {
		s1  []int
		s2  []int
		res []int
	}{
		{
			s1:  []int{1, 2, 3, 4},
			s2:  []int{1, 4},
			res: []int{1, 2, 3, 4},
		},
		{
			s1:  []int{1, 4},
			s2:  []int{1, 2, 3, 4},
			res: []int{1, 2, 3, 4},
		},
		{
			s1:  []int{1, 2, 4, 5, 6},
			s2:  []int{0, 9, 21},
			res: []int{0, 1, 2, 4, 5, 6, 9, 21},
		},
		{
			s1:  []int{1, 2, 4, 6, 7, 12, 123, 1222},
			s2:  []int{},
			res: []int{1, 2, 4, 6, 7, 12, 123, 1222},
		},
		{
			s1:  []int{100, 200, 300, 400},
			s2:  []int{100, 200, 300, 400},
			res: []int{100, 200, 300, 400},
		},
	}
	for i, tt := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i), func(t *testing.T) {
			actual := mergeSlices(tt.s1, tt.s2)
			require.Equal(t, tt.res, actual)

		})
	}
}

func TestSubstractSlices(t *testing.T) {
	cases := []struct {
		s1  []int
		s2  []int
		res []int
	}{
		{
			s1:  []int{1, 2, 3, 4},
			s2:  []int{1, 4},
			res: []int{2, 3},
		},
		{
			s1:  []int{1, 4},
			s2:  []int{1, 2, 3, 4},
			res: []int{2, 3},
		},
		{
			s1:  []int{1, 2, 4, 5, 6},
			s2:  []int{0, 9, 21},
			res: []int{1, 2, 4, 5, 6},
		},
		{
			s1:  []int{1, 2, 4, 6, 7, 12, 123, 1222},
			s2:  []int{},
			res: []int{1, 2, 4, 6, 7, 12, 123, 1222},
		},
		{
			s1:  []int{100, 200, 300, 400},
			s2:  []int{100, 200, 300, 400},
			res: []int{},
		},
	}
	for i, tt := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i), func(t *testing.T) {
			actual := subtractSlices(tt.s1, tt.s2)
			require.Equal(t, tt.res, actual)

		})
	}
}
