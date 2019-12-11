package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func(){
		var ch chan int
		//ch = make(chan int, 2)
		ch <- 1
		fmt.Println("goroutine done.")
		wg.Done()
	}()
	wg.Wait()
	fmt.Println("main done")

}