package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Akhilbisht798/cloud-text-editor/go-server/internal/database"
	"github.com/Akhilbisht798/cloud-text-editor/go-server/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	database.DbConnect()
	router := routes.NewRouter()

	var port string
	if os.Getenv("APP_ENV") == "production" {
		port = os.Getenv("PORT")
	} else {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Server listening on http://localhost%s\n", addr)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		panic(err)
	}
}
