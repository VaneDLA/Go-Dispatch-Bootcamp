package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type PokemonController interface {
	GetPokemonById(w http.ResponseWriter, r *http.Request)
	GetAllPokemon(w http.ResponseWriter, r *http.Request)
	CreatePokemon(w http.ResponseWriter, r *http.Request)
	CatchPokemon(w http.ResponseWriter, r *http.Request)
	FilterPokemons(w http.ResponseWriter, r *http.Request)
}

func Setup(c PokemonController) *mux.Router {
	r := mux.NewRouter()

	v1 := r.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/pokemons", c.GetAllPokemon).Methods(http.MethodGet).Name("GetAllPokemons")

	v1.HandleFunc("/pokemons/{id}", c.GetPokemonById).Methods(http.MethodGet).Name("GetPokemonById")

	v1.HandleFunc("/pokemons", c.CreatePokemon).Methods(http.MethodPost).Name("CreatePokemon")

	v1.HandleFunc("/catch", c.CatchPokemon).Methods(http.MethodPost).Name("CatchPokemon")

	v1.HandleFunc("/filter_pokemons", c.FilterPokemons).
		Queries("type", "{type:odd|even}", "items", "{items:[1-9][0-9]*}", "items_per_worker", "{items_per_worker:[1-9][0-9]*}").
		Methods(http.MethodGet).
		Name("FilterPokemons")

	return r
}
