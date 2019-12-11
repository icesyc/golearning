package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
	"runtime"
)

type empty struct{}
var sem = make(chan empty, 2)

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	fileSize := make(chan int64)
	var syn sync.WaitGroup
	ticker := time.Tick(5 * time.Millisecond)
	var nfiles, nbytes int64

	syn.Add(1)
	go walkDir(root, &syn, fileSize)

	go func() {
		syn.Wait()
		close(fileSize)
	}()

	loop: 
	for {
		select {
		case size, ok := <-fileSize:
			if !ok {
				break loop
			}
			nfiles++
			nbytes += size
		case <-ticker:
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes)
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files, %.1f GB, %d goroutines\n", nfiles, float64(nbytes)/1e9, runtime.NumGoroutine())
}

func walkDir(dir string, syn *sync.WaitGroup, fileSize chan int64) {
	defer func(){
		syn.Done()
	}()
	for _, ent := range dirents(dir) {
		if ent.IsDir() {
			subDir := filepath.Join(dir, ent.Name())
			syn.Add(1)
			go walkDir(subDir, syn, fileSize)
		} else {
			fileSize <- ent.Size()
		}
	}

}

func dirents(dir string) []os.FileInfo {
	sem <- empty{}
	defer func() { <-sem }()
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du error:%v", err)
		return nil
	}
	return entries
}
