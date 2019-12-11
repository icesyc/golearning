package main

import (
	"sort"
	"time"
	"fmt"
	"net/http"
	"html/template"
)

type Track struct {
	Title string
	Artist string
	Album string
	Year int
	Length time.Duration
}

func main() {
	http.HandleFunc("/", sortTable)
	http.ListenAndServe("localhost:8080", nil)
}

var sortBy []string
func sortTable(w http.ResponseWriter, r *http.Request) {
	var tracks = []*Track{
		{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
		{"Go", "Moby", "Moby", 1992, length("3m37s")},
    	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
    	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
	}
	newField := r.FormValue("s")
	for i, f := range sortBy {
		if f == newField {
			sortBy = append(sortBy[:i], sortBy[i+1:]...)
		}
	}
	sortBy = append([]string{newField}, sortBy...)
	fmt.Fprintf(w, "<p>%s</p>", sortBy)
	trackSort := TrackSort{tracks, sortBy}
	sort.Sort(trackSort)
	PrintTracks(w, tracks)
}

func length(s string) time.Duration{
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
}

func PrintTracks(w http.ResponseWriter, tracks []*Track) {
	tpl := `
		<table>
			<tr>
				<th><a href="/?s=Title">Title</a></th>
				<th><a href="/?s=Artist">Artist</a></th>
				<th><a href="/?s=Album">Album</a></th>
				<th><a href="/?s=Year">Year</a></th>
				<th><a href="/?s=Length">Length</a></th>
			</tr>
			{{range .}}
			<tr>
				<td>{{.Title}}</td>
				<td>{{.Artist}}</td>
				<td>{{.Album}}</td>
				<td>{{.Year}}</td>
				<td>{{.Length}}</td>
			</tr>
			{{end}}
		</table>`
	html := template.Must(template.New("html").Parse(tpl))
	if err := html.Execute(w, tracks); err != nil {
		panic(err)
	}

}

type TrackSort struct {
	tracks []*Track
	sortBy []string

}

func (x TrackSort) Less(i, j int) bool {
	for _, field := range x.sortBy {
		switch field {
		case "Title":
			if x.tracks[i].Title != x.tracks[j].Title {
				return x.tracks[i].Title < x.tracks[j].Title
			}
		case "Artist":
			if x.tracks[i].Artist != x.tracks[j].Artist {
				return x.tracks[i].Artist < x.tracks[j].Artist
			}
		case "Album":
			if x.tracks[i].Album != x.tracks[j].Album {
				return x.tracks[i].Album < x.tracks[j].Album
			}
		case "Year":
			if x.tracks[i].Year != x.tracks[j].Year {
				return x.tracks[i].Year < x.tracks[j].Year
			}
		case "Length":
			if x.tracks[i].Length != x.tracks[j].Length {
				return x.tracks[i].Length < x.tracks[j].Length
			}
		}
	}
	return false
} 
func (x TrackSort) Len() int { return len(x.tracks) }
func (x TrackSort) Swap(i, j int) { x.tracks[i], x.tracks[j] = x.tracks[j], x.tracks[i] }

