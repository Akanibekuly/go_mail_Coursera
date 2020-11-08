package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// const filePath string = "./data/users.txt"

// func main() {
// 	FastSearch(os.Stdout)
// }

type Data struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:"-"`
	Country  string   `json:"-"`
	Email    string   `json:"email"`
	Job      string   `json:"-"`
	Name     string   `json:"name"`
	Phone    string   `json:"-"`
}

var dataPool = sync.Pool{
	New: func() interface{} {
		return &Data{}
	},
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	seenBrowsers := make(map[string]bool, 200)
	var isAndroid, isMSIE bool
	in := bufio.NewScanner(file)
	count := 0
	fmt.Fprintln(out, fmt.Sprintf("found users:"))
	for in.Scan() {
		count++

		data := in.Bytes() // читаем по строчно из файла
		if bytes.Contains(data, []byte("Android")) == false && bytes.Contains(data, []byte("MSIE")) == false {
			continue
		}

		d := dataPool.Get().(*Data)
		err := json.Unmarshal(data, d)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// fmt.Printf("struct:\n\t%#v\n\n", d)
		isAndroid = false
		isMSIE = false
		for _, browser := range d.Browsers {

			if strings.Contains(browser, "Android") {
				isAndroid = true
				if _, ok := seenBrowsers[browser]; !ok {
					seenBrowsers[browser] = true
				}
			}

			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				if _, ok := seenBrowsers[browser]; !ok {
					seenBrowsers[browser] = true
				}
			}
		}
		dataPool.Put(d)
		if !(isAndroid && isMSIE) {
			continue
		}

		email := strings.Replace(d.Email, "@", " [at] ", -1)
		fmt.Fprintln(out, fmt.Sprintf("[%d] %s <%s>", count-1, d.Name, email))
	}

	// fmt.Fprintln(out, "")
	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}
