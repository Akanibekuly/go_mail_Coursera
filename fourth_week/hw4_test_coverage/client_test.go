package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// код писать тут
type TestCase struct {
	ID         string
	Response   string
	StatusCode int
}

type MyUser struct {
	Id        int    `xml: "id"`
	Name      string `xml: "-"`
	FirstName string `xml: "first_name"`
	LastName  string `xml: "last_name"`
	Age       int    `xml: "age"`
	About     string `xml: "about"`
	Gender    string `xml: "gender"`
}

type SearchServer struct {
	pathToFile string
}

// func main() {
// 	result, err := getUsersFromFile("dataset.xml")
// 	if err != nil {
// 		fmt.Println("Error happend ", err)
// 	}
// 	for _, u := range result {
// 		fmt.Printf("%#v", u)
// 		fmt.Println()
// 	}

// }

func (ss *SearchServer) GetUsers(params SearchRequest) ([]MyUser, error) {
	raw, err := getUsersFromFile(ss.pathToFile)
	if err != nil {
		return nil, err
	}

	var result []MyUser

	// фильтр по содержанию
	if params.Query != "" {
		for _, u := range raw {
			nameContains := strings.Contains(u.Name, params.Query)
			aboutContains := strings.Contains(u.About, params.Query)

			if nameContains || aboutContains {
				result = append(result, u)
			}
		}
	}

	// сортировка
	if params.OrderBy != 0 && params.OrderField != "" {
		sortUsers(result, params.OrderBy, params.OrderField)
	}

	// начало конец, если в сумме больше длины результата
	if params.Offset+params.Limit > len(result) {
		return result[params.Offset:], nil
	}

	return result[params.Offset:params.Limit], nil
}

func sortUsers(users []MyUser, orderBy int, orderField string) {
	sort.Slice(users, func(i, j int) bool {
		if orderField == "Id" {
			if orderBy == -1 {
				return users[i].Id > users[j].Id
			}
			return users[i].Id < users[j].Id

		} else if orderField == "Age" {
			if orderBy == -1 {
				return users[i].Age > users[j].Age
			}
			return users[i].Age < users[j].Age
		} else if orderField == "Name" {
			if orderBy == -1 {
				return users[i].Name > users[j].Name
			}
			return users[i].Name < users[j].Name
		}
		return users[i].Id > users[j].Id
	})
}

func getUsersFromFile(pathToFile string) ([]MyUser, error) {
	var result = []MyUser{}
	var temp = MyUser{}

	file, err := os.Open(pathToFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := xml.NewDecoder(file)

	for {
		tok, tokenErr := decoder.Token()
		if tokenErr != nil && tokenErr != io.EOF {
			fmt.Println("error happend", tokenErr)
			return nil, tokenErr
		} else if tokenErr == io.EOF {
			break
		}
		if tok == nil {
			fmt.Println("t is nil break")
		}

		switch tok := tok.(type) {
		case xml.StartElement:
			if tok.Name.Local == "id" {
				if err := decoder.DecodeElement(&temp.Id, &tok); err != nil {
					return nil, err
				}
			}
			if tok.Name.Local == "first_name" {
				if err := decoder.DecodeElement(&temp.FirstName, &tok); err != nil {
					return nil, err
				}
			}
			if tok.Name.Local == "last_name" {
				if err := decoder.DecodeElement(&temp.LastName, &tok); err != nil {
					return nil, err
				}
				temp.Name = temp.FirstName + " " + temp.LastName
			}
			if tok.Name.Local == "age" {
				if err := decoder.DecodeElement(&temp.Age, &tok); err != nil {
					return nil, err
				}
			}
			if tok.Name.Local == "about" {
				if err := decoder.DecodeElement(&temp.About, &tok); err != nil {
					return nil, err
				}
			}
			if tok.Name.Local == "gender" {
				if err := decoder.DecodeElement(&temp.Gender, &tok); err != nil {
					return nil, err
				}
				result = append(result, temp)
			}
		}
	}
	return result, nil
}
