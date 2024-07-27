package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

type ProcessConfig struct {
	pixel_size   int
	pixel_width  int
	pixel_height int
	scaling      int
	bg_color     color.RGBA
	palette      []color.RGBA
}

var config = ProcessConfig{
	pixel_size:   16,
	pixel_width:  16,
	pixel_height: 16,
	scaling:      1,
	bg_color:     color.RGBA{0, 0, 0, 0},
	palette: []color.RGBA{
		{255, 0, 0, 255},     // red
		{255, 255, 0, 255},   // yellow
		{0, 255, 0, 255},     // green
		{255, 128, 0, 255},   // orange
		{0, 0, 255, 255},     // blue
		{0, 255, 255, 255},   // cyan
		{255, 0, 127, 255},   // magenta
		{255, 0, 255, 255},   // pink
		{0, 0, 0, 255},       // black
		{255, 255, 255, 255}, // white
	},
}

func main() {
	var input_path = "./input"
	var input_filename = "input"
	var tmp_path = "./tmp"
	var tmp_filename = "working"
	var output_path = "./output"
	var output_filename = "output"

	flag_input_file := flag.String("input", "raw_input", "name of the input file excluding file type extension")
	flag_output_file := flag.String("output", "output", "name of the output file excluding file type extension")

	flag_pixelise := flag.Bool("pixelise", true, "pixelise the input image")
	flag_apply_palette := flag.Bool("apply-palette", false, "apply a palette to the input image")
	flag_invert_colors := flag.Bool("invert", false, "invert the colors in the input image")

	flag_pixel_size := flag.Int("pixel-size", 16, "the size of the resulting pixels as a portion of the input")
	flag_pixel_width := flag.Int("pixel-width", 16, "the width of the resulting pixels as a portion of the input")
	flag_pixel_height := flag.Int("pixel-height", 16, "the height of the resulting pixels as a portion of the input")
	flag_scaling := flag.Int("scale", 1, "invert the colors in the input image")

	config.pixel_size = *flag_pixel_size
	config.pixel_width = *flag_pixel_width
	config.pixel_height = *flag_pixel_height
	config.scaling = *flag_scaling

	input_filename = *flag_input_file
	output_filename = *flag_output_file

	flag.Parse()
	fmt.Printf("image processing started, working image stored in %s/%s\n", tmp_path, tmp_filename)
	var working_image image.Image
	var err error
	working_image, err = get_base_image(input_path, input_filename)
	if err != nil {
		log.Fatalf("Failed to get base image: %v", err)
	}

	crop_required, crop_width, crop_height := get_required_crop_lengths(working_image, config)
	if crop_required {
		working_image, err = crop_image(working_image, crop_width, crop_height)
		if err != nil {
			log.Fatalf("Failed to crop input PNG: %v", err)
		}
	}

	if *flag_pixelise {
		working_image, err = process_png_pixelise(working_image, config)
		if err != nil {
			log.Fatalf("Failed to process PNG: %v", err)
		}
		create_png(working_image, output_path, "pixelated")
	}

	if *flag_apply_palette {
		working_image, err = process_png_apply_palette(working_image, config)
		if err != nil {
			log.Fatalf("Failed to process PNG: %v", err)
		}
		create_png(working_image, output_path, "palette_applied")
	}

	if *flag_invert_colors {
		working_image, err = process_png_invert_colors(working_image, config)
		if err != nil {
			log.Fatalf("Failed to process PNG: %v", err)
		}
		create_png(working_image, output_path, "inversion_applied")
	}

	create_png(working_image, output_path, output_filename)
}

func get_base_image(base_path string, file_name string) (image.Image, error) {
	var path = fmt.Sprintf("%s/%s.png", base_path, file_name)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()
	base_img, err := png.Decode(file)
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

func calculate_color_diff(color1, color2 color.Color) int {
	r1, g1, b1, _ := color1.RGBA()
	fmt.Printf("r1:%v, g1:%v, b1:%v\n", r1, g1, b1)
	r2, g2, b2, _ := color2.RGBA()
	fmt.Printf("r2:%v, g2:%v, b2:%v\n", r2, g2, b2)

	r1s, r2s := int(r1>>8), int(r2>>8)
	g1s, g2s := int(g1>>8), int(g2>>8)
	b1s, b2s := int(b1>>8), int(b2>>8)

	r_diff := r1s - r2s
	g_diff := g1s - g2s
	b_diff := b1s - b2s

	euclid_dist := math.Sqrt(float64(r_diff*r_diff + g_diff*g_diff + b_diff*b_diff))
	fmt.Printf("euclid_dist: %v\n", euclid_dist)

	return int(euclid_dist)
}

func get_closest_color_in_palette(pixel color.Color, palette []color.RGBA) color.RGBA {
	var closest_palette_color = color.RGBA{0, 0, 0, 0}
	var smallest_color_diff int = 99999999999999
	for _, color := range palette {
		color_diff := calculate_color_diff(color, pixel)
		if color_diff < smallest_color_diff {
			smallest_color_diff = color_diff
			closest_palette_color = color
		}
	}
	return closest_palette_color
}

// func set_image_block(img image.Image, x1 int, y1 int, x2 int, y2 int) image.Image {
// helper function to set all pixels a single color from topleft x1,y1 to bottom right x2,y2
// }

func process_png_apply_palette(base_img image.Image, config ProcessConfig) (image.Image, error) {
	var width = base_img.Bounds().Dx()
	var height = base_img.Bounds().Dy()
	fmt.Printf("%v", config)
	colored_img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := base_img.At(x, y)
			closest_color_in_palette := get_closest_color_in_palette(pixel, config.palette)
			fmt.Printf("closest color:%v\n", closest_color_in_palette)
			colored_img.Set(x, y, closest_color_in_palette)
		}
	}
	return colored_img, nil
}

func get_required_crop_lengths(base_img image.Image, config ProcessConfig) (bool, int, int) {
	var crop_required bool
	var required_crop_width, required_crop_height int
	var width = base_img.Bounds().Dx()
	var height = base_img.Bounds().Dy()

	var width_remainder int = width % config.pixel_width
	var height_remainder int = height % config.pixel_height

	if width_remainder == 0 {
		required_crop_width = width
	} else {
		required_crop_width = width - width_remainder
		crop_required = true
	}
	if height_remainder == 0 {
		required_crop_height = height
	} else {
		required_crop_height = height - height_remainder
		crop_required = true
	}
	return crop_required, required_crop_width, required_crop_height
}

func crop_image(base_img image.Image, cropped_width int, cropped_height int) (image.Image, error) {
	var width = base_img.Bounds().Dx()
	var height = base_img.Bounds().Dy()
	cropped_image := image.NewRGBA(image.Rect(0, 0, width, height))
	if cropped_width >= width && cropped_height >= height {
		fmt.Printf("image is smaller than crop size, returning early with base image")
		return base_img, nil
	}
	if cropped_width > width {
		fmt.Printf("cropped_width is greater than input width, not cropping width")
		cropped_width = width
	}
	if cropped_height > height {
		fmt.Printf("cropped_height is greater than input height, not cropping height")
		cropped_height = height
	}

	for y := 0; y < cropped_height; y++ {
		for x := 0; x < cropped_width; x++ {
			pixel := base_img.At(x, y)
			cropped_image.Set(x, y, pixel)
		}
	}

	return cropped_image, nil
}

func get_inverted_pixel(pixel color.Color) color.Color {
	r, g, b, a := pixel.RGBA()
	rs, gs, bs, as := r>>8, g>>8, b>>8, a>>8
	fmt.Printf("rs:%v,gs:%v,bs:%v,as:%v\n", rs, gs, bs, as)
	var inverted_pixel = color.RGBA{uint8(255 - rs), uint8(255 - gs), uint8(255 - bs), uint8(as)}
	return inverted_pixel
}

func process_png_invert_colors(base_img image.Image, config ProcessConfig) (image.Image, error) {
	var width = base_img.Bounds().Dx()
	var height = base_img.Bounds().Dy()
	fmt.Printf("%v", config)
	colored_img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := base_img.At(x, y)
			var inverted_pixel = get_inverted_pixel(pixel)
			colored_img.Set(x, y, inverted_pixel)
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
