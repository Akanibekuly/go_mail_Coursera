package main

import (
	"fmt"
	"time"
)

func longSQLQuery() chan bool {
	ch := make(chan bool, 1)
	go func() {
		time.Sleep(2 * time.Second)
		ch <- true
	}()
	return ch
}

func main() {
	timer := time.NewTimer(3 * time.Second)
	select {
	case <-timer.C:
		fmt.Println("timer.C timout happend")
	case <-time.After(time.Minute):
		fmt.Println("time.After timeout happend")
	case result := <-longSQLQuery():
		// освобаждает ресурс
		if !timer.Stop() {
			<-timer.C
		}
		fmt.Println("operation result:", result)
	}
	
}
