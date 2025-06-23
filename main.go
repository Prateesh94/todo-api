package main

import (
	"fmt"
	"log"
	"net/http"
	"todo-api/data"
	"todo-api/encrypt"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.Use(encrypt.Limitmid)
	router.HandleFunc("/register", data.AddNewUserEndpoint).Methods("POST")
	router.HandleFunc("/login", data.LoginEndpoint).Methods("POST")
	router.HandleFunc("/todos", data.AddTodoEndpoint).Methods("POST")
	router.HandleFunc("/todos/{id}", data.UpdateEndpoint).Methods("PUT")
	router.HandleFunc("/todos/{id}", data.DeleteEndpoint).Methods("DELETE")
	router.HandleFunc("/todos", data.FetchtodoEndpoint).Methods("GET")
	fmt.Println("Server listening on PORT:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
