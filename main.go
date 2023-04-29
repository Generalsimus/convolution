package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

// func main() {

// }

//	type Shape2D struct {
//		inputs []float64
//		widget int
//		height int
//	}
type Input2D struct {
	inputs []float64
	widget int
	height int
}

type Kernel2D struct {
	inputs []float64
	widget int
	height int
}

func (kernel Kernel2D) ConvolutionImage(img image.Image) image.Image {
	kernelLen := len(kernel.inputs)
	// fmt.Println("kernelLen", kernelLen)
	bounds := img.Bounds()
	// maxColor := float64(257 * 255)
	centerX := int(math.Floor(float64(kernel.widget) / 2))
	centerY := int(math.Floor(float64(kernel.height) / 2))

	newImage := image.NewRGBA(image.Rect(0, 0, bounds.Max.X, bounds.Max.Y))
	for row := 0; row < bounds.Max.Y; row++ {
		for column := 0; column < bounds.Max.X; column++ {
			var red float64 = 0
			var green float64 = 0
			var blue float64 = 0
			var alpha float64 = 0
			for k := 0; k < kernelLen; k++ {
				kx := k % kernel.widget
				ky := (k - kx) / kernel.widget

				x := (column + kx - centerX)
				y := (row + ky - centerY)
				r, g, b, a := img.At(x, y).RGBA()
				kernelInput := kernel.inputs[k]

				red = red + (float64(r>>8) * kernelInput)
				green = green + (float64(g>>8) * kernelInput)
				blue = blue + (float64(b>>8) * kernelInput)
				if column == x && row == y {
					alpha = float64(a >> 8)
				}

			}

			newImage.Set(column, row, color.RGBA{
				uint8(math.Max(math.Min(red, 255), 0)),
				uint8(math.Max(math.Min(green, 255), 0)),
				uint8(math.Max(math.Min(blue, 255), 0)),
				uint8(alpha),
			})

		}
	}

	return newImage
}

type Layer struct {
	inputs []float64
	widget int
	height int
}

type Network struct {
	Layers []Layer
}

func (net Network) addLayer(inputs []float64) {
	net.Layers = append(net.Layers, Layer{
		inputs: inputs,
	})
}

func main() {
	dirPath := "./images/"
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range entries {
		fileName := file.Name()

		kernel := Kernel2D{
			inputs: []float64{
				-1, -1, -1,
				-1, 8, -1,
				-1, -1, -1,
				// -1, -2, -1,
				// 0, 0, 0,
				// 1, 2, 1,
			},
			widget: 3,
			height: 3,
		}
		fmt.Println("CONVOLVED: ", fileName)

		img, _ := getImageFromFilePath(dirPath + fileName)
		saveImageAt(kernel.ConvolutionImage(img), "./save/"+fileName)

	}
	fmt.Println("END")
}
func saveImageAt(image image.Image, path string) {
	var imageBuf bytes.Buffer
	png.Encode(&imageBuf, image)

	// Write to file.
	outfile, err := os.Create(path)
	if err != nil {
		// replace this with real error handling
		panic(err.Error())
	}
	defer outfile.Close()
	png.Encode(outfile, image)
}
func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _ := png.Decode(f)
	return image, err
}

func rgbaToNum[T int | uint | uint32](r T, g T, b T, a T) T {
	rgb := r
	rgb = (rgb << 8) + g
	rgb = (rgb << 8) + b
	rgb = (rgb << 8) + a
	return rgb
}

func numToRgba(num int) (int, int, int, int) {

	red := (num >> 24) & 0xFF
	green := (num >> 16) & 0xFF
	blue := (num >> 8) & 0xFF
	a := num & 0xFF

	return red, green, blue, a
}
