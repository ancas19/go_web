package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)
	Endpoints  struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateUserRequest struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	ErrorResponse struct {
		Error string `json:"error"`
	}
)

func MakeEnpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndPoint(s),
		Delete: makeDeleteEndPoint(s),
		Get:    makeGetEndPoint(s),
		GetAll: makeGetAllEndPoint(s),
		Update: makeUpdateEndPoint(s),
	}
}

func makeDeleteEndPoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("delete user")
		time.Sleep(5 * time.Second)
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}

func makeCreateEndPoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request format"})
			return
		}
		if req.Firstname == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Firstname si required"})
			return
		}
		if req.Lastname == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Lastname si required"})
			return
		}
		err = s.Create(req)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
			return
		}
		json.NewEncoder(w).Encode(req)
	}
}

func makeUpdateEndPoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("update user")
		time.Sleep(5 * time.Second)
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}

func makeGetEndPoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("get user")
		time.Sleep(5 * time.Second)
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}

func makeGetAllEndPoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Get All users")
		time.Sleep(5 * time.Second)
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}
