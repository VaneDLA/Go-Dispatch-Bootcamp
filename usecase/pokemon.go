package usecase

import (
	"log"

	"github.com/carlos-garibay/Go-Dispatch-Bootcamp/model"
)

type PokemonDBService interface {
	GetAllPokemon() (model.Pokemons, error)
	GetPokemonById(id int) (*model.Pokemon, error)
	CreatePokemon(p *model.Pokemon) error
	CatchPokemon(p *model.PokemonAPI) (*model.Pokemon, error)
	FilterPokemons(typ string, items int, itemPerWorker int) (model.Pokemons, error)
}

type pokemonUsecase struct {
	db PokemonDBService
}

func New(db PokemonDBService) pokemonUsecase {
	log.Println("In usecase.pokemon.New")

	return pokemonUsecase{
		db: db,
	}
}

func (pu pokemonUsecase) GetAllPokemon() (model.Pokemons, error) {
	log.Println("In usecase.pokemon.GetAllPokemon")

	pokemons, err := pu.db.GetAllPokemon()
	if err != nil {
		return nil, err
	}

	return pokemons, nil
}

func (pu pokemonUsecase) GetPokemonById(id int) (*model.Pokemon, error) {
	log.Println("In usecase.pokemon.GetPokemonById")

	pokemon, err := pu.db.GetPokemonById(id)
	if err != nil {
		return nil, err
	}

	return pokemon, nil
}

func (pu pokemonUsecase) CreatePokemon(p *model.Pokemon) error {
	log.Println("In usecase.pokemon.CreatePokemon")

	if err := pu.db.CreatePokemon(p); err != nil {
		return err
	}

	return nil
}

func (pu pokemonUsecase) CatchPokemon(p *model.PokemonAPI) (*model.Pokemon, error) {
	log.Println("In usecase.pokemon.CatchPokemon")

	pokemon, err := pu.db.CatchPokemon(p)
	if err != nil {
		return nil, err
	}

	return pokemon, nil
}

func (pu pokemonUsecase) FilterPokemons(typ string, items int, itemsPerWorker int) (model.Pokemons, error) {
	log.Println("In usecase.pokemon.FilterPokemons")

	pokemons, err := pu.db.FilterPokemons(typ, items, itemsPerWorker)
	if err != nil {
		log.Printf("Usecase: error in FilterPokemons: %v", err)
		return nil, err
	}

	return pokemons, nil
}
