package main

import (
	"fmt"
	"net/http"
	"math"
	"log"
)

const (
	width, height = 600, 320			//canvas size in pixels
	cells 		  = 100		 			//number of grid cells
	xyrange 	  = 30.0				//axis range (-xyrange..+xyrange)
	xyscale		  = width / 2 / xyrange	//pixels per x or y unit
	zscale		  = height * 0.4		//pixels per z unit
	angle		  = math.Pi / 6			//angle of x, y axes (=30°)

)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

func main() {
	http.HandleFunc("/svg", svg)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func svg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: gray; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)
			if math.IsNaN(ax) || math.IsNaN(ay) || math.IsNaN(bx) || math.IsNaN(by) || math.IsNaN(cx) || math.IsNaN(cy) || math.IsNaN(dx) || math.IsNaN(dy) {
				continue
			}
			fmt.Fprintf(w, "<polygon points='%g,%g,%g,%g,%g,%g,%g,%g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Fprintf(w, "</svg>")
}


func corner(i, j int) (sx, sy float64) {
	// Find point (x,y) at corner of cell (i,j)
	x := xyrange * (float64(i) / cells - 0.5)
	y := xyrange * (float64(j) / cells - 0.5)

	//Compute surface height z
	z := bessel(x, y)

	//Project(x,y,z) isometrically onto 2-D SVG canvas (sx, sy).
	sx = width / 2 + (x - y) * cos30 * xyscale
	sy = height / 2 + (x + y) * sin30 * xyscale - z * zscale
	return
}

func eggbox(x, y float64) float64 {
	return 0.2 * (math.Cos(x) + math.Cos(y))
}

func saddle(x, y float64) float64 {
	a := 25.0
	b := 17.0
	a2 := a*a
	b2 := b*b
	return (y*y/a2 - x*x/b2)
}

func bessel(x, y float64) float64 {
	return 0.1 * (math.Sin(x) + math.Sin(y))
}

func f(x, y float64) float64 { 
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}