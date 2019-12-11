package main

import (
	"os"
	"fmt"
	"net/http"
	"bufio"
	"strings"
	"golang.org/x/net/html"
)

func main() {
	url := "http://www.tenholes.com"
	words, images, err := CountWordsAndImages(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("words: %d, images: %d\n", words, images)
}

func CountWordsAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parse html error: %s", err)
		return
	}

	words, images = _CountWordsAndImages(doc)
	return
}

func _CountWordsAndImages(node *html.Node) (words, images int) {
	if node.Type == html.ElementNode && (node.Data == "script" || node.Data == "style") {
		return
	}
	if node.Type == html.ElementNode && node.Data == "img" {
		images++
	}
	if node.Type == html.TextNode {
		scanner := bufio.NewScanner(strings.NewReader(node.Data))
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			words++
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "scanner error:", err)
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		w, i := _CountWordsAndImages(c)
		words += w
		images += i
	}
	return
}