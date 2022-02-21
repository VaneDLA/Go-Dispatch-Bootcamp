package model

const NumOfAttributes = 5

const (
	Id = iota
	Number
	Name
	Image
	Abilities
)

type Pokemons []Pokemon

type Pokemon struct {
	Id        int      `json:"id"`
	Number    int      `json:"number"`
	Name      string   `json:"name"`
	Image     string   `json:"image_url"`
	Abilities []string `json:"abilities"`
}
