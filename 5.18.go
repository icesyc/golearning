package main

import (
	"fmt"
	"os"
	"net/http"
	"path"
	"io"
)

func main() {
	_, _, err := fetch(os.Args[1])
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func fetch(url string) (filename string, n int64, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("fetch %s error: %v", url, resp.Status)
	}	
	filename = path.Base(resp.Request.URL.Path)
	if filename == "/" || filename == "." {
		filename = "index.html"
	}
	f, err := os.Create(filename)
	close := func() {
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}
	defer close()
	if err != nil {
		return "", 0, fmt.Errorf("create file %s error: %v", filename, err)
	}
	n, err = io.Copy(f, resp.Body)
	return filename, n, err
}