package model

// From the swapi documentation: https://github.com/peterhellberg/swapi/blob/master/person.go
// A Person is an individual person or character within the Star Wars universe.
type Person struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Height       string   `json:"height"`
	Mass         string   `json:"mass"`
	HairColor    string   `json:"hair_color"`
	SkinColor    string   `json:"skin_color"`
	EyeColor     string   `json:"eye_color"`
	BirthYear    string   `json:"birth_year"`
	Gender       string   `json:"gender"`
	Homeworld    string   `json:"homeworld"`
	FilmURLs     []string `json:"films"`
	SpeciesURLs  []string `json:"species"`
	VehicleURLs  []string `json:"vehicles"`
	StarshipURLs []string `json:"starships"`
	Created      string   `json:"created"`
	Edited       string   `json:"edited"`
	URL          string   `json:"url"`
}