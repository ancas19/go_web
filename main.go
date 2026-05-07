package main

import (
	"courses/internal/users"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	router := mux.NewRouter()
	godotenv.Load()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DATABSE_HOST"),
		os.Getenv("DATABSE_USER"),
		os.Getenv("DATABSE_PASS"),
		os.Getenv("DATABSE_NAME"),
		os.Getenv("DATABSE_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect:", err)
	}

	// Optional: configure the underlying connection pool
	database := db.Debug()
	if err != nil {
		log.Fatal(err)
	}
	log := log.New(os.Stdout, " ", log.LstdFlags|log.Lshortfile)
	userRepo := users.NewRepos(log, database)
	usersService := users.NewService(log, userRepo)
	userEnpoints := users.MakeEnpoints(usersService)
	database.AutoMigrate(&users.User{})
	router.HandleFunc("/users", userEnpoints.Create).Methods("POST")
	router.HandleFunc("/users", userEnpoints.Update).Methods("PUT")
	router.HandleFunc("/users", userEnpoints.Delete).Methods("DELETE")
	router.HandleFunc("/users", userEnpoints.Get).Methods("GET")
	router.HandleFunc("/users-all", userEnpoints.GetAll).Methods("GET")
	server := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:3333",
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
