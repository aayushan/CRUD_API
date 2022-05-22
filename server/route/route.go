package route

import (
	"example/server/controller"
	"example/server/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/user", controller.Signup).Methods("POST")
	r.HandleFunc("/user", middleware.Athutenticate(controller.Read)).Methods("GET")
	r.HandleFunc("/update/{id}", middleware.Athutenticate(controller.UpdateUser)).Methods("PATCH")
	r.HandleFunc("/delete/{id}", controller.DeleteUser).Methods("DELETE")
	r.HandleFunc("/login", middleware.Athutenticate(controller.Login)).Methods("POST")
	r.HandleFunc("/logout", middleware.Athutenticate(controller.Logout)).Methods("GET")
	return r
}
