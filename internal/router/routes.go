package router

import (
	"github.com/gorilla/mux"
	"github.com/snipep/iot/internal/handler"
)

func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()


	router.HandleFunc("/user/{id}", handler.GetUserData).Methods("GET")
	router.HandleFunc("/log/{id}", handler.GetUserLog).Methods("GET")
	router.HandleFunc("/register", handler.RegisterUser).Methods("POST")
	router.HandleFunc("/edit/user", handler.EditUser).Methods("PUT")
	router.HandleFunc("/change/access", handler.ChangeAccess).Methods("PUT")
	router.HandleFunc("/delete/{id}", handler.DeleteUser).Methods("DELETE")
	router.HandleFunc("/get/users", handler.GetAllUsers).Methods("GET")

	return router
}