package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
	"text/template"
)

// структура для обработки и сохранения данных для обработчиков
type wrapper struct {
	In     string
	Name   string
	Url    string
	Auth   bool
	Method string
}

// структкра для создания валидационных функции
type tag struct {
	Required  bool
	Paramname string
	Default   string
	Enum      string
	Min       string
	Max       string
}

type param struct {
	Name string
	Type string
	Tags tag
}

type validate struct {
	Name   string
	Params []param
}

// структура для создания ServeHTTP роутинга для каждой из методов
type serveHTTP struct {
	Serve map[string][]wrapper
}

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	out, _ := os.Create(os.Args[2])

	fmt.Fprintln(out, `package `+node.Name.Name+`

	import (
		"encoding/json"
		"fmt"
		"net/http"
		"strconv"
	)
	
// response that we return back
type response struct {`)

	fmt.Fprintln(out, "\tError    string      `json:\"error\"`")
	fmt.Fprintln(out, "\tResponse interface{} `json:\"response,omitempty\"`")
	fmt.Fprintln(out, "}\n")
	// инициализируем мапу для серверов
	servers := serveHTTP{
		Serve: make(map[string][]wrapper),
	}

	for _, f := range node.Decls {

		// обработка функции и получения данных для обработки
		fnc, ok := f.(*ast.FuncDecl)
		if ok {
			if fnc.Doc != nil {

				// проверка коммнентов на содержание ключевого слова
				if !strings.HasPrefix(fnc.Doc.Text(), "apigen:api") {
					// fmt.Println("SKIP comments doesn't have \"apigen:api\" prefix")
					continue
				}

				w := wrapper{}
				w.Name = fnc.Name.Name
				in := (*fnc.Type.Params.List[1]).Type.(*ast.Ident).Name
				w.In = in
				u := struct {
					Url    string `json:"url"`
					Auth   bool   `json:"auth"`
					Method string `json:"method"`
				}{}
				index := strings.Index(fnc.Doc.List[0].Text, "{")
				jsonRaw := fnc.Doc.List[0].Text[index:]
				json.Unmarshal([]byte(jsonRaw), &u)
				// fmt.Println(jsonRaw, u)
				w.Auth = u.Auth
				w.Method = u.Method
				w.Url = u.Url
				// fmt.Println(w)
				temp := fnc.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name
				servers.Serve[temp] = append(servers.Serve[temp], w)
				handlerTpl.Execute(out, w)
			}
		}

		// обработка структур
		g, ok := f.(*ast.GenDecl)
		if !ok {
			// fmt.Printf("SKIP %T is not *ast.GenDecl\n", f)
			continue
		}
		// SPECS_LOOP:
		for _, spec := range g.Specs {
			currType, ok := spec.(*ast.TypeSpec)
			if !ok {
				// fmt.Printf("SKIP %T is not ast.TypeSpec\n", spec)
				continue
			}

			// type assertion for structure Type
			currStruct, ok := currType.Type.(*ast.StructType)
			if ok {
				// обработка структур для валидации с содержанием в имени Params
				if strings.Contains(currType.Name.Name, "Params") {
					v := validate{}
					v.Name = currType.Name.Name
					// переходим к параметрам структуры
					for _, field := range currStruct.Fields.List {
						fieldName := field.Names[0].Name
						fileType := field.Type.(*ast.Ident).Name
						p := param{}
						p.Name = fieldName
						p.Type = fileType
						// обробатываем тэги параметров
						if field.Tag != nil {
							// присваиваем значения тэгов
							tagsRaw := field.Tag.Value
							t := tag{}
							index := strings.Index(tagsRaw, "\"")
							// убираем префикс apivalidator:" и "
							//и полцчаем в чистом виде enum=user|moderator|admin,default=user
							tagsString := tagsRaw[index+1 : len(tagsRaw)-2]
							// fmt.Println(tagsString)
							// разделяем по параметром тэгов
							tagsArr := strings.Split(tagsString, ",")
							for _, val := range tagsArr {

								vArr := strings.Split(val, "=")
								// fmt.Print(vArr, "\t")
								if vArr[0] == "required" {
									t.Required = true
								}
								if vArr[0] == "paramname" {
									t.Paramname = vArr[1]
								}
								if vArr[0] == "enum" {
									t.Enum = vArr[1]
								}
								if vArr[0] == "default" {
									t.Default = vArr[1]
								}
								if vArr[0] == "min" {
									t.Min = vArr[1]
								}
								if vArr[0] == "max" {
									t.Max = vArr[1]
								}

							}

							p.Tags = t
						}
						v.Params = append(v.Params, p)
					}
					fmt.Println(v)
					err := validatorTpl.Execute(out, v)
					fmt.Println("DEBUG template error %s", err)
				}

			}

		}

	}

	// парсинг структур закончен
	serverTpl.Execute(out, servers)

}

var (
	handlerTpl = template.Must(template.New("handlerTpl").Parse(`
func (srv *{{.In}}) handle{{.Name}}(w http.ResponseWriter, r *http.Request) {
		{{if .Auth }}
	//TODO authorization
	if key := r.Header.Get("X-Auth"); key != "100500" {
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(response{
			Error: "unauthorized",
		}); err != nil {
			panic(err)
		}
		return
	}
	{{end}}
	{{if eq .Method "POST"}}
	//TODO method validation
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotAcceptable)
		if err := json.NewEncoder(w).Encode(response{
			Error: "bad method",
		}); err != nil {
			panic(err)
		}
		return
	}
	{{end}}
	p := {{.In}}{}
	err := p.validator(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(response{
			Error: err.Error(),
		}); err != nil {
			panic(err)
		}
		return
	}
	u, err := srv.{{.Name}}(r.Context(), p)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
	if newErr := json.NewEncoder(w).Encode(response{
		Error:    "",
		Response: u,
	}); newErr != nil {
		panic(err)
	}
}
`))

	validatorTpl = template.Must(template.New("handlerTpl").Parse(`
func (in *{{.Name}}) validator(r *http.Request) error {
	{{range .Params}}
	{{if ne .Tags.Default ""}}in.{{.Name}}=r.FormValue("{{.Tags.Default}}"){{else}}in.{{.Name}}=r.FormValue(strings.ToLower("{{.Name}}")){{end}}
	{{if .Tags.Required}}if in.{{.Name}}==""{
		return fmt.Errorf("login must me not empty")
	}{{end}}
	{{if ne .Tags.Min ""}}
		{{if eq .Type "string"}}
			if len(in.{{.Name}})<{{.Tags.Min}}{
				return fmt.Errorf("login len must be >= {{.Tags.Min}}")
			}
		{{else}}
			if in.Name
		{{end}}
	{{end}}
	{{end}}
	return nil
	}
	`))

	serverTpl = template.Must(template.New("serverrTpl").Parse(
		`{{ range $key, $value := .Serve }}
	// {{$key}} server router
	func (srv *{{$key}}) ServeHTTP(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
			{{range $value}}
			case "{{.Url}}":
				 srv.handle{{.Name}}(w, r)
			{{end}}
		 	default:
		 		w.WriteHeader(http.StatusNotFound)
		 		if err := json.NewEncoder(w).Encode(response{
		 			Error: "unknown method",
		 		}); err != nil {
		 			panic(err)
		 		}
		 	}
	}
	{{ end }}
	`))
)
