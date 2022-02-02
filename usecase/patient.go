package usecase

import (
	"github.com/BernardoGR/Go-Dispatch-Bootcamp/model"
)

// PatientService is the interface that wraps the service's methods
// GetAllPatients, GetPatientByID.
type PatientService interface {
	GetAllPatients() (model.Patients, error)
	GetPatientByID(id int) (*model.Patient, error)
}

// PatientUSecase implements PatientService interface.
type PatientUsecase struct {
	service PatientService
}

// New returns a new PatientUsecase instance.
func New(s PatientService) *PatientUsecase {
	return &PatientUsecase{
		service: s,
	}
}

// GetAllPatients calls the service to returns all patients.
func (pu *PatientUsecase) GetAllPatients() (model.Patients, error) {
	patients, err := pu.service.GetAllPatients()
	if err != nil {
		return nil, err
	}
	return patients, nil
}

// GetPatientByID calls the service to returns an patient by its ID.
func (pu *PatientUsecase) GetPatientByID(id int) (*model.Patient, error) {
	patient, err := pu.service.GetPatientByID(id)
	if err != nil {
		return nil, err
	}
	return patient, nil
}
