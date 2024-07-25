# ImgToPix

`ImgToPix` is a Go application that converts standard images into pixel art images.

## Supported Image Input

`ImgToPix` supports the following image formats for input:
- PNG

Ensure that the input file is accessible and has the proper permissions for reading.

## Supported Image Output

The application can generate pixel art images in the following formats:

- PNG (Recommended for preserving transparency and quality)


## Configuration

`ImgToPix` offers several configuration options to customize the output image:

- **Resolution**: Set the resolution of the output pixel art. Lower resolutions result in a more pronounced pixel effect.
<!-- - **Palette**: Choose from predefined color palettes or create your own to define the color scheme of the output image. -->
<!-- - **Dithering**: Enable or disable dithering to achieve different stylistic effects in the color transition of the pixel art. -->

To use these configurations, pass the respective flags or settings through the command line or a configuration file. For example:

```bash
./ImgToPix --input=path/to/input.jpg --output=path/to/output.png --resolution=32x32
