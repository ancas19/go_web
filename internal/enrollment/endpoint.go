package enrollment

import (
	"courses/internal/commons"
	"encoding/json"
	"net/http"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		Create Controller
	}

	EnrollmentRequest struct {
		UserId   string `json:"userId"`
		CourseId string `json:"courseId"`
	}
)

func MakeEnpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndPoint(s),
	}
}

func makeCreateEndPoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req EnrollmentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		}
		if req.CourseId == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&commons.GeneralResponse{Status: 400, Error: "Course id is required"})
			return
		}
		if req.UserId == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&commons.GeneralResponse{Status: 400, Error: "User id is required"})
			return
		}

		enrollmentCreated, err := s.Create(req)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&commons.GeneralResponse{Status: 400, Error: err.Error()})
			return
		}
		json.NewEncoder(w).Encode(commons.GeneralResponse{Data: enrollmentCreated, Status: 202})

	}

}
