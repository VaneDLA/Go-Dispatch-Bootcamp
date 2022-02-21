package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	errz "github.com/carlos-garibay/Go-Dispatch-Bootcamp/errors"
	"github.com/carlos-garibay/Go-Dispatch-Bootcamp/model"
	"github.com/gorilla/mux"
)

type PokemonUsecase interface {
	GetAllPokemon() (model.Pokemons, error)
	GetPokemonById(id int) (*model.Pokemon, error)
	CreatePokemon(p *model.Pokemon) error
	CatchPokemon(p *model.PokemonAPI) (*model.Pokemon, error)
	FilterPokemons(typ string, items int, items_per_worker int) (model.Pokemons, error)
}

type pokemonController struct {
	usecase PokemonUsecase
}

func New(uc PokemonUsecase) pokemonController {
	log.Println("In controller.pokemon.New")
	return pokemonController{
		usecase: uc,
	}
}

func (pc pokemonController) GetPokemonById(w http.ResponseWriter, r *http.Request) {
	log.Println("In controller.pokemon.GetPokemonById")

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid id: %v", err.Error())
		log.Printf("Invalid id: %v\n", err.Error())
		return
	}

	pokemon, err := pc.usecase.GetPokemonById(id)
	if err != nil {
		switch {
		case errors.Is(err, errz.ErrEmptyData), errors.Is(err, errz.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "%v", err.Error())
			log.Println(err.Error())
			return
		case errors.Is(err, errz.ErrDataNotInitialized):
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err.Error())
			log.Println(err.Error())
			return
		}
	}

	if pokemon == nil {
		log.Println("pokemon not found")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "pokemon not found")
		return
	}

	jsonPokemon, err := json.Marshal(pokemon)
	if err != nil {
		log.Println("error marshalling pokemon")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error marshalling pokemon")
		return
	}

	log.Printf("pokemon found: %+v\n", pokemon)
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonPokemon)
	w.WriteHeader(http.StatusOK)
}

func (pc pokemonController) GetAllPokemon(w http.ResponseWriter, r *http.Request) {
	log.Println("In controller.pokemon.GetAllPokemon")

	pokemons, err := pc.usecase.GetAllPokemon()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error getting pokemons: %v", err.Error())
		log.Printf("error getting pokemons: %v\n", err.Error())
		return
	}

	if len(pokemons) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "no pokemons found")
		log.Println("no pokemons found")
		return
	}

	jsonPokemons, err := json.Marshal(pokemons)
	if err != nil {
		log.Printf("error marshalling pokemons: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error marshalling pokemons: %v\n", err.Error())
		return
	}
	log.Printf("pokemons found: %+v\n", pokemons)

	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonPokemons)
	w.WriteHeader(http.StatusOK)
}

func (pc pokemonController) CreatePokemon(w http.ResponseWriter, r *http.Request) {
	log.Println("In controller.pokemon.CreatePokemon")

	pokemon := model.Pokemon{}
	err := json.NewDecoder(r.Body).Decode(&pokemon)
	if err != nil {
		log.Printf("r.Body: %+v\n", r.Body)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request body")
		log.Printf("decoding request body: %v\n", err)
		return
	}

	err = pc.usecase.CreatePokemon(&pokemon)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error creating pokemon: %v", err.Error())
		log.Printf("error creating pokemon: %v", err.Error())
		return
	}

	log.Printf("pokemon created: %+v\n", pokemon)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "pokemon created")
}

func (pc pokemonController) CatchPokemon(w http.ResponseWriter, r *http.Request) {
	log.Println("In controller.pokemon.CatchPokemon")

	var pokemonNumber model.Pokemon
	err := json.NewDecoder(r.Body).Decode(&pokemonNumber)
	if err != nil {
		log.Printf("r.Body: %+v\n", r.Body)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request body")
		log.Printf("decoding request body: %v\n", err)
		return
	}

	log.Printf("%+v\n", pokemonNumber)

	resp, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d", pokemonNumber.Number))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not catch new pokemon: %v", err.Error())
		log.Printf("Could not catch new pokemon: %v", err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not catch new pokemon: %v", err.Error())
		log.Printf("Could not catch new pokemon: %v", err.Error())
		return
	}

	var pokemonApi model.PokemonAPI
	if err = json.Unmarshal(body, &pokemonApi); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not catch new pokemon: %v", err.Error())
		log.Printf("Could not catch new pokemon: %v", err.Error())
		return
	}

	pokemon, err := pc.usecase.CatchPokemon(&pokemonApi)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error catching pokemon: %v", err.Error())
		log.Printf("Error catching pokemon: %v", err.Error())
		return
	}

	jsonPokemon, err := json.Marshal(*pokemon)
	if err != nil {
		log.Printf("error marshalling pokemon: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error marshalling pokemon: %v\n", err.Error())
		return
	}

	log.Printf("pokemon catched: %v - %v\n", pokemon.Id, pokemon.Name)
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonPokemon)
	w.WriteHeader(http.StatusCreated)
}

func (pc pokemonController) FilterPokemons(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	items, err := strconv.Atoi(params["items"][0])
	if err != nil {
		log.Printf("r.URL: %+v\n", r.URL)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid value for items param: %v", err)
		log.Printf("invalid value for items param: %v\n", err)
		return
	}

	itemsPerWorker, err := strconv.Atoi(params["items_per_worker"][0])
	if err != nil {
		log.Printf("r.URL: %+v\n", r.URL)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid value for items_per_worker param: %v", err)
		log.Printf("invalid value for items_per_worker param: %v\n", err)
		return
	}

	if items < itemsPerWorker {
		log.Printf("Items can't be less than items_per_worker. items: %v items_per_worker: %v\n", items, itemsPerWorker)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Items can't be less than items_per_worker. items: %v items_per_worker: %v\n", items, itemsPerWorker)
		return
	}

	typ := params["type"][0]

	pokemons, err := pc.usecase.FilterPokemons(typ, items, itemsPerWorker)
	if err != nil {
		log.Printf("Error filtering pokemons: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error filtering pokemons: %v\n", err)
		return
	}

	jsonPokemons, err := json.Marshal(pokemons)
	if err != nil {
		log.Printf("error marshalling pokemons: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error marshalling pokemons: %v\n", err.Error())
		return
	}
	log.Printf("pokemons found: %+v\n", pokemons)

	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonPokemons)
	w.WriteHeader(http.StatusOK)
}
