package main

import (
	"fmt"
	"os"
	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "html parse error: %s", err)
	}
	els := ElementsByTagName(doc, "meta")
	for _, el := range els {
		fmt.Printf("%v\n", el)
	}
}

func ForEachNode(n *html.Node, visitor func(n *html.Node) bool) bool {
	if visitor != nil && !visitor(n) {
		return false
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if !ForEachNode(c, visitor) {
			return false
		}
	}
	return true
}

func InArray(arr []string, item string) bool{
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

func ElementsByTagName(doc *html.Node, names...string) []*html.Node {
	var nodeList []*html.Node
	queryFunc := func(n *html.Node) bool {
		if n.Type == html.ElementNode && InArray(names, n.Data) {
			nodeList = append(nodeList, n)
		}
		return true
	}
	ForEachNode(doc, queryFunc)
	return nodeList
}
