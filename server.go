package main

import (
	"log"
	"net/http"

	"github.com/rahul/backend-go/db"
	"github.com/rahul/backend-go/middleware"
	"github.com/rahul/backend-go/router"
)

func main()  {
	r := router.SetupRouter()
	cors := middleware.EnableCORS(r)
	db.InitDB()
	log.Println("Server is running on port :8080")
	log.Fatal(http.ListenAndServe(":8080",cors))
}