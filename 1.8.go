package main

import (
	"os"
	"fmt"
	"net/http"
	"io/ioutil"
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

		data, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close();
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch error: %s\n", err)
			os.Exit(2)
		}

		fmt.Printf("%s\n", data)
	}
}