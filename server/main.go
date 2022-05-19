package main

import (
	"example/server/database"
	Logger "example/server/logger"
	"example/server/route"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
	database.DbMigrator()
	r := route.Router()
	methods := handlers.AllowedMethods([]string{"POST", "GET", "PATCH", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	err := http.ListenAndServe(":8081", handlers.CORS(methods, origins)(r))
	if err != nil {
		Logger.ErrorLog.Println(fmt.Errorf("error :%s", err))
	}
	Logger.CommonLog.Println("Server is Running")
	Logger.CommonLog.Println("checking")
}
