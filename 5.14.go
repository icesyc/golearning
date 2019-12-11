package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func main() {
	listDir := func(file string) []string {
		path, err := filepath.Abs(file)
		fmt.Printf("%v\n", path)
		fileInfo, err := os.Stat(file)
		if err != nil {
			fmt.Printf("%s stat error: %v", file, err)
		}
		var files []string
		if fileInfo.IsDir() {
			fileList, _ := ioutil.ReadDir(file)
			for _, f := range fileList {
				files = append(files, file + "/" + f.Name())
			}
			return files
		}
		return nil
	}
	BreadthFirstSearch(listDir, os.Args[1:])
}

func BreadthFirstSearch(f func(string) []string, workList []string) {
	seen := make(map[string]bool)
	for len(workList) > 0 {
		items := workList
		workList = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				workList = append(workList, f(item)...)
			}
		}
	}
}
