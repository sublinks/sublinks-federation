package routes

import "github.com/gorilla/mux"

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	SetupUserRoutes(r)
	SetupPostRoutes(r)
	SetupApubRoutes(r)
	return r
}
