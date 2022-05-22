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
	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"POST", "GET", "PATCH", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	credential:= handlers.AllowCredentials()
	err := http.ListenAndServe(":8081", handlers.CORS(methods, header, origins,credential)(r))
	if err != nil {
		Logger.ErrorLog.Println(fmt.Errorf("error :%s", err))
	}
}
