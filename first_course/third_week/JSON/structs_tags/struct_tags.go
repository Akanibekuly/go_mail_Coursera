package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID       int `json:"user_id,string"`
	Username string
	Address  string `json:",omitempty"`
	Company  string `json:"-"`
}

func main() {
	u := &User{
		ID:       42,
		Username: "Akanibekuly",
		Address:  "Astana",
		Company:  "Techno Partners",
	}
	result, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	fmt.Printf("json string %s\n", string(result))
}
