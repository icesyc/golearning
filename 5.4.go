package main

import (
	"fmt"
	"os"
	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Printf("html parse error: %v\n", err)
		os.Exit(1)
	}
	tags := visit(nil, doc)
	for _, tag := range tags {
		fmt.Println(tag)
	}
}

func visit(tags []string, node *html.Node) []string{
	if node.Type == html.ElementNode {
		tags = append(tags, "<" + node.Data + ">")
	}
	if node.FirstChild != nil {
		tags = visit(tags, node.FirstChild)
	}
	if node.NextSibling != nil {
		tags = visit(tags, node.NextSibling)
	}
	return tags
}