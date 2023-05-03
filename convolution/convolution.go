package convolution

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
	Inputs []float64
	Widget int
	Height int
}

func (kernel Kernel2D) ConvolutionImage(img image.Image) image.Image {
	kernelLen := len(kernel.Inputs)
	bounds := img.Bounds()
	// maxColor := float64(257 * 255)
	centerX := int(math.Floor(float64(kernel.Widget) / 2))
	centerY := int(math.Floor(float64(kernel.Height) / 2))

	newImage := image.NewRGBA(image.Rect(0, 0, bounds.Max.X, bounds.Max.Y))
	for row := 0; row < bounds.Max.Y; row++ {
		for column := 0; column < bounds.Max.X; column++ {
			var red float64 = 0
			var green float64 = 0
			var blue float64 = 0
			var alpha float64 = 0
			for k := 0; k < kernelLen; k++ {
				kx := k % kernel.Widget
				ky := (k - kx) / kernel.Widget

				x := (column + kx - centerX)
				y := (row + ky - centerY)
				r, g, b, a := img.At(x, y).RGBA()
				kernelInput := kernel.Inputs[k]

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
func (kernel Kernel2D) ConvolutionPoolImage(img image.Image) image.Image {
	kernelLen := len(kernel.Inputs)
	bounds := img.Bounds()
	//////////////////////
	newWidget := bounds.Max.X - kernel.Widget
	newHeight := bounds.Max.Y - kernel.Height

	newImage := image.NewRGBA(image.Rect(0, 0, newWidget, newHeight))
	for row := 0; row < newHeight; row++ {
		for column := 0; column < newWidget; column++ {
			var red float64 = 0
			var green float64 = 0
			var blue float64 = 0
			var alpha float64 = 0
			// if ro
			for k := 0; k < kernelLen; k++ {
				kx := k % kernel.Widget
				ky := (k - kx) / kernel.Widget

				x := column + kx
				y := row + ky
				r, g, b, a := img.At(x, y).RGBA()
				kernelInput := kernel.Inputs[k]

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
func MaxPoolImage(img image.Image, widget int, height int) image.Image {
	bounds := img.Bounds()
	//////////////////////
	newWidget := int(math.Round(float64(bounds.Max.X / widget)))
	newHeight := int(math.Round(float64(bounds.Max.Y / height)))

	newImage := image.NewRGBA(image.Rect(0, 0, newWidget, newHeight))

	downPixelSize := float64(height * widget)
	for row := 0; row < newHeight; row++ {
		for column := 0; column < newWidget; column++ {
			var maxRed float64 = 0
			var maxGreen float64 = 0
			var maxBlue float64 = 0
			var averageAlpha float64 = 0

			for rowY := (row * height); rowY < ((row * height) + height); rowY++ {

				for columnX := (column * widget); columnX < ((column * widget) + widget); columnX++ {
					r, g, b, a := img.At(columnX, rowY).RGBA()

					maxRed = math.Max(maxRed, float64(r))
					maxGreen = math.Max(maxGreen, float64(g))
					maxBlue = math.Max(maxBlue, float64(b))
					averageAlpha += float64(a)
				}
			}

			newImage.Set(column, row, color.RGBA{
				uint8(maxRed),
				uint8(maxGreen),
				uint8(maxBlue),
				uint8(averageAlpha / downPixelSize),
			})

		}
	}

	return newImage
}
func AveragePoolImage(img image.Image, widget int, height int) image.Image {
	bounds := img.Bounds()
	//////////////////////
	newWidget := int(math.Round(float64(bounds.Max.X / widget)))
	newHeight := int(math.Round(float64(bounds.Max.Y / height)))

	newImage := image.NewRGBA(image.Rect(0, 0, newWidget, newHeight))

	downPixelSize := uint32(height * widget)
	for row := 0; row < newHeight; row++ {
		for column := 0; column < newWidget; column++ {
			var maxRed uint32 = 0
			var maxGreen uint32 = 0
			var maxBlue uint32 = 0
			var averageAlpha uint32 = 0

			for rowY := (row * height); rowY < ((row * height) + height); rowY++ {

				for columnX := (column * widget); columnX < ((column * widget) + widget); columnX++ {
					r, g, b, a := img.At(columnX, rowY).RGBA()

					maxRed = (maxRed + uint32(r>>8))
					maxGreen = (maxGreen + uint32(g>>8))
					maxBlue = (maxBlue + uint32(b>>8))
					averageAlpha = (averageAlpha + uint32(a>>8))

				}
			}

			newImage.Set(column, row, color.RGBA{
				uint8(maxRed / downPixelSize),
				uint8(maxGreen / downPixelSize),
				uint8(maxBlue / downPixelSize),
				uint8(averageAlpha / downPixelSize),
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

func maiwn() {
	dirPath := "./images/"
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range entries {
		fileName := file.Name()

		kernel := Kernel2D{
			Inputs: []float64{
				// -1, -1, -1,
				// -1, 8, -1,
				// -1, -1, -1,
				// -1, -2, -1,
				// 0, 0, 0,
				// 1, 2, 1,
				0.111, 0.111, 0.111,
				0.111, 0.111, 0.111,
				0.111, 0.111, 0.111,
			},
			Widget: 3,
			Height: 3,
		}
		fmt.Println("CONVOLVED: ", kernel)

		// img, _ := getImageFromFilePath(dirPath + fileName)
		// saveImageAt(kernel.ConvolutionImage(img), "./save/"+fileName)
		// fmt.Println("POOLED: ", fileName)

		// imgPool, _ := getImageFromFilePath(dirPath + fileName)
		// saveImageAt(kernel.ConvolutionPoolImage(imgPool), "./save/pool"+fileName)

		imgMaxPool, _ := getImageFromFilePath(dirPath + fileName)
		saveImageAt(MaxPoolImage(imgMaxPool, 30, 30), "./save/MaxPool"+fileName)

		imgAveragePool, _ := getImageFromFilePath(dirPath + fileName)
		saveImageAt(AveragePoolImage(imgAveragePool, 30, 30), "./save/AveragePool"+fileName)

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
