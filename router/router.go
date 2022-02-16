package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

// PatientController is the interface that wraps the controller's methods
// GetPatientByID, GetAllPatients, GetPersonByID.
type PatientController interface {
	GetPatientByID(w http.ResponseWriter, r *http.Request)
	GetAllPatients(w http.ResponseWriter, r *http.Request)
	GetRemotePatientByID(w http.ResponseWriter, r *http.Request)
}

// Setup returns a router instance
func Setup(c PatientController) *mux.Router {
	r := mux.NewRouter()

	// versioning api
	v1 := r.PathPrefix("/api/v1").Subrouter()

	// patients endpoints
	p := v1.PathPrefix("/patients").Subrouter()

	p.HandleFunc("/", c.GetAllPatients).
		Methods(http.MethodGet).Name("GetAllPatients")

	p.HandleFunc("/{id}", c.GetPatientByID).
		Methods(http.MethodGet).Name("GetPatientByID")

	p.HandleFunc("/remote/{id}", c.GetRemotePatientByID).
		Methods(http.MethodGet).Name("GetRemotePatientByID")

	return r
}
