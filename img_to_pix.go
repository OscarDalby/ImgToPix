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

type PixelData struct {
	x int
	y int
	c color.RGBA
}

var config = ProcessConfig{pixel_size: 32, scaling: 1, bg_color: color.RGBA{0, 0, 0, 0}}

func main() {
	base_img, err := get_base_image("./input", "selfie", config)
	if err != nil {
		log.Fatalf("Failed to get base image: %v", err)
	}
	processed_img, scaled_img, err := process_png(base_img, config)
	if err != nil {
		log.Fatalf("Failed to process PNG: %v", err)
	}
	create_png(processed_img, "./output", "output_image")
	create_png(scaled_img, "./output", "scaled_output_image")
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

func average_color(colors []color.RGBA) color.RGBA {
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

func process_png(base_img image.Image, config ProcessConfig) (image.Image, image.Image, error) {
	fmt.Printf("%v", config)
	var width = base_img.Bounds().Dx()
	var height = base_img.Bounds().Dy()
	var scaled_width = width / config.pixel_size
	var scaled_height = height / config.pixel_size

	if width%config.pixel_size != 0 {
		panic("pixel_size must divide width")
	}
	if height%config.pixel_size != 0 {
		panic("pixel_size must divide height")
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	scaled_img := image.NewRGBA(image.Rect(0, 0, scaled_width, scaled_height))

	fmt.Printf("new image instantiated\n")

	if config.bg_color.A != 0 { // if the bg_color is not completely transparent, then fill the background before the remaining processing
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				img.Set(x, y, config.bg_color)
			}
		}
	}

	fmt.Printf("background color set\n")
	var total_pixels_to_process = height * width
	total_count := 0
	fmt.Printf("x * y %v\n", total_pixels_to_process)

	for y := 0; y < height; y += config.pixel_size {
		for x := 0; x < width; x += config.pixel_size {
			var color_list []color.RGBA
			var pixel_list []PixelData
			var avg_color color.RGBA

			count := 0

			for dy := 0; dy < config.pixel_size; dy++ {
				for dx := 0; dx < config.pixel_size; dx++ {
					pixel := base_img.At(x+dx, y+dy)
					r, g, b, a := pixel.RGBA()
					color_list = append(color_list, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
					avg_color = average_color(color_list)
					pixel_list = append(pixel_list, PixelData{x: x, y: y, c: avg_color})
					// fmt.Printf("avg_color: %v\n", avg_color)

					for j := 0; j < config.pixel_size; j++ {
						for i := 0; i < config.pixel_size; i++ {
							img.Set(x+i, y+j, avg_color)
						}
					}

					for j := 0; j < scaled_height/config.pixel_size; j++ {
						for i := 0; i < scaled_width/config.pixel_size; i++ {
							scaled_img.Set(x+i, y+j, avg_color)
						}
					}

					count++
				}
			}
			if (total_pixels_to_process-total_count)%1000 == 0 {
				fmt.Printf("%v%% of pixels processed\n", total_count*100/total_pixels_to_process)
			}
			total_count++

			fmt.Printf("pixel_list: %v\n", pixel_list)
			fmt.Printf("number of pixels in pixel_list: %v\n", len(pixel_list))

		}
	}
	fmt.Printf("returning proccessed img\n")

	return img, scaled_img, nil
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

	fmt.Println("PNG image created successfully!")
}
