package router

import (
	"github.com/gorilla/mux"
	"github.com/rahul/backend-go/handlers"
)

func SetupRouter() *mux.Router{
	r := mux.NewRouter()
	r.HandleFunc("/users",handlers.GetUsers).Methods("GET")
	r.HandleFunc("/users",handlers.CreateUsers).Methods("POST")
	r.HandleFunc("/users/{id}",handlers.UpdateUsers).Methods("PUT")
	r.HandleFunc("/users/{id}",handlers.DeleteUsers).Methods("DELETE")
	return r
}