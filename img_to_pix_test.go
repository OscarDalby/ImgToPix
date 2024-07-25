package main

import (
	"image/color"
	"testing"
)

func TestAverageColor(t *testing.T) {
	result := average_color(
		color.RGBA{1, 1, 1, 1},
		color.RGBA{2, 2, 2, 2},
		color.RGBA{3, 3, 3, 3},
		color.RGBA{4, 4, 4, 4},
		color.RGBA{5, 5, 5, 5},
	)
	expected := color.RGBA{3, 3, 3, 3}
	if result != expected {
		t.Errorf("average_color(1,2,3,4,5) = %d; want %v", result, expected)
	}
}
