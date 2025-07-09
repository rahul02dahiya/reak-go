package db

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}

	CreateUserTable := `CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	mobile TEXT NOT NULL UNIQUE, 
	email TEXT NOT NULL UNIQUE,
	CONSTRAINTS check_mobile CHECK (mobile GLOB '[6-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]') 
);`

	_, err = DB.Exec(CreateUserTable)

	if err != nil {
		log.Fatalf("Error while creating table %q : %s\n", err, CreateUserTable)
	}
}
