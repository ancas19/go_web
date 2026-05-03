package main

import (
	"courses/internal/users"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	usersService := users.NewService()
	userEnpoints := users.MakeEnpoints(usersService)
	router := mux.NewRouter()
	router.HandleFunc("/users", userEnpoints.Create).Methods("POST")
	router.HandleFunc("/users", userEnpoints.Update).Methods("PUT")
	router.HandleFunc("/users", userEnpoints.Delete).Methods("DELETE")
	router.HandleFunc("/users", userEnpoints.Get).Methods("GET")
	router.HandleFunc("/users-all", userEnpoints.GetAll).Methods("GET")
	server := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:3333",
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
