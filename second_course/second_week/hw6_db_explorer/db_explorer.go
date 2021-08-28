package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Handler struct {
	DB     *sql.DB
	Tables map[string]*Table
}

type Response struct {
	Error    string      `json:"error,omitempty"`
	Response interface{} `json:"response,omitempty"`
}

type Table struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name            string
	Type            string
	IsPrimary       bool
	IsNullable      bool
	IsAutoincrement bool
}

func NewDbExplorer(db *sql.DB) (http.Handler, error) {
	h := Handler{}
	h.DB = db
	h.Tables = make(map[string]*Table)

	if err := h.getAllTabels(); err != nil {
		return nil, err
	}

	if err := h.getTableFields(); err != nil {
		return nil, err
	}

	return &h, nil
}

func (h *Handler) getAllTabels() error {
	rows, err := h.DB.Query("SHOW TABLES")
	if err != nil {
		log.Printf("error with query all tables info: %s\n", err)
		return err
	}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			log.Printf("scanning error with query all tables info: %s\n", err)
			return err
		}

		h.Tables[name] = &Table{
			Name: name,
		}
	}

	return nil
}

func (h *Handler) getTableFields() error {
	t := func(i string) string { // is func that converts type of sql type name into golang type name
		if strings.HasPrefix(i, "int") {
			return "int"
		} else if strings.HasPrefix(i, "varchar") || i == "text" {
			return "string"
		} else {
			panic("unsupported field type")
		}
	} // end of t

	for _, table := range h.Tables {
		query := fmt.Sprintf("SHOW FULL COLUMNS FROM %s", table.Name)
		stm, err := h.DB.Prepare(query) // we need to prepare cause of sql injection cases
		if err != nil {
			log.Printf("prepare error show full columns: %s\n", err)
			return err
		}
		defer stm.Close() // always close statements

		rows, err := stm.Query()
		if err != nil {
			log.Printf("query error show full columns: %s\n", err)
			return err
		}
		defer rows.Close() // and close rows

		var w interface{}
		var isNullable, isPrimary, isAutoIncrement string
		for rows.Next() {
			var f Field
			err := rows.Scan(&f.Name, &f.Type, &w, &isNullable, &isPrimary, &w, &isAutoIncrement, &w, &w)
			if err != nil {
				log.Printf("scan error show full columns %s: %s\n", table.Name, err)
				return err
			}
			f.Type = t(f.Type)
			if isNullable == "YES" {
				f.IsNullable = true
			}
			if isPrimary == "PRI" {
				f.IsPrimary = true
			}
			if isAutoIncrement == "auto_increment" {
				f.IsAutoincrement = true
			}

			table.Fields = append(table.Fields, f)
		}
	}

	return nil
}

// GET / - возвращает список все таблиц (которые мы можем использовать в дальнейших запросах)
func (h *Handler) handleTables(w http.ResponseWriter, r *http.Request) {
	tables := make([]string, 0, len(h.Tables))

	// will return all the tables that we got when we start the server
	for name, _ := range h.Tables {
		tables = append(tables, name)
	}

	sort.Strings(tables)

	j, err := json.Marshal(Response{"", map[string]interface{}{"tables": tables}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(j)
}

type getQuery struct {
	tableName string
	limit     int
	offset    int
	id        *int
}

func parseQuery(r *http.Request) *getQuery {
	arr := strings.Split(r.URL.Path, "/")

	query := &getQuery{
		tableName: arr[1],
		limit:     5,
	}
	if r.Method == "PUT" {
		return query
	}

	switch len(arr) {
	case 3:
		strID := arr[2]
		id, err := strconv.Atoi(strID)
		query.id = &id
		if err != nil {
			log.Printf("query error: id should be int: %s\n", err)
			return nil
		}
		return query
	case 2:
		break
	default:
		return nil
	}

	var err error
	q := r.URL.Query()
	if offsetStr, ok := q["offset"]; ok {
		query.offset, err = strconv.Atoi(offsetStr[0])
		if err != nil {
			log.Printf("query error: offset should be int: %s\n", err)
			return nil
		}
	}

	if limitStr, ok := q["limit"]; ok {
		query.limit, err = strconv.Atoi(limitStr[0])
		if err != nil {
			log.Printf("query error: limit should be int: %s\n", err)
			return nil
		}
	}

	return query
}

// GET /$table?limit=5&offset=7 - возвращает список из 5 записей (limit) начиная с 7-й (offset) из таблицы $table. limit по-умолчанию 5, offset 0
// GET /$table/$id - возвращает информацию о самой записи или 404
func (h *Handler) handleShow(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) < 1 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	q := parseQuery(r)
	if q == nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	table, ok := h.Tables[q.tableName]
	if !ok {
		j, err := json.Marshal(map[string]interface{}{"error": "unknown table"})
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(j)
		return
	}

	if q.id != nil { // if id exists in the query
		var primaryKey string
		for _, f := range table.Fields {
			if f.IsPrimary {
				primaryKey = f.Name
				break
			}
		}

		query := fmt.Sprintf("SELECT * FROM %s WHERE %s=?", q.tableName, primaryKey)
		columns := make([]interface{}, len(table.Fields))
		columnsPtr := make([]interface{}, len(columns))
		for i := range columns {
			columnsPtr[i] = &columns[i]
		}
		err := h.DB.QueryRow(query, *q.id).Scan(columnsPtr...)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				j, err := json.Marshal(Response{"record not found", nil})
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				w.WriteHeader(http.StatusNotFound)
				w.Write(j)
				return
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		record := make(map[string]interface{})
		for i, f := range table.Fields {
			colNmae := strings.ToLower(f.Name)
			value := columns[i]
			bytes, ok := columns[i].([]byte)
			if ok {
				value = string(bytes)
			}
			record[colNmae] = value
		}

		j, err := json.Marshal(Response{"", map[string]interface{}{"record": record}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(j)
		return
	}

	stm, err := h.DB.Prepare(fmt.Sprintf("SELECT * FROM %s LIMIT ?, ?", q.tableName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stm.Close()

	rows, err := stm.Query(q.offset, q.limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	columns := make([]interface{}, len(table.Fields))
	columnsPtr := make([]interface{}, len(columns))
	for i := range columns {
		columnsPtr[i] = &columns[i]
	}

	var records []map[string]interface{}
	for rows.Next() {
		err := rows.Scan(columnsPtr...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		record := make(map[string]interface{})
		for i, f := range table.Fields {
			colName := strings.ToLower(f.Name)
			value := columns[i]
			columns[i] = nil

			bytes, ok := value.([]byte)
			if ok {
				strValue := string(bytes)
				switch f.Type {
				case "int":
					intValue, err := strconv.Atoi(strValue)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					value = intValue
				case "string":
					value = strValue
				}
			}
			record[colName] = value
		}

		records = append(records, record)
	}

	j, err := json.Marshal(Response{"", map[string]interface{}{"records": records}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(j)
}

// PUT /$table - создаёт новую запись, данный по записи в теле запроса (POST-параметры)
func (h *Handler) handlePut(w http.ResponseWriter, r *http.Request) {
	q := parseQuery(r)
	t, exists := h.Tables[q.tableName]
	if q.id != nil || !exists {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	jsonBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	body := make(map[string]interface{})
	if err = json.Unmarshal(jsonBody, &body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var primaryKey string
	var keys, placeholders []string
	var values []interface{}
	for _, f := range t.Fields {
		if f.IsPrimary {
			primaryKey = f.Name
			continue
		}

		value, ok := body[f.Name]
		if !ok {
			if f.IsNullable {
				value = nil
			} else {
				switch f.Type {
				case "string":
					value = ""
				case "int":
					value = 0
				default:
					http.Error(w, "unknown type", http.StatusInternalServerError)
					return
				}
			}
		} else {
			switch value.(type) {
			case string:
				if f.Type == "int" {
					invalidField(f.Name, w, r)
					return
				}
			case float64:
				if f.Type == "string" {
					invalidField(f.Name, w, r)
					return
				}
			default:
				invalidField(f.Name, w, r)
				return
			}
		}

		keys = append(keys, f.Name)
		placeholders = append(placeholders, "?")
		values = append(values, value)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", q.tableName, strings.Join(keys, ","), strings.Join(placeholders, ","))
	stm, err := h.DB.Prepare(query)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stm.Close()

	res, err := stm.Exec(values...)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{primaryKey: id}
	j, err := json.Marshal(Response{"", response})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(j)
}

func invalidField(field string, w http.ResponseWriter, r *http.Request) {
	errMsg := fmt.Sprintf("field %s have invalid type", field)
	c := Response{errMsg, nil} // wrap in 'response'
	j, err := json.Marshal(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return // ?
	}

	w.WriteHeader(http.StatusBadRequest) // 400
	w.Write(j)                           // the content
}

// POST /$table/$id - обновляет запись, данные приходят в теле запроса (POST-параметры)
func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request) {
	q := parseQuery(r)

	t, exists := h.Tables[q.tableName]
	if q.id == nil || !exists {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	jsonBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	body := make(map[string]interface{})
	if err = json.Unmarshal(jsonBody, &body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var primaryKey string
	var keys []string
	var values []interface{}
	for _, f := range t.Fields {
		if f.IsPrimary {
			primaryKey = f.Name
		}
		value, exists := body[f.Name]
		if !exists {
			continue
		}

		if f.IsPrimary {
			invalidField(f.Name, w, r)
			return
		}

		switch value.(type) {
		case string:
			if f.Type == "int" {
				invalidField(f.Name, w, r)
				return
			}
		case float64:
			if f.Type == "string" {
				invalidField(f.Name, w, r)
				return
			}
		case nil:
			if !f.IsNullable {
				invalidField(f.Name, w, r)
				return
			}
		default:
			invalidField(f.Name, w, r)
			return
		}

		keys = append(keys, fmt.Sprintf("%s = ?", f.Name))
		values = append(values, value)
	}

	values = append(values, q.id)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", q.tableName, strings.Join(keys, ","), primaryKey)
	stm, err := h.DB.Prepare(query)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stm.Close()

	res, err := stm.Exec(values...)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(Response{"", map[string]interface{}{"updated": count}})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(j)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	if r.Method == "GET" && r.URL.Path == "/" {
		h.handleTables(w, r)
		return
	}

	switch r.Method {
	case "GET":
		h.handleShow(w, r)
	case "PUT":
		h.handlePut(w, r)
	case "POST":
		h.handlePost(w, r)
	default:
		http.Error(w, "bad request", http.StatusBadRequest)
	}
}
