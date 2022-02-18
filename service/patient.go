package service

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	errz "github.com/BernardoGR/Go-Dispatch-Bootcamp/errors"
	"github.com/BernardoGR/Go-Dispatch-Bootcamp/model"
	"github.com/BernardoGR/Go-Dispatch-Bootcamp/repository"
)

const csvPath = "./resources/patients.csv"

// PatientServicestruct implements PatientService interface.
type PatientService struct {
	data     model.Patients
	nextID   int
}

// New returns a new PatientService instance.
func New() PatientService {
	raw_data, err := repository.ReadCsvFile(csvPath)
	
	if err != nil {
		log.Println("Error reading csv: ", err)
	}

	var data = parseCSVPatients(raw_data)


	return PatientService {
		data: data,
		nextID: len(data) + 1,
	}
}

// GetAllPatients returns all patients data.
func (ps PatientService) GetAllPatients() (model.Patients, error) {
	if err := ps.dataValidation(); err != nil {
		return nil, err
	}

	return ps.data, nil
}

// GetPatientByID returns a patient by its ID.
func (ps PatientService) GetPatientByID(id int) (model.Patient, error) {
	if err := ps.dataValidation(); err != nil {
		return model.Patient{}, err
	}

	// find the patient in the data
	for _, p := range ps.data {
			if p.ID == id {
				return p, nil
			}
	}
	return model.Patient{}, errz.ErrNotFound
}

// create new patient.
func (ps *PatientService) CreatePatientFromRemote(resp *http.Response) (model.Patient, error) {
	// parse remote patient
	var person model.Person
	json.NewDecoder(resp.Body).Decode(&person)
	person.ID = ps.nextID

	// convert remote patient to patient
	patient := personToPatient(person)

	// create patient
	err := createPatient(patient)

	if err != nil {
		return model.Patient{}, errz.ErrCreationFailed
	}

	// update next id for further patient creations
	ps.nextID += 1
	ps.data = append(ps.data, patient)
	return patient, nil
}

func parseCSVPatients(raw_data [][]string) model.Patients {
	var patient model.Patient
	var patientSlice model.Patients

	for i, r := range raw_data {
		if i == 0 {
			continue
		}
		patient.ID, _ = strconv.Atoi(r[0])
		patient.Name = r[1]
		patient.Age, _ = strconv.Atoi(r[2])
		patientSlice = append(patientSlice, patient)
	}

	return patientSlice
}

func createPatient(patient model.Patient) (error) {
  id := strconv.Itoa(patient.ID)
  age := strconv.Itoa(patient.Age)
	return repository.AddLine(csvPath, []string{id, patient.Name, age})
}

// parse patient from remote patient format.
func personToPatient(person model.Person) (model.Patient) {
	var patient model.Patient
	patient.ID = person.ID
	patient.Name = person.Name
	patient.Age = ageFromBirth(person.BirthYear)
	return patient
}

func ageFromBirth(birth_str string) (int) {
	age_num := 0
	if birth_str == "unknown" {
		return age_num
	}
	to_index := len(birth_str) - 3
	if to_index >= 0 {
		decimal_index := strings.Index(birth_str, ".")
		if decimal_index > -1 {
			to_index = decimal_index
		}
		age := birth_str[0:to_index]
		age_num, _ = strconv.Atoi(age)
	}
	return age_num
}


// dataValidation is an auxiliary function that checks if the data has been initialized or if it is empty
// returns a matching ServiceError if any of these conditions are met.
func (ps PatientService) dataValidation() error {
	// special handling if data is nil
	if ps.data == nil {
		return errz.ErrDataNotInitialized
	}

	// special handling if data is empty
	if len(ps.data) == 0 {
		return errz.ErrEmptyData
	}

	return nil
}
