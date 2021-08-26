package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

type Handler struct {
	DB *sql.DB
}

func NewDbExplorer(db *sql.DB) (http.Handler, error) {
	h := Handler{}
	h.DB = db
	return &h, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()
}

func (h *Handler) hondleShow(w http.ResponseWriter, r *http.Request) {

}
