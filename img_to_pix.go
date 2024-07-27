package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

type ProcessConfig struct {
	pixel_size int
	scaling    int
	bg_color   color.RGBA
}

var input_path = "./input"
var input_filename = "input"
var tmp_path = "./tmp"
var tmp_filename = "working"
var output_path = "./output"
var output_filename = "output"

var config = ProcessConfig{pixel_size: 32, scaling: 3, bg_color: color.RGBA{0, 0, 0, 0}}

func main() {
	fmt.Printf("image processing started, working image stored in %s/%s\n", tmp_path, tmp_filename)
	var working_image image.Image
	var err error
	base_img, err := get_base_image(input_path, input_filename, config)
	if err != nil {
		log.Fatalf("Failed to get base image: %v", err)
	}
	working_image, err = process_png_pixelise(base_img, config)
	if err != nil {
		log.Fatalf("Failed to process PNG: %v", err)
	}

	create_png(working_image, output_path, "pixelated")

	working_image, err = process_png_apply_palette(working_image, config)
	if err != nil {
		log.Fatalf("Failed to process PNG: %v", err)
	}
	create_png(working_image, output_path, output_filename)
}

func get_base_image(base_path string, file_name string, config ProcessConfig) (image.Image, error) {
	var path = fmt.Sprintf("%s/%s.png", base_path, file_name)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()
	base_img, err := png.Decode(file)
	if base_img.Bounds().Dx()%config.pixel_size != 0 {
		panic("pixel_size must divide width")
	}
	if base_img.Bounds().Dy()%config.pixel_size != 0 {
		panic("pixel_size must divide height")
	}
	if err != nil {
		fmt.Println("Error decoding file:", err)
		return nil, err
	}
	return base_img, nil
}

func get_average_color(colors []color.RGBA) color.RGBA {
	total_r := 0
	total_g := 0
	total_b := 0
	total_a := 0
	num_colors := len(colors)
	if num_colors == 0 {
		return color.RGBA{0, 0, 0, 0}
	}
	for _, color := range colors {
		total_r += int(color.R)
		total_g += int(color.G)
		total_b += int(color.B)
		total_a += int(color.A)
	}
	average_r := total_r / num_colors
	average_g := total_g / num_colors
	average_b := total_b / num_colors
	average_a := total_a / num_colors
	return color.RGBA{uint8(average_r), uint8(average_g), uint8(average_b), uint8(average_a)}
}

func process_png_pixelise(base_img image.Image, config ProcessConfig) (image.Image, error) {
	var width = base_img.Bounds().Dx()
	var height = base_img.Bounds().Dy()
	var scaled_width = width * config.scaling / config.pixel_size
	var scaled_height = height * config.scaling / config.pixel_size
	scaled_img := image.NewRGBA(image.Rect(0, 0, scaled_width, scaled_height))

	for y := 0; y < height; y += config.pixel_size {
		for x := 0; x < width; x += config.pixel_size {
			var totalR, totalG, totalB, totalA, count uint32
			var colors_in_block []color.RGBA
			for dy := 0; dy < config.pixel_size; dy++ {
				for dx := 0; dx < config.pixel_size; dx++ {
					pixel := base_img.At(x+dx, y+dy)
					r, g, b, a := pixel.RGBA()
					colors_in_block = append(colors_in_block, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
					totalR += r
					totalG += g
					totalB += b
					totalA += a
					count++
				}
			}
			avg_color := get_average_color(colors_in_block)

			for sy := 0; sy < config.scaling; sy++ {
				for sx := 0; sx < config.scaling; sx++ {
					scaled_x := (x/config.pixel_size)*config.scaling + sx
					scaled_y := (y/config.pixel_size)*config.scaling + sy
					scaled_img.Set(scaled_x, scaled_y, avg_color)
				}
			}
		}
	}
	return scaled_img, nil
}

func get_closest_palette_color(pixel color.Color) color.RGBA {
	r, g, b, a := pixel.RGBA()

	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func process_png_apply_palette(base_img image.Image, config ProcessConfig) (image.Image, error) {
	var width = base_img.Bounds().Dx()
	var height = base_img.Bounds().Dy()
	fmt.Printf("%v", config)
	colored_img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := base_img.At(x, y)
			closest_palette_color := get_closest_palette_color(pixel)
			colored_img.Set(x, y, closest_palette_color)
		}
	}
	return colored_img, nil
}

func create_png(img image.Image, base_path string, file_name string) {
	var path = fmt.Sprintf("%s/%s.png", base_path, file_name)
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}

	fmt.Printf("PNG image created successfully!\n")
}
