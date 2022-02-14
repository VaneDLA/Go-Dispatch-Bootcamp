package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type controllers interface {
	GetAllVehicles(w http.ResponseWriter, r *http.Request)
}

func Setup(c controllers) *mux.Router {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/vehicles", c.GetAllVehicles).Methods(http.MethodGet)

	return rtr
}
