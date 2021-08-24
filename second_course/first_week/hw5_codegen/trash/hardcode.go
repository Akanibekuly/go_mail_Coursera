package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// response that we return back
type response struct {
	Error    string      `json:"error"`
	Response interface{} `json:"response,omitempty"`
}

// Vaidator validates ProfileParams
func (in *ProfileParams) validator(r *http.Request) error {
	in.Login = r.FormValue("login")
	if in.Login == "" {
		return fmt.Errorf("login must me not empty")
	}
	return nil
}

// Validator validates CreateParams
func (in *CreateParams) validator(r *http.Request) error {
	in.Login = r.FormValue("login")
	if in.Login == "" {
		return fmt.Errorf("login must me not empty")
	}
	if len(in.Login) < 10 {
		return fmt.Errorf("login len must be >= 10")
	}

	in.Name = r.FormValue("full_name")

	if r.FormValue("status") == "" {
		in.Status = "user"
	} else {
		in.Status = r.FormValue("status")
	}
	if !(in.Status == "user" || in.Status == "moderator" || in.Status == "admin") {
		return fmt.Errorf("status must be one of [user, moderator, admin]")
	}
	age, err := strconv.Atoi(r.FormValue("age"))
	if err != nil {
		return fmt.Errorf("age must be int")
	}
	if age < 0 {
		return fmt.Errorf("age must be >= 0")
	}
	if age > 128 {
		return fmt.Errorf("age must be <= 128")
	}
	in.Age = age
	return nil
}

// ServeHTTP routes api for MyApi
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

// handles Profile data
func (srv *MyApi) handleProfile(w http.ResponseWriter, r *http.Request) {
	//TODO params
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

	ctx := r.Context()
	u, err := srv.Profile(ctx, p)
	if err != nil {
		//type assertions
		v, ok := err.(ApiError)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(v.HTTPStatus)
		}
		if newErr := json.NewEncoder(w).Encode(response{
			Error: err.Error(),
		}); newErr != nil {
			panic(err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response{
		Error:    "",
		Response: u,
	}); err != nil {
		panic(err)
	}
}

func (srv *MyApi) handleCreate(w http.ResponseWriter, r *http.Request) {
	// TODO: authorization
	if key := r.Header.Get("X-Auth"); key != "100500" {
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(response{
			Error: "unauthorized",
		}); err != nil {
			panic(err)
		}
		return
	}
	// TODO: method validation
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotAcceptable)
		if err := json.NewEncoder(w).Encode(response{
			Error: "bad method",
		}); err != nil {
			panic(err)
		}
		return
	}

	// TODO: params
	c := CreateParams{}
	err := c.validator(r)
	// handle errors from validator
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if newErr := json.NewEncoder(w).Encode(response{
			Error: err.Error(),
		}); newErr != nil {
			panic(err)
		}
		return
	}

	ctx := r.Context()
	u, err := srv.Create(ctx, c)
	// handle errors from create function
	if err != nil {
		v, ok := err.(ApiError)
		var status int
		if !ok {
			status = 500
		} else {
			status = v.HTTPStatus
		}
		w.WriteHeader(status)
		if newErr := json.NewEncoder(w).Encode(response{
			Error: err.Error(),
		}); newErr != nil {
			panic(err)
		}
		return
	}

	// if all is OK
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response{
		Error:    "",
		Response: u,
	}); err != nil {
		panic(err)
	}
}

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

func (in *OtherCreateParams) validator(r *http.Request) error {
	in.Username = r.FormValue("username")
	if in.Username == "" {
		return fmt.Errorf("username must me not empty")
	}
	if len(in.Username) < 3 {
		return fmt.Errorf("username len must be >= 3")
	}

	in.Name = r.FormValue("account_name")

	in.Class = "warrior"
	if r.FormValue("class") != "" {
		in.Class = r.FormValue("class")
	}

	if !(in.Class == "warrior" || in.Class == "sorcerer" || in.Class == "rouge") {
		return fmt.Errorf("class must be one of [warrior, sorcerer, rouge]")
	}

	level, err := strconv.Atoi(r.FormValue("level"))
	if err != nil {
		return fmt.Errorf("level must be int")
	}
	if level < 0 {
		return fmt.Errorf("level must be >= 1")
	}
	if level > 128 {
		return fmt.Errorf("level must be <= 50")
	}
	in.Level = level
	return nil
}

func (srv *OtherApi) handleCreate(w http.ResponseWriter, r *http.Request) {
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
