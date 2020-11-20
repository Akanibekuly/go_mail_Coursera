package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

// код писать тут
type TestCase struct {
	ID         string
	Response   string
	StatusCode int
}

type MyUser struct {
	Id        int    `xml:"id" json:"id"`
	Name      string `xml:"-" json:"-"`
	FirstName string `xml:"first_name" json:"-"`
	LastName  string `xml:"last_name" json:"-"`
	Age       int    `xml:"age" json:"age"`
	About     string `xml:"about" json:"about"`
	Gender    string `xml:"gender" json:"gender"`
}

type SearchServer struct {
	pathToFile string
}

const (
	testToken = "1234567"
)

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

// SearchServerHandler Обработчик сервера
func SearchServerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// проверяем токен на наличие
	token := r.Header.Get("AccesToken")
	if token == "" || token != testToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// берем данные для поиска и обробатывем ошибки
	searchRequest, err := getValidInput(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf(`"StatusCode": 400, "Error": "%s"`, err.Error()))
		return
	}

	// создаем своего сервера поиска
	searchServer := SearchServer{"./dataset.xml"}

	users, err := searchServer.GetUsers(searchRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf(`"StatusCode": 500, "Error": "%s"`, err.Error()))
		return
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf(`"StatusCode: 500", "Error": "Invalid data for json encoding"`))
		return
	}

	io.WriteString(w, string(usersJSON))
}

// Обпработка правильности данных для поиска
func getValidInput(r *http.Request) (SearchRequest, error) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		return SearchRequest{}, errors.New("limit")
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		return SearchRequest{}, errors.New("offset")
	}

	orderBy, err := strconv.Atoi(r.URL.Query().Get("order_by"))
	if err != nil || orderBy < -1 || orderBy > 1 {
		return SearchRequest{}, errors.New("order_by")
	}

	orderField := r.URL.Query().Get("order_field")
	if !isValidOrderField(orderField) {
		return SearchRequest{}, errors.New("ErrorBadOrderField")
	}
	query := r.URL.Query().Get("query")

	return SearchRequest{
		Limit:      limit,
		Offset:     offset,
		OrderBy:    orderBy,
		Query:      query,
		OrderField: orderField,
	}, nil
}

// Обработка правильности полей для поиска
func isValidOrderField(orderField string) bool {
	if orderField == "Id" || orderField == "Name" || orderField == "Age" {
		return true
	}
	return false
}



