package main

// race 1
import (
	"fmt"
	"sync"
)

func main() {
	var counters = map[int]int{}
	my := &sync.Mutex{}
	for i := 0; i < 5; i++ {
		go func(counters map[int]int, th int, mu *sync.Mutex) {
			for j := 0; j < 5; j++ {
				my.Lock()
				counters[th*10+j]++
				my.Unlock()
			}
		}(counters, i, my)
	}
	fmt.Scanln()
	my.Lock()
	fmt.Println("counters result", counters)
	my.Unlock()
}
