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

// Defines the structure for error responses
type ErrorResponse struct {
	Message string `json:"message"`
}

// sendError sends a JSON error response with the specified status code and message
func sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}

// Get all users from users DB and send 
func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, name, mobile, email FROM users;")
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Database error: failed to fetch users")
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Mobile, &u.Email); err != nil {
			sendError(w, http.StatusInternalServerError, "Error while fetching user data")
			return
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		sendError(w, http.StatusInternalServerError, "Error while iterating user data")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Getting data from frontend and save to users DB after validation
func CreateUsers(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if u.Name == "" || u.Mobile == "" || u.Email == "" {
		sendError(w, http.StatusBadRequest, "Name, mobile, and email are required")
		return
	}

	sqlStmt := `INSERT INTO users (name, mobile, email) VALUES (?, ?, ?);`
	result, err := db.DB.Exec(sqlStmt, u.Name, u.Mobile, u.Email)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			if strings.Contains(err.Error(), "users.email") {
				sendError(w, http.StatusConflict, "Email already registered")
				return
			}
			if strings.Contains(err.Error(), "users.mobile") {
				sendError(w, http.StatusConflict, "Mobile number already registered")
				return
			}
			sendError(w, http.StatusConflict, "Details already registered")
			return
		}
		sendError(w, http.StatusInternalServerError, "Unable to create user: database error")
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Unable to retrieve created user ID")
		return
	}

	u.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

// Updating user details based on user ID
func UpdateUsers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if u.Name == "" || u.Mobile == "" || u.Email == "" {
		sendError(w, http.StatusBadRequest, "Name, mobile, and email are required")
		return
	}

	sqlStmt := `UPDATE users SET name = ?, mobile = ?, email = ? WHERE id = ?;`
	result, err := db.DB.Exec(sqlStmt, u.Name, u.Mobile, u.Email, id)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			if strings.Contains(err.Error(), "users.email") {
				sendError(w, http.StatusConflict, "Email already registered")
				return
			}
			if strings.Contains(err.Error(), "users.mobile") {
				sendError(w, http.StatusConflict, "Mobile number already registered")
				return
			}
			sendError(w, http.StatusConflict, "Details already registered")
			return
		}
		sendError(w, http.StatusInternalServerError, "Unable to update user: database error")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Unable to verify update operation")
		return
	}
	if rowsAffected == 0 {
		sendError(w, http.StatusNotFound, "User not found")
		return
	}

	u.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

// Searching for existing user based on user ID 
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var u models.User
	err = db.DB.QueryRow(`SELECT id, name, mobile, email FROM users WHERE id=?;`, id).Scan(&u.ID, &u.Name, &u.Mobile, &u.Email)
	if err == sql.ErrNoRows {
		sendError(w, http.StatusNotFound, "User not found")
		return
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Database error: failed to fetch user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

// Searching and deleting user by ID
func DeleteUsers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	sqlStmt := `DELETE FROM users WHERE id=?;`
	result, err := db.DB.Exec(sqlStmt, id)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Unable to delete user: database error")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Unable to verify delete operation")
		return
	}
	if rowsAffected == 0 {
		sendError(w, http.StatusNotFound, "User not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ErrorResponse{Message: "User deleted"})
}