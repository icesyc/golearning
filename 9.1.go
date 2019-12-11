package main

import (
	"fmt"
)

func main() {
	Deposit(100)
	fmt.Printf("%v\n", WithDraw(110))
	fmt.Printf("%v\n", Balance())
}

var deposit = make(chan int)
var balances = make(chan int)
var drawResult = make(chan bool)

func Deposit(amount int) {
	deposit <- amount
}
func WithDraw(amount int) bool{
	deposit <- -amount
	return <-drawResult
}
func Balance() int {
	return <-balances
}
func teller() {
	var balance int
	for {
		select {
		case amount := <-deposit:
			if amount < 0 {
				if amount + balance > 0 {
					balance += amount
					drawResult <- true
				} else {
					drawResult <- false
				}
			} else {
				balance += amount
			}
		case balances <- balance:
		}
	}
}

func init() {
	go teller()
}