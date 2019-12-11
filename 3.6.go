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
		epsX = (xmax - xmin) / width
		epsY = (ymax - ymin) / height
	)
	offx := []float64{-epsX, epsX}
	offy := []float64{-epsY, epsY}
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for px := 0; px < width; px++ {
		x := float64(px) / width * (xmax - xmin) + xmin;
		for py := 0; py < height; py++ {
			y := float64(py) / width * (ymax - ymin) + ymin;
			subPixels := make([]color.Color, 0)
			for i := 0; i < 2; i++ {
				for j := 0; j < 2; j++ {
					z := complex(x + offx[i], y + offy[j])
					subPixels = append(subPixels, mandelbrot(z))
				}
			}
			img.Set(px, py, avg(subPixels))
		}
	}
	png.Encode(w, img)
}

func avg(colors []color.Color) color.Color {
	var r, g, b, a uint8
	n := uint32(len(colors))
	for _, c := range colors {
		tr, tg, tb, ta := c.RGBA()	
		r += uint8(tr / n)
		g += uint8(tg / n)
		b += uint8(tb / n)
		a += uint8(ta / n)
	}
	return color.RGBA{r, g, b, a}
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