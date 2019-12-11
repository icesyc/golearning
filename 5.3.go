package main

import (
	"fmt"
	"os"
	"golang.org/x/net/html"
	"strings"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Printf("html parse error: %v\n", err)
		os.Exit(1)
	}

	visit(doc)
}

func visit(node *html.Node){
	if node.Type == html.ElementNode && (node.Data == "script" || node.Data == "style") {
		return
	}
	if node.Type == html.TextNode {
		text := strings.Trim(node.Data, " \t\n\r")
		if text != "" {
			fmt.Printf("%s\n", text)
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		visit(c)
	}
}
