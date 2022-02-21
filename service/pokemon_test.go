package service

import (
	"testing"

	"github.com/carlos-garibay/Go-Dispatch-Bootcamp/model"
	"github.com/google/go-cmp/cmp"
)

var testSliceToPokemon = []struct {
	name     string
	data     []string
	expected *model.Pokemon
	isErr    bool
}{
	{"Invalid number of attributes", []string{"1", "2", "pikachu"}, &model.Pokemon{}, true},
	{"Invalid id attribute", []string{"A", "2", "pikachu", "ImgUrl", "a1,a2,a3"}, &model.Pokemon{}, true},
	{"Invalid number attribute", []string{"1", "B", "pikachu", "ImgUrl", "a1,a2,a3"}, &model.Pokemon{}, true},
	{"Valid data", []string{"1", "2", "pikachu", "ImgUrl", "a1,a2,a3"}, &model.Pokemon{Id: 1, Number: 2, Name: "pikachu", Image: "ImgUrl", Abilities: []string{"a1", "a2", "a3"}}, false},
}

func TestSliceToPokemon(t *testing.T) {
	for _, tt := range testSliceToPokemon {
		got, err := sliceToPokemon(tt.data)
		if tt.isErr {
			if err == nil {
				t.Errorf("Expected an error: %v", tt.name)
			}
		} else {
			if err != nil {
				t.Errorf("Did not expected an error: %v. Error: %v", tt.name, err.Error())
			}
		}
		if got != nil && !cmp.Equal(*got, *tt.expected) {
			t.Errorf("Expected %v but got %v", *tt.expected, *got)
		}
	}
}

var testPokemonToSlice = []struct {
	name     string
	data     *model.Pokemon
	expected []string
}{
	{"Pokemon without abilities", &model.Pokemon{Id: 1, Number: 1, Name: "pikachu", Image: "ImgUrl", Abilities: []string{}}, []string{"1", "1", "pikachu", "ImgUrl", ""}},
	{"Pokemon with one ability", &model.Pokemon{Id: 1, Number: 1, Name: "pikachu", Image: "ImgUrl", Abilities: []string{"agility"}}, []string{"1", "1", "pikachu", "ImgUrl", "agility"}},
	{"Pokemon with two abilities", &model.Pokemon{Id: 1, Number: 1, Name: "pikachu", Image: "ImgUrl", Abilities: []string{"agility", "unnerve"}}, []string{"1", "1", "pikachu", "ImgUrl", "agility,unnerve"}},
}

func TestPokemonToSlice(t *testing.T) {
	for _, tt := range testPokemonToSlice {
		got := pokemonToSlice(tt.data)

		if !cmp.Equal(got, tt.expected) {
			t.Errorf("Expected %v but got %v", tt.expected, got)
		}
	}
}

var testNextId = []struct {
	name     string
	data     []int
	expected int
	isErr    bool
}{
	{"Empty input array", []int{}, 0, true},
	{"Input array with data", []int{1, 4}, 5, false},
}

func TestNextId(t *testing.T) {
	for _, tt := range testNextId {
		got, err := nextId(tt.data)
		if tt.isErr {
			if err == nil {
				t.Errorf("Expected an error: %v", tt.name)
			}
		} else {
			if err != nil {
				t.Errorf("Did not expected an error: %v. Error: %v", tt.name, err.Error())
			}
		}
		if got != tt.expected {
			t.Errorf("Expected %v but got %v", tt.expected, got)
		}
	}
}

var pokemonApiToPokemon = []struct {
	name       string
	inputModel *model.PokemonAPI
	inputSlice []int
	expected   *model.Pokemon
	isErr      bool
}{
	{"Empty input slice", &model.PokemonAPI{}, []int{}, nil, true},
	{
		"Valid Data",
		&model.PokemonAPI{
			ID:   39,
			Name: "jigglypuff",
			Sprites: struct {
				BackDefault      string      `json:"back_default"`
				BackFemale       interface{} `json:"back_female"`
				BackShiny        string      `json:"back_shiny"`
				BackShinyFemale  interface{} `json:"back_shiny_female"`
				FrontDefault     string      `json:"front_default"`
				FrontFemale      interface{} `json:"front_female"`
				FrontShiny       string      `json:"front_shiny"`
				FrontShinyFemale interface{} `json:"front_shiny_female"`
			}{FrontDefault: "ImgUrl"},
			Abilities: []struct {
				Ability struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"ability"`
				IsHidden bool `json:"is_hidden"`
				Slot     int  `json:"slot"`
			}{
				{Ability: struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{Name: "cute-charm"}},
				{Ability: struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{Name: "competitive"}},
			},
		},
		[]int{1, 2, 3},
		&model.Pokemon{
			Id:        4,
			Number:    39,
			Name:      "jigglypuff",
			Image:     "ImgUrl",
			Abilities: []string{"cute-charm", "competitive"},
		},
		false,
	},
}

func TestPokemonAPIToPokemon(t *testing.T) {
	for _, tt := range pokemonApiToPokemon {
		got, err := pokemonAPIToPokemon(tt.inputModel, tt.inputSlice)
		if tt.isErr {
			if err == nil {
				t.Errorf("Expected an error: %v", tt.name)
			}
		} else {
			if err != nil {
				t.Errorf("Did not expected an error: %v. Error: %v", tt.name, err.Error())
			}
		}
		if got != nil && !cmp.Equal(*got, *tt.expected) {
			t.Errorf("Expected %v but got %v", *tt.expected, *got)
		}
	}
}

var testInitDB = []struct {
	name     string
	data     [][]string
	expected PokemonMap
}{
	{"Empty input slice", [][]string{}, PokemonMap{}},
	{
		"Input slice with single row",
		[][]string{{"1", "1", "pikachu", "ImgUrl", "agility,unnerve"}},
		PokemonMap{
			1: model.Pokemon{
				Id:        1,
				Number:    1,
				Name:      "pikachu",
				Image:     "ImgUrl",
				Abilities: []string{"agility", "unnerve"},
			},
		},
	},
	{
		"Input slice with two rows",
		[][]string{
			{"1", "1", "pikachu", "ImgUrl", "agility,unnerve"},
			{"2", "2", "raiku", "ImgUrl", "agility,unnerve"},
		},
		PokemonMap{
			1: model.Pokemon{
				Id:        1,
				Number:    1,
				Name:      "pikachu",
				Image:     "ImgUrl",
				Abilities: []string{"agility", "unnerve"},
			},
			2: model.Pokemon{
				Id:        2,
				Number:    2,
				Name:      "raiku",
				Image:     "ImgUrl",
				Abilities: []string{"agility", "unnerve"},
			},
		},
	},
	{
		"Input slice with invalid row",
		[][]string{
			{"A", "1", "pikachu", "ImgUrl", "agility,unnerve"},
			{"2", "2", "raiku", "ImgUrl", "agility,unnerve"},
		},
		PokemonMap{
			2: model.Pokemon{
				Id:        2,
				Number:    2,
				Name:      "raiku",
				Image:     "ImgUrl",
				Abilities: []string{"agility", "unnerve"},
			},
		},
	},
}

func TestInitDB(t *testing.T) {
	for _, tt := range testInitDB {
		got := initDB(tt.data)

		if !cmp.Equal(got, tt.expected) {
			t.Errorf("Expected %v but got %v", tt.expected, got)
		}
	}
}

var testRemain = []struct {
	name     string
	total    int
	assigned int
	limit    int
	expected int
}{
	{"Remaining more than the limit", 5, 2, 2, 2},
	{"Remaining equal to the limit", 5, 3, 2, 2},
	{"Remaining less than the limit", 5, 4, 2, 1},
}

func TestRemain(t *testing.T) {
	for _, tt := range testRemain {
		got := remain(tt.total, tt.assigned, tt.limit)

		if got != tt.expected {
			t.Errorf("Expected %v but got %v", tt.expected, got)
		}
	}
}
