package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID       int
	Username string
	phone    string
}

var jsonStr = `{"id": 42, "username": "Akzhol", "phone": "123"}`

func main() {
	fmt.Println(jsonStr)
	data := []byte(jsonStr)

	u := &User{}
	json.Unmarshal(data, u)
	fmt.Printf("struct:\n\t%#v\n\n", u)

	u.phone = "87478440780"
	result, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("json string: \n\t%s\n", string(result))
}
