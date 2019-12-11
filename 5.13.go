package main

import (
	"os"
	"fmt"
	"net/http"
	"net/url"
	"golang.org/x/net/html"
	"log"
	"io/ioutil"
	"strings"
)

const maxPage = 200

func main() {
	var u = os.Args[1]
	uu, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}
	host := uu.Hostname()
	os.Mkdir(host, 0755)
	err = os.Chdir(host)
	if err != nil {
		log.Fatal(err)
	}

	pages := make(map[string]bool)
	crawl := func(u string) []string {
		uu, err := url.Parse(u)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s error: %v", u, err)
		}
		//只抓原始host的内容
		if uu.Hostname() != host {
			return nil
		}
		fmt.Println(u)
		content, links, err := Extract(u)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
		}
		var values []string
		for _, v := range uu.Query() {
			values = append(values, v...)
		}
		fname := strings.Join(values, "_") + ".html"
		fname = strings.Replace(fname, "/", "_", -1)
		if fname == ".html" {
			fname = "index.html"
		}
		if !pages[fname] {
			pages[fname] = true
			saveFile(fname, content)
		}
		return links
	}
	BreadthFirstSearch(crawl, []string{u})
}

func saveFile(name string, data string) error {
	f, err := os.Create(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "-> can not create file: %s, %v ", name, err)
		return err
	}
	_, err = f.WriteString(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "-> can not write file: %s, %v ", name, err)
		return err
	}
	err = f.Close()
	return err
}

func ForeachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ForeachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

func Extract(url string) (string, []string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		err := fmt.Errorf("getting %s error: %s", url, resp.Status)
		return "", nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	doc, err := html.Parse(strings.NewReader(string(content)))
	resp.Body.Close()
	if err != nil {
		return "", nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link, err := resp.Request.URL.Parse(attr.Val)
					if err != nil {
						continue
					}
					links = append(links, link.String())
				}
			}
		}
	}
	ForeachNode(doc, visitNode, nil)
	return string(content), links, nil
}

func BreadthFirstSearch(f func(string) []string, workList []string) {
	seen := make(map[string]bool)
	var n int
	for len(workList) > 0 {
		items := workList
		workList = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				workList = append(workList, f(item)...)
				n++
			}
			if n > maxPage {
				return
			}
		}
	}
}
