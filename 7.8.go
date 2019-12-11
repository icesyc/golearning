package main

import (
	"fmt"
	"sort"
	"time"
	"os"
	"text/tabwriter"
)

type Track struct {
	Title string
	Artist string
	Album string
	Year int
	Length time.Duration
}

func main() {
	var tracks = []*Track{
		{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
		{"Go", "Moby", "Moby", 1992, length("3m37s")},
    	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
    	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
	}
	sortBy := os.Args[1:]
	trackSort := TrackSort{tracks, sortBy}
	sort.Sort(trackSort)
	PrintTracks(tracks)
}

func length(s string) time.Duration{
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
}

func PrintTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush()
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

