package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Handler struct {
	DB *sql.DB
}

type Response struct {
	Error    string      `json:"error,omitempty"`
	Response interface{} `json:"response,omitempty"`
}

type Tables struct {
}

type Fields struct {
	Name string
}

func NewDbExplorer(db *sql.DB) (http.Handler, error) {
	h := Handler{}
	h.DB = db
	return &h, nil
}

func (h *Handler) handleTables(w http.ResponseWriter, r *http.Request) {
	query := "SHOW TABLES;"
	rows, err := h.DB.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var tables []string

	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			panic(err)
		}
		tables = append(tables, table)
	}

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
	id        int
}

func parseQuery(path string) *getQuery {
	params := strings.Split(path, "?")
	fmt.Println(params)
	if len(params) == 1 {
		return &getQuery{
			tableName: params[0],
		}
	}
	return nil
}

func (h *Handler) handleShow(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) < 1 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	q := parseQuery(r.URL.Path[1:])
	fmt.Println(q)
	rows, err := h.DB.Query("SHOW FULL COLUMNS FROM " + q.tableName)
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), "doesn't exist") {
			j, err := json.Marshal(map[string]interface{}{"error": "unknown table"})
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			w.Write(j)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	for rows.Next() {

	}
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
	default:
		http.Error(w, "bad request", http.StatusBadRequest)
	}
}
