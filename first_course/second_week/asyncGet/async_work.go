package main

import (
	"fmt"
	"time"
)

func getComments() chan string {
	resultCH := make(chan string, 1)
	go func(out chan<- string) {
		time.Sleep(2 * time.Second)
		fmt.Println("async operations ready, return comments")
		out <- "32 комментария"
	}(resultCH)
	return resultCH
}

func getPage() {
	resultCH := getComments()

	time.Sleep(time.Second)
	fmt.Println("get related pages")

	commentsData := <-resultCH
	fmt.Println("main gourrutine:", commentsData)
}

func main() {
	for i := 0; i <= 3; i++ {
		getPage()
	}
}
