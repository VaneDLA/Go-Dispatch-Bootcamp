package usecase

import (
	"github.com/edgarduran/Go-Dispatch-Bootcamp/models"
	"github.com/edgarduran/Go-Dispatch-Bootcamp/service"
)

func GetAllVehicles() ([]models.Vehicle, error) {
	vehicles, err := service.GetAllVehicles()
	if err != nil {
		return nil, err
	}

	return vehicles, nil
}

func GetVehicleById(id int) ([]models.Vehicle, error) {
	vehicle, err := service.GetVehicleById(id)
	if err != nil {
		return nil, err
	}

	return vehicle, nil
}
