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

	zMin, zMax := zrange()
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)
			if math.IsNaN(ax) || math.IsNaN(ay) || math.IsNaN(bx) || math.IsNaN(by) || math.IsNaN(cx) || math.IsNaN(cy) || math.IsNaN(dx) || math.IsNaN(dy) {
				continue
			}
			color := zToColor(i, j, zMin, zMax)
			fmt.Fprintf(w, "<polygon points='%g,%g,%g,%g,%g,%g,%g,%g' style='stroke: %s; fill:%s'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy, color, color)
		}
	}
	fmt.Fprintf(w, "</svg>")
}

func zrange() (float64, float64) {
	var min, max = math.NaN(), math.NaN()
	for i := 0;  i < cells; i++ {
		for j := 0; j < cells; j++ {
			for xoff := 0; xoff < 2; xoff++ {
				for yoff := 0; yoff < 2; yoff++ {
					z := zvalue(i + xoff, j + yoff)
					if math.IsNaN(min) || z < min {
						min = z
					}
					if math.IsNaN(max) || z > max {
						max = z
					}
				}
			}
		}
	}
	return min, max
}

func zToColor(i int, j int, zMin, zMax float64) string {
	x := xyrange * (float64(i) / cells - 0.5)
	y := xyrange * (float64(j) / cells - 0.5)
	//Compute surface height z
	z := f(x, y)
	percent := (z - zMin) / (zMax - zMin)
	return fmt.Sprintf("#%02x00%02x", int(255 * percent), int(255 * (1 - percent)))
}

func zvalue(i int, j int) float64 {
	x := xyrange * (float64(i) / cells - 0.5)
	y := xyrange * (float64(j) / cells - 0.5)
	//Compute surface height z
	z := f(x, y)
	return z
}

func corner(i, j int) (float64, float64) {
	// Find point (x,y) at corner of cell (i,j)
	x := xyrange * (float64(i) / cells - 0.5)
	y := xyrange * (float64(j) / cells - 0.5)

	//Compute surface height z
	z := f(x, y)

	//Project(x,y,z) isometrically onto 2-D SVG canvas (sx, sy).
	sx := width / 2 + (x - y) * cos30 * xyscale
	sy := height / 2 + (x + y) * sin30 * xyscale - z * zscale
	return sx, sy
}

func f(x, y float64) float64 { 
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}