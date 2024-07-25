package main

import (
	"image/color"
	"testing"
)

func TestAverageColor(t *testing.T) {
	color_list := []color.RGBA{
		{1, 1, 1, 1},
		{2, 2, 2, 2},
		{3, 3, 3, 3},
		{4, 4, 4, 4},
		{5, 5, 5, 5},
	}
	result := average_color(color_list)
	expected := color.RGBA{3, 3, 3, 3}
	if result != expected {
		t.Errorf("average_color(1,2,3,4,5) = %d; want %v", result, expected)
	}
}
