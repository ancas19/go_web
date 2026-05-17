package main

import (
	"courses/internal/course"
	"courses/internal/users"
	"courses/pkg/bootstrap"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	router := mux.NewRouter()
	godotenv.Load()
	log := bootstrap.InitLogger()
	database, err := bootstrap.DBConnection()
	if err != nil {
		log.Fatal("failed to connect:", err)
	}

	userRepo := users.NewRepos(log, database)
	courseRepo := course.NewRepo(database, log)

	usersService := users.NewService(log, userRepo)
	courseService := course.NewService(log, courseRepo)

	userEnpoints := users.MakeEnpoints(usersService)
	courseEndpoits := course.MakeEnpoints(courseService)

	router.HandleFunc("/users", userEnpoints.Create).Methods("POST")
	router.HandleFunc("/users/{id}", userEnpoints.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEnpoints.Delete).Methods("DELETE")
	router.HandleFunc("/users/{id}", userEnpoints.Get).Methods("GET")
	router.HandleFunc("/users-all", userEnpoints.GetAll).Methods("GET")
	router.HandleFunc("/course", courseEndpoits.Create).Methods("POST")
	router.HandleFunc("/course-all", courseEndpoits.GetAll).Methods("GET")
	router.HandleFunc("/course/{id}", courseEndpoits.Get).Methods("GET")
	router.HandleFunc("/course/{id}", courseEndpoits.Delete).Methods("DELETE")

	server := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:3333",
	}
	log.Printf("Server started on http://%s\n", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
