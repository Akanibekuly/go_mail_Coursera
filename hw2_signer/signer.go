package main

import "fmt"

// ExecutePipeline is
func ExecutePipeline(jobs ...job) {
	// fmt.Println("JOBS: ", jobs)
	in := make(chan interface{}, 1)
	out := make(chan interface{}, 1)
	for _, j := range jobs {
		j(in, out)
	}
}

// SingleHash is
func SingleHash(in, out chan interface{}) {
	data := <-in
	fmt.Println("DATA:", data)
}

// MultiHash is
func MultiHash(in, out chan interface{}) {
	data := <-in
	fmt.Println("DATA:", data)
}

// CombineResults is
func CombineResults(in, out chan interface{}) {

}
