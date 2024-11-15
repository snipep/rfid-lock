package router

import (
	"github.com/gorilla/mux"
	"github.com/snipep/iot/internal/handler"
)

func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", handler.HelloWorld).Methods("GET")
	router.HandleFunc("/data", handler.GetData).Methods("GET")
	router.HandleFunc("/log", handler.GetUserLog).Methods("GET")

	return router
}