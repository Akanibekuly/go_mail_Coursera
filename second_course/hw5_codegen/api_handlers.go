package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// response that we return back
type response struct {
	Error    string      `json:"error"`
	Response interface{} `json:"response,omitempty"`
}

func (in *ProfileParams) validator(r *http.Request) error {

	in.Login = r.FormValue(strings.ToLower("Login"))
	if in.Login == "" {
		return fmt.Errorf("login must me not empty")
	}

	return nil
}

func (in *CreateParams) validator(r *http.Request) error {

	in.Login = r.FormValue(strings.ToLower("Login"))
	if in.Login == "" {
		return fmt.Errorf("login must me not empty")
	}

	in.Name = r.FormValue(strings.ToLower("Name"))

	in.Status = r.FormValue("user")

	in.Age = r.FormValue(strings.ToLower("Age"))

	return nil
}

func (srv *ProfileParams) handleProfile(w http.ResponseWriter, r *http.Request) {

	p := ProfileParams{}
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
	u, err := srv.Profile(r.Context(), p)
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

func (srv *CreateParams) handleCreate(w http.ResponseWriter, r *http.Request) {

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

	p := CreateParams{}
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
	u, err := srv.Create(r.Context(), p)
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

func (in *OtherCreateParams) validator(r *http.Request) error {

	in.Username = r.FormValue(strings.ToLower("Username"))
	if in.Username == "" {
		return fmt.Errorf("login must me not empty")
	}

	in.Name = r.FormValue(strings.ToLower("Name"))

	in.Class = r.FormValue("warrior")

	in.Level = r.FormValue(strings.ToLower("Level"))

	return nil
}

func (srv *OtherCreateParams) handleCreate(w http.ResponseWriter, r *http.Request) {

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

	p := OtherCreateParams{}
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
	u, err := srv.Create(r.Context(), p)
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

// MyApi server router
func (srv *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {

	case "/user/profile":
		srv.handleProfile(w, r)

	case "/user/create":
		srv.handleCreate(w, r)

	default:
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(response{
			Error: "unknown method",
		}); err != nil {
			panic(err)
		}
	}
}

// OtherApi server router
func (srv *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {

	case "/user/create":
		srv.handleCreate(w, r)

	default:
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(response{
			Error: "unknown method",
		}); err != nil {
			panic(err)
		}
	}
}
