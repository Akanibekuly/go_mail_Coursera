package main

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"
)

func main() {
	sess := &Session{
		Login:     "rvasily",
		Useragent: "Chrome",
	}

	dataJson, _ := json.Marshal(sess)

	fmt.Printf("dataJson\nlen %d\n%s\n", len(dataJson), dataJson)

	dataPb, _ := proto.Marshal(sess)

	fmt.Printf("dataPb\nlen %d\n%v\n", len(dataPb), dataPb)
}
