package main

import (
	"os"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	for _, url := range os.Args[1:] {

		if !strings.HasPrefix(url, "http://") {
			url = "http://" + url
		}
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch error: %s\n", err)
			os.Exit(1)
		}

		status := resp.Status

		fmt.Printf("http status is %v\n", status)
	}
}