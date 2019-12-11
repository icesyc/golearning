package main

import (
	"fmt"
	"os"
	"golang.org/x/net/html"
)

type bailout struct{}

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Printf("parse error: %v\n", err)
	}
	title, err := soleTitle(doc)
	if err != nil {
		fmt.Printf("soleTitle error: %v\n", err)
	}
	fmt.Printf("title: %s\n", title)
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

func soleTitle(doc *html.Node) (title string, err error) {
	defer func() {
		bailout := bailout{}
		if p := recover(); p != nil {
			if p == bailout {
				err = fmt.Errorf("multiple title elements")
			} else {
				panic(p)
			}

		}
	}()
	findTitle := func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			if title != "" {
				panic(bailout{})
			}
			title = n.FirstChild.Data
		}
		return true
	}
	ForEachNode(doc, findTitle)
	if title == "" {
		err = fmt.Errorf("no title element")
	}
	return title, err
}