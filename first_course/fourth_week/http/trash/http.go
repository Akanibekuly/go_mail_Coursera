package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world!")
	w.Write([]byte("wtf!!"))
}

func main() {
	http.HandleFunc("/", handler)

	fmt.Println("Strating server on port: 8080")
	http.ListenAndServe(":8080", nil)
}
