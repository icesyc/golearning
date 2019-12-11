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
	el := ElementByID(doc, "pjax-container")
	fmt.Printf("%v\n", el)
}

func foreachNode(n *html.Node, pre, post func(n *html.Node) bool) bool{
	if pre != nil && !pre(n) {
		return false
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if !foreachNode(c, pre, post) {
			return false
		}
	}
	if post != nil && !post(n){
		return false
	}
	return true
}

func ElementByID(doc *html.Node, id string) *html.Node{
	var node *html.Node
	queryFunc := func(n *html.Node) bool {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if attr.Key == "id" && attr.Val == id {
					node = n
					return false
				}
			}
		}
		return true
	}
	foreachNode(doc, queryFunc, nil)
	return node
}
