package main

import "fmt"

func main() {
	in := make(chan int)
	go func(out chan<- int) {
		for i := 0; i <= 4; i++ {
			fmt.Println("before", i)
			out <- i
			fmt.Println("after", i)
		}
		close(out)
		fmt.Println("generatpr finish")
	}(in)

	for i := range in {
		fmt.Println("\tget", i)
	}
}
