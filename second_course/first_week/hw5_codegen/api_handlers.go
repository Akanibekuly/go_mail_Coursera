package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
		return fmt.Errorf(strings.ToLower("Login") + " must me not empty")
	}

	return nil
}

func (in *CreateParams) validator(r *http.Request) error {

	in.Login = r.FormValue(strings.ToLower("Login"))

	if in.Login == "" {
		return fmt.Errorf(strings.ToLower("Login") + " must me not empty")
	}

	if len(in.Login) < 10 {
		return fmt.Errorf(strings.ToLower("Login") + " len must be >= 10")
	}

	in.Name = r.FormValue("full_name")

	in.Status = "user"
	if r.FormValue(strings.ToLower("Status")) != "" {
		in.Status = r.FormValue(strings.ToLower("Status"))
	}

	tArr := strings.Split("user|moderator|admin", "|")
	isConatin := false
	for _, v := range tArr {
		if v == in.Status {
			isConatin = true
			break
		}
	}
	if !isConatin {
		return fmt.Errorf(strings.ToLower("Status") + " must be one of [" + strings.Join(tArr, ", ") + "]")
	}

	temp, err := strconv.Atoi(r.FormValue(strings.ToLower("Age")))
	in.Age = temp
	if err != nil {
		return fmt.Errorf(strings.ToLower("Age") + " must be int")
	}

	if in.Age < 0 {
		return fmt.Errorf(strings.ToLower("Age") + " must be >= 0")
	}

	if in.Age > 128 {
		return fmt.Errorf(strings.ToLower("Age") + " must be <= 128")
	}
	return nil
}

func (srv *MyApi) handleProfile(w http.ResponseWriter, r *http.Request) {

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
	if newErr := json.NewEncoder(w).Encode(response{
		Error:    "",
		Response: u,
	}); newErr != nil {
		panic(err)
	}
}

func (srv *MyApi) handleCreate(w http.ResponseWriter, r *http.Request) {

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
	if newErr := json.NewEncoder(w).Encode(response{
		Error:    "",
		Response: u,
	}); newErr != nil {
		panic(err)
	}
}

// teststs 
func (in *OtherCreateParams) validator(r *http.Request) error {

	in.Username = r.FormValue(strings.ToLower("Username"))

	if in.Username == "" {
		return fmt.Errorf(strings.ToLower("Username") + " must me not empty")
	}

	if len(in.Username) < 3 {
		return fmt.Errorf(strings.ToLower("Username") + " len must be >= 3")
	}

	in.Name = r.FormValue("account_name")

	in.Class = "warrior"
	if r.FormValue(strings.ToLower("Class")) != "" {
		in.Class = r.FormValue(strings.ToLower("Class"))
	}

	tArr := strings.Split("warrior|sorcerer|rouge", "|")
	isConatin := false
	for _, v := range tArr {
		if v == in.Class {
			isConatin = true
			break
		}
	}
	if !isConatin {
		return fmt.Errorf(strings.ToLower("Class") + " must be one of [" + strings.Join(tArr, ", ") + "]")
	}

	temp, err := strconv.Atoi(r.FormValue(strings.ToLower("Level")))
	in.Level = temp
	if err != nil {
		return fmt.Errorf(strings.ToLower("Level") + " must be int")
	}

	if in.Level < 1 {
		return fmt.Errorf(strings.ToLower("Level") + " must be >= 1")
	}

	if in.Level > 50 {
		return fmt.Errorf(strings.ToLower("Level") + " must be <= 50")
	}
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
