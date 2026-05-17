package users

import (
	"courses/internal/commons"
	"courses/pkg/meta"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
		path := mux.Vars(r)
		id := path["id"]
		err := s.Delete(id)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(commons.GeneralResponse{Error: err.Error(), Status: 400})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func makeCreateEndPoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(commons.GeneralResponse{Error: "invalid request format"})
			return
		}
		if req.Firstname == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(commons.GeneralResponse{Error: "Firstname si required"})
			return
		}
		if req.Lastname == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(commons.GeneralResponse{Error: "Lastname si required"})
			return
		}
		userCreated, err := s.Create(req)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(commons.GeneralResponse{Error: err.Error()})
			return
		}
		json.NewEncoder(w).Encode(commons.GeneralResponse{Data: userCreated, Status: 202})
	}
}

func makeUpdateEndPoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(commons.GeneralResponse{Error: "invalid request format"})
			return
		}
		path := mux.Vars(r)
		id := path["id"]
		userUpdated, err := s.Update(id, req)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(commons.GeneralResponse{Error: err.Error()})
			return
		}
		json.NewEncoder(w).Encode(commons.GeneralResponse{Data: userUpdated, Status: 200})
	}
}

func makeGetEndPoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		userFound, err := s.GetById(id)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(commons.GeneralResponse{Error: err.Error()})
			return
		}
		json.NewEncoder(w).Encode(commons.GeneralResponse{Data: userFound, Status: 200})
	}
}

func makeGetAllEndPoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		filters := Filters{
			Firtsname: params.Get("first_name"),
			Email:     params.Get("email"),
		}
		limit, _ := strconv.ParseInt(params.Get("limit"), 10, 64)
		page, _ := strconv.ParseInt(params.Get("page"), 10, 64)
		totalUsers, err := s.Count(filters)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(commons.GeneralResponse{Error: err.Error()})
			return
		}

		meta, _ := meta.New(totalUsers, page, limit)
		usersFound, err := s.GetAll(filters, meta.Offset(), meta.Limit())
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(commons.GeneralResponse{Error: err.Error()})
			return
		}
		json.NewEncoder(w).Encode(commons.GeneralResponse{Data: usersFound, Status: 200, Meta: meta})
	}
}
