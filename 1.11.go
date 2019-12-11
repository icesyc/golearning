package main

import (
	"log"
	"net/http"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"math"
	"math/rand"
	"strconv"
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/lissajous", lissajous)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}


var palette = []color.Color{
	color.White,
	color.RGBA{0xff, 0x00, 0x00, 1},
	color.RGBA{0x00, 0xff, 0x00, 1},
	color.RGBA{0x00, 0x00, 0xff, 1},
	color.RGBA{0xff, 0xff, 0x00, 1},
	color.RGBA{0xff, 0x00, 0xff, 1},
	color.RGBA{0x00, 0xff, 0xff, 1},
}

func lissajous(w http.ResponseWriter, r *http.Request){
	const (
		res = 0.001
		size = 100
		nframes = 64
		delay = 8
	)

	cycles, _ := strconv.Atoi(r.FormValue("cycles"))
	if cycles == 0 {
		cycles = 5
	}
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2 * size + 1, 2 * size + 1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < float64(cycles) * 2 * math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t * freq + phase)
			img.SetColorIndex(size + int(x * size + 0.5), size + int(y * size + 0.5), uint8(rand.Intn(7) + 1))
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(w, &anim)
}