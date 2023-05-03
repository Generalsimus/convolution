package main

import (
	"convolution/convolution"
	"fmt"

	"gocv.io/x/gocv"
)

// C:\Program Files (x86)\CMake\bin;C:\mingw-w64\x86_64-8.1.0-posix-seh-rt_v6-rev0\mingw64\bin
func main() {
	webcam, _ := gocv.OpenVideoCapture(1)
	window := gocv.NewWindow("Hello")
	img := gocv.NewMat()

	kernel := convolution.Kernel2D{
		Inputs: []float64{
			-1, -1, -1,
			-1, 8, -1,
			-1, -1, -1,
			// -1, -2, -1,
			// 0, 0, 0,
			// 1, 2, 1,
			// 0.111, 0.111, 0.111,
			// 0.111, 0.111, 0.111,
			// 0.111, 0.111, 0.111,
		},
		Widget: 3,
		Height: 3,
	}

	// window.setMouseCallback
	// w := gocv.NewWindow("test")
	// mouse_callback := func(event int, x int, y int, flags int) {
	// 	if event == cv.EVENT_LBUTTONDOWN {
	// 		Prln("mouse click at: " + I2S(x) + "," + I2S(y))
	// 	}
	// }
	// w.SetMouseCallback(mouse_callback)
	fmt.Println("CONVOLVED: ", kernel)

	// ioutil.WriteFile("output.png", , 0644)
	// convolvedImg := kernel.ConvolutionImage(image)

	// conv, _ := gocv.ImageToMatRGBA(convolvedImg)
	// ioutil.WriteFile("output.png", , 0644)

	for {
		webcam.Read(&img)
		image, _ := img.ToImage()
		convolvedImg := kernel.ConvolutionImage(image)
		convMapPool := convolution.MaxPoolImage(convolvedImg, 5, 5)
		updatedIMG, _ := gocv.ImageToMatRGBA(convMapPool)
		// windoww.get
		// window.WaitKey(2)
		// fmt.Println("CLICK")
		window.IMShow(updatedIMG)
		window.WaitKey(1)
	}
}
