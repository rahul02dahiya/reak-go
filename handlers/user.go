package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rahul/backend-go/models"
)

var users []models.User = []models.User{
	{
		ID:     "1",
		Name:   "Rahul",
		Mobile: "9303392881",
		Email:  "rahul@gmail.com",
	},
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

func CreateUsers(w http.ResponseWriter, r *http.Request){
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	user.ID = "id"
	users = append(users, user)
	json.NewEncoder(w).Encode(user)
}
func UpdateUsers(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	for i,u := range users {
		if u.ID == params["id"] {
			var user models.User
			json.NewDecoder(r.Body).Decode(&user)
			users[i]=user
			return
		}
	}
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
}

func GetUserByID(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	for _, u := range users {
		if (u.ID == params["id"]) {
			json.NewEncoder(w).Encode(u)
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}


func DeleteUsers(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	for i, u := range users {
		if u.ID == params["id"] {
			users = append(users[:i],users[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}
