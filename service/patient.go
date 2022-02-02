package service

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	errz "github.com/BernardoGR/Go-Dispatch-Bootcamp/errors"
	"github.com/BernardoGR/Go-Dispatch-Bootcamp/model"
)

// PatientServicestruct implements PatientService interface.
type PatientService struct {
	data   model.Patients
}

// New returns a new PatientService instance.
func New(csvPath string) *PatientService {
	data, err := readCsvFile(csvPath)
	if err != nil {
		log.Fatal("Error reading csv: ", err)
	}
	// TODO check for errors reading file

	return &PatientService{
		data: data,
	}
}

// GetAllPatients returns all patients data.
func (ps *PatientService) GetAllPatients() (model.Patients, error) {
	if err := ps.dataValidation(); err != nil {
		return nil, err
	}

	return ps.data, nil
}

// GetPatientByID returns an patient by its ID.
func (ps *PatientService) GetPatientByID(id int) (*model.Patient, error) {
	if err := ps.dataValidation(); err != nil {
		return nil, err
	}

	// find the patient in the data
	for _, p := range ps.data {
			if p.ID == id {
				return &p, nil
			}
	}
	return nil, errz.ErrNotFound
}


// dataValidation is an auxiliary function that checks if the data has been initialized or if it is empty
// returns a matching ServiceError if any of these conditions are met.
func (ps *PatientService) dataValidation() error {
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

func readCsvFile(filePath string) (model.Patients, error) {
	f, err := os.Open(filePath)
	if err != nil {
			fmt.Errorf("Unable to read input file: %v", err)
			return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.FieldsPerRecord = 3

	records, err := csvReader.ReadAll()

	if err != nil {
			fmt.Errorf("Unable to parse file: %v", err)
			return nil, err
	}

	var patient model.Patient
	var patientSlice model.Patients

	for i, r := range records {
		if i == 0 {
			continue
		}
		patient.ID, _ = strconv.Atoi(r[0])
		patient.Name = r[1]
		patient.Age, _ = strconv.Atoi(r[2])
		patientSlice = append(patientSlice, patient)
	}

	return patientSlice, nil
}
