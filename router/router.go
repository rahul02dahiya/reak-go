package router

import (
	"github.com/gorilla/mux"
	"github.com/rahul/backend-go/handlers"
)

func SetupRouter() *mux.Router{
	r := mux.NewRouter()
	r.HandleFunc("/api/users",handlers.GetUsers).Methods("GET")
	r.HandleFunc("/api/users/{id}",handlers.GetUserByID).Methods("GET")
	r.HandleFunc("/api/users",handlers.CreateUsers).Methods("POST")
	r.HandleFunc("/api/users/{id}",handlers.UpdateUsers).Methods("PUT")
	r.HandleFunc("/api/users/{id}",handlers.DeleteUsers).Methods("DELETE")
	return r
}