package course

import (
	"courses/internal/commons"
	"courses/pkg/meta"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateCourseReq struct {
		Name      string `json:"name"`
		StartDate string `json:"startDate"`
		EndtDate  string `json:"endDate"`
	}
)

func MakeEnpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoints(s),
		Delete: makeDeleteEndPoint(s),
		Get:    makeGetEndPoint(s),
		GetAll: makeGetAllEndPoint(s),
	}
}

func makeCreateEndpoints(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateCourseReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&commons.GeneralResponse{Status: 400, Error: fmt.Sprintf("Invalid request format")})
			return
		}
		if req.Name == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&commons.GeneralResponse{Status: 400, Error: fmt.Sprintf("Name is required")})
			return
		}
		if req.StartDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&commons.GeneralResponse{Status: 400, Error: fmt.Sprintf("StartDate is required")})
			return
		}
		if req.EndtDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&commons.GeneralResponse{Status: 400, Error: fmt.Sprintf("EndDate is required")})
			return
		}

		courseCreated, err := s.Create(req)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&commons.GeneralResponse{Status: 400, Error: err.Error()})
			return
		}
		json.NewEncoder(w).Encode(commons.GeneralResponse{Data: courseCreated, Status: 202})
	}
}

func makeDeleteEndPoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		idCourse := path["id"]
		if idCourse == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&commons.GeneralResponse{Status: 400, Error: fmt.Sprintf("EndDate is required")})
			return
		}
		err := s.Delete(idCourse)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&commons.GeneralResponse{Status: 400, Error: fmt.Sprintf("EndDate is required")})
			return
		}
		w.WriteHeader(http.StatusNoContent)
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
		nameCourse := params.Get("name")
		limit, _ := strconv.ParseInt(params.Get("limit"), 10, 64)
		page, _ := strconv.ParseInt(params.Get("page"), 10, 64)
		totalUsers, err := s.Count(nameCourse)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(commons.GeneralResponse{Error: err.Error()})
			return
		}

		meta, _ := meta.New(totalUsers, page, limit)
		usersFound, err := s.GetAllCourses(nameCourse, meta.Offset(), meta.Limit())
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(commons.GeneralResponse{Error: err.Error()})
			return
		}
		json.NewEncoder(w).Encode(commons.GeneralResponse{Data: usersFound, Status: 200, Meta: meta})
	}
}
