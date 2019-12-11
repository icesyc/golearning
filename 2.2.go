package main

import (
	"./unitconv"
	"fmt"
	"os"
	"strconv"
)

func main() {
	for _, arg := range os.Args[1:] {
		t, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		m := unitconv.Meter(t)
		ft := unitconv.Foot(t)
		fmt.Printf("%s=%s, %s=%s\n", m, unitconv.MToFt(m), ft, unitconv.FtToM(ft))
	}
}