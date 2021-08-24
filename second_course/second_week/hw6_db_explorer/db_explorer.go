package main

import (
	"database/sql"
	"net/http"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные
func NewDbExplorer(db *sql.DB) http.Handler{

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

	})
}