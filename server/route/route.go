package route

import (
	"example/server/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/user", controller.Signup).Methods("POST")
	r.HandleFunc("/user", controller.Athutenticate(controller.Read)).Methods("GET")
	r.HandleFunc("/user/{id}", controller.Athutenticate(controller.UpdateUser)).Methods("PATCH")
	r.HandleFunc("/user/{id}", controller.DeleteUser).Methods("DELETE")
	r.HandleFunc("/login", controller.Login).Methods("POST")
	r.HandleFunc("/logout", controller.Logout).Methods("GET")
	return r
}
