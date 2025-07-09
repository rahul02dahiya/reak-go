package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rahul/backend-go/db"
	"github.com/rahul/backend-go/models"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {

	rows, err := db.DB.Query("SELECT id, name, mobile, email FROM users;")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Mobile, &u.Email); err != nil {
			http.Error(w, "Error while fetching data", http.StatusInternalServerError)
		}
		users = append(users, u)
	}

	json.NewEncoder(w).Encode(users)
}

func CreateUsers(w http.ResponseWriter, r *http.Request) {
	var u models.User
	json.NewDecoder(r.Body).Decode(&u)
	sqlStmt := `INSERT INTO users (name, mobile, email) VALUES (?, ?, ?);`
	result, err := db.DB.Exec(sqlStmt, u.Name, u.Mobile, u.Email)
	if err != nil {

		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			parts := strings.Split(err.Error(), ".")
			field := parts[1]
			switch field {
			case "email":
				http.Error(w, "Email already registered", http.StatusInternalServerError)
				return
			case "mobile":
				http.Error(w, "Mobile number already registered", http.StatusInternalServerError)
				return
			default:
				http.Error(w, "Details already registered", http.StatusInternalServerError)
				return
			}
		}

		http.Error(w, "Unable to create user", http.StatusInternalServerError)
		return
	}
	id, _ := result.LastInsertId()
	u.ID = int(id)
	json.NewEncoder(w).Encode(u)
}

func UpdateUsers(w http.ResponseWriter, r *http.Request) {

	var u models.User
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	json.NewDecoder(r.Body).Decode(&u)
	sqlStmt := `UPDATE users SET name = ?, mobile = ?, email = ? where id = ? ;`
	_, err = db.DB.Exec(sqlStmt, u.Name, u.Mobile, u.Email, id)
	if err != nil {
		http.Error(w, "Unable to update user", http.StatusInternalServerError)
		return
	}
	u.ID = id
	json.NewEncoder(w).Encode(u)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {

	var u models.User
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	err = db.DB.QueryRow(`SELECT id, name, mobile, email FROM users WHERE id=? ;`, id).Scan(&u.ID, &u.Name, &u.Mobile, &u.Email)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(u)
}

func DeleteUsers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	sqlStmt := `DELETE FROM users WHERE id=? ;`
	_, err = db.DB.Exec(sqlStmt, id)
	if err != nil {
		http.Error(w, "Unable to detele user", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("User deleted"))
}
