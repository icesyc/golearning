package main

import (
	"net/http"
	"math/cmplx"
	"image"
	"image/color"
	"image/png"
	"log"
)

func main() {
	http.HandleFunc("/mandel", mandel)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func mandel(w http.ResponseWriter, r *http.Request) {
	const (
		xmin, ymin, xmax, ymax = -2, -2, 2, 2
		width, height = 1024, 1024
	)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for px := 0; px < width; px++ {
		x := float64(px) / width * (xmax - xmin) + xmin;
		for py := 0; py < height; py++ {
			y := float64(py) / width * (ymax - ymin) + ymin;
			z := complex(x, y)
			img.Set(px, py, mandelbrot(z))
		}
	}
	png.Encode(w, img)
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15
	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v * v + z
		if cmplx.Abs(v) > 2 {
			r := 255 - contrast * n
			g := r * 3 % 255
			b := r * 6 % 255
			return color.RGBA{r, g, b, 255}
		}
	}
	return color.RGBA{0, 0, 0, 255}
}