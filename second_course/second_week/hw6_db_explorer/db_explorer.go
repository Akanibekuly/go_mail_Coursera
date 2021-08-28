package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

func parseQuery(r *http.Request) *getQuery {
	arr := strings.Split(r.URL.Path, "/")
	if len(arr) != 2 {
		return nil
	}
	query := &getQuery{
		tableName: arr[1],
		limit:     5,
	}
	var err error
	q := r.URL.Query()
	if strID, ok := q["id"]; ok {
		query.id, err = strconv.Atoi(strID[0])
		if err != nil {
			log.Printf("query error: id should be int: %s\n", err)
			return nil
		}
	}

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

	rows, err := h.DB.Query("SHOW FULL COLUMNS FROM " + q.tableName)
	if err != nil {
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
