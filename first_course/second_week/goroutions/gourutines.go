package main

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	iterattionsNum = 7
	goroutinesNum  = 5
)

func doSomeWork(in int) {
	for j := 0; j < iterattionsNum; j++ {
		fmt.Printf(formatWork(in, j))
		runtime.Gosched()
	}
}

func main() {
	for i := 0; i < goroutinesNum; i++ {
		go doSomeWork(i)
	}
	fmt.Scanln()
}

func formatWork(in, i int) string {
	res := "" + strings.Repeat(" ", in) + "*"
	res += fmt.Sprintf("th %v inter %v ", in, i)
	res += strings.Repeat("*", i)
	res += "\n"
	return res
}
