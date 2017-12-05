package main

import (
	"github.com/jshodd/guestMapper/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/clear", handlers.ClearDB).Methods("POST")
	router.HandleFunc("/add", handlers.Add).Methods("POST")
	router.HandleFunc("/generate-test-data", handlers.GenerateTest).Methods("POST")
	router.HandleFunc("/graph", handlers.GraphDB).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	http.Handle("/",router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
