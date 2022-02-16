package model

// Patients type is an alias for a slice of Patient.
type Patients []Patient

// Patient struct represents a single employee.
type Patient struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Age     int     `json:"age"`
}
