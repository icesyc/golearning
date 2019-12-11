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

	counts := make(map[string]int)
	for tag, n:= range countTag(doc, counts) {
		fmt.Printf("%s\t%d\n", tag, n)
	}
}

func countTag(node *html.Node, counts map[string]int) map[string]int{
	if node.Type == html.ElementNode {
		counts[node.Data]++
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		countTag(c, counts)
	}
	return counts
}
