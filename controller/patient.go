package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	errz "github.com/BernardoGR/Go-Dispatch-Bootcamp/errors"
	"github.com/BernardoGR/Go-Dispatch-Bootcamp/model"

	"github.com/gorilla/mux"
)

// usecase is the interface that wraps the usecase's methods
// GetPatientByID, GetAllPatients.
type usecase interface {
	GetPatientByID(id int) (model.Patient, error)
	GetAllPatients() (model.Patients, error)
}

// patientController implements PatientUsecase interface.
type patientController struct {
	usecase usecase
}

// New returns a new PatientController instance.
func New(uc usecase) patientController {
	return patientController{
		usecase: uc,
	}
}

// GetAllPatients calls the usecase to return all patients.
func (pc patientController) GetAllPatients(w http.ResponseWriter, r *http.Request) {
	// get all patients from the usecase
	patients, err := pc.usecase.GetAllPatients()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error getting patients")
	}

	// special handling if patients is empty
	if len(patients) == 0 {
		log.Println("no patients found")
		w.WriteHeader(http.StatusNotFound)

		fmt.Fprintln(w, "no patients found")
		return
	}

	jsonData, err := json.Marshal(patients)
	if err != nil {
		log.Println("error marshalling patients")
		w.WriteHeader(http.StatusInternalServerError)

		fmt.Fprintf(w, "error marshalling patients: %v\n", err)
	}

	// this is fine
	log.Printf("patients found: %+v\n", patients)

	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonData)
	w.WriteHeader(http.StatusOK)
}

// GetPatientByID returns an patient by its ID.
func (pc patientController) GetPatientByID(w http.ResponseWriter, r *http.Request) {
	// extract the path parameters
	vars := mux.Vars(r)

	// convert the id param into an int
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid id: %v", err)
	}

	// get the patient from the usecase
	patient, err := pc.usecase.GetPatientByID(id)
	if err != nil {
		switch {
		case errors.Is(err, errz.ErrNotFound), errors.Is(err, errz.ErrEmptyData):
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "patient not found")

		case errors.Is(err, errz.ErrDataNotInitialized):
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "data not initialized")
		}
	}

	if (patient == model.Patient{}) {
		log.Println("no patient found")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "no patient found")

		return
	}

	jsonData, err := json.Marshal(patient)
	if err != nil {
		log.Println("error marshalling patients")
		w.WriteHeader(http.StatusInternalServerError)

		fmt.Fprintf(w, "error marshalling patients")
	}

	// this is fine
	log.Printf("patient found: %+v", patient)

	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonData)
	w.WriteHeader(http.StatusOK)
}
