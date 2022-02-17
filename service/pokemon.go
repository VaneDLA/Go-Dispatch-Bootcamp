package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/carlos-garibay/Go-Dispatch-Bootcamp/errors"
	"github.com/carlos-garibay/Go-Dispatch-Bootcamp/model"
	"github.com/carlos-garibay/Go-Dispatch-Bootcamp/utils"
)

// PokemonMap is an alias for a map of pokemons.
type PokemonMap map[int]model.Pokemon

type pokemonDBService struct {
	data PokemonMap
}

type worker struct {
	id           int
	recProcessed int
	itemsAdded   int
	maxItems     int
	typ          string
}

var pokemonOrder []int = []int{}

func (w *worker) Id() int {
	return w.id
}

func (w *worker) Execute(in <-chan []string, out chan<- *model.Pokemon) {
	log.Printf("Service: started processing for worker %v\n", w.id)

	for data := range in {
		w.recProcessed++
		pokemon, err := sliceToPokemon(data)
		if err != nil {
			log.Printf("Service: error parsing data to pokemon: %v\n", err)
			continue
		}

		if (pokemon.Id%2 == 0 && w.typ == "even") || (pokemon.Id%2 == 1 && w.typ == "odd") {
			w.itemsAdded++
			out <- pokemon

			if w.itemsAdded == w.maxItems {
				log.Printf("Service: worker %v finished - reached max lines.\n", w.id)
				return
			}
		}
	}
	log.Printf("Service: worker %v finished - input channel closed.\n", w.id)
}

func New(pm PokemonMap) pokemonDBService {
	log.Println("In service.pokemon.New")
	if pm == nil {
		//Init db with csv file
		data, err := utils.ReadLines("pokemons.csv")
		if err != nil && data == nil {
			log.Fatal(err.Error())
		}
		pm = initDB(data)
	}

	return pokemonDBService{
		data: pm,
	}
}

func (ps pokemonDBService) validateDB() error {
	log.Println("In service.pokemon.validateDB")

	if ps.data == nil {
		return errors.ErrDataNotInitialized
	}

	if len(ps.data) == 0 {
		return errors.ErrEmptyData
	}

	return nil
}

func (ps pokemonDBService) GetAllPokemon() (model.Pokemons, error) {
	log.Println("In service.pokemon.GetAllPokemon")

	if err := ps.validateDB(); err != nil {
		return nil, err
	}

	pokemonArray := make(model.Pokemons, 0, len(ps.data))
	for _, id := range pokemonOrder {
		pokemonArray = append(pokemonArray, ps.data[id])
	}

	return pokemonArray, nil
}

func (ps pokemonDBService) GetPokemonById(id int) (*model.Pokemon, error) {
	log.Println("In service.pokemon.GetPokemonById")

	if err := ps.validateDB(); err != nil {
		return nil, err
	}

	pokemon, ok := ps.data[id]
	if !ok {
		return nil, errors.ErrNotFound
	}
	return &pokemon, nil
}

func (ps pokemonDBService) CreatePokemon(p *model.Pokemon) error {
	log.Println("In service.pokemon.CreatePokemon")

	if err := ps.validateDB(); err != nil {
		return err
	}

	if _, ok := ps.data[p.Id]; ok {
		return errors.ErrPokemonExists
	}

	pokemonData := [][]string{pokemonToSlice(p)}
	if err := utils.WriteLines("pokemons.csv", pokemonData); err != nil {
		return err
	}

	pokemonOrder = append(pokemonOrder, p.Id)
	ps.data[p.Id] = *p

	return nil
}

func (ps pokemonDBService) CatchPokemon(p *model.PokemonAPI) (*model.Pokemon, error) {
	log.Println("In service.pokemon.CatchPokemon")

	if err := ps.validateDB(); err != nil {
		return nil, err
	}

	pokemon, err := pokemonAPIToPokemon(p)
	if err != nil {
		return nil, err
	}

	pokemonData := [][]string{pokemonToSlice(pokemon)}
	if err := utils.WriteLines("pokemons.csv", pokemonData); err != nil {
		return nil, err
	}

	pokemonOrder = append(pokemonOrder, pokemon.Id)
	ps.data[pokemon.Id] = *pokemon

	return pokemon, nil
}

func (ps pokemonDBService) FilterPokemons(typ string, items int, itemsPerWorker int) (model.Pokemons, error) {
	log.Println("In service.pokemon.FilterPokemons")

	//To make sure we have enough workers to process the requested items
	numWorkers := int(math.Ceil(float64(items) / float64(itemsPerWorker)))
	wg := &sync.WaitGroup{}

	resultChan := make(chan *model.Pokemon, items)
	dataChan := make(chan []string)

	wp := utils.NewWorkerPool(wg, numWorkers, dataChan, resultChan)
	addWorkers(wp, typ, items, itemsPerWorker)
	wp.Run()

	file, err := utils.OpenFile("pokemons.csv")
	if err != nil {
		log.Println("Service: error opening file pokemon.csv")
		return nil, err
	}
	defer utils.CloseFile(file)

	csvReader := csv.NewReader(file)

	// read the csv and write it to dataChan
	go func() {
		log.Println("Started routine to read csv file.")
		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				line, _ := csvReader.FieldPos(0)
				log.Printf("Error reading line %v from csv.\n", line)
				continue
			}
			if line, _ := csvReader.FieldPos(0); line == 1 {
				continue
			}
			dataChan <- record
		}
		close(dataChan)
	}()

	// wait for worker group to finish and close out channel
	go func() {
		wg.Wait()
		wp.Close()
	}()

	result := model.Pokemons{}
	for pokemon := range resultChan {
		result = append(result, *pokemon)
	}

	return result, nil
}

func addWorkers(wp utils.IworkerPool, typ string, items, itemsPerWorker int) {
	i := 0
	id := 0
	for i < items {
		id++
		itemsAssigned := remain(items, i, itemsPerWorker)
		wp.AddWorker(&worker{
			id:           id,
			recProcessed: 0,
			itemsAdded:   0,
			maxItems:     itemsAssigned,
			typ:          typ,
		})
		log.Printf("Service: Added worker %v to worker pool.\n", id)
		i += itemsAssigned
	}
}

func remain(total, assigned, limit int) int {
	if total-assigned >= limit {
		return limit
	}
	return total - assigned
}

func initDB(data [][]string) PokemonMap {
	result := PokemonMap{}

	for _, pokemonData := range data {
		pokemon, err := sliceToPokemon(pokemonData)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		pokemonOrder = append(pokemonOrder, pokemon.Id)
		result[pokemon.Id] = *pokemon
	}
	return result
}

func NextId() (int, error) {
	if len(pokemonOrder) == 0 {
		return 0, errors.ErrDataNotInitialized
	}
	return pokemonOrder[len(pokemonOrder)-1] + 1, nil
}

func pokemonAPIToPokemon(p *model.PokemonAPI) (*model.Pokemon, error) {
	nextId, err := NextId()
	if err != nil {
		log.Println("Uninitialized DB - NextId")
		return nil, errors.ErrDataNotInitialized
	}

	abilities := []string{}
	for _, ability := range p.Abilities {
		abilities = append(abilities, ability.Ability.Name)
	}

	return &model.Pokemon{
		Id:        nextId,
		Number:    p.ID,
		Name:      p.Name,
		Image:     p.Sprites.FrontDefault,
		Abilities: abilities,
	}, nil
}

func pokemonToSlice(p *model.Pokemon) []string {
	pokemonData := make([]string, model.NumOfAttributes)
	for i := 0; i < model.NumOfAttributes; i++ {
		switch i {
		case model.Id:
			pokemonData[model.Id] = fmt.Sprint(p.Id)
		case model.Number:
			pokemonData[model.Number] = fmt.Sprint(p.Number)
		case model.Name:
			pokemonData[model.Name] = p.Name
		case model.Image:
			pokemonData[model.Image] = p.Image
		case model.Abilities:
			pokemonData[model.Abilities] = strings.Join(p.Abilities, ",")
		}
	}
	return pokemonData
}

func sliceToPokemon(data []string) (*model.Pokemon, error) {
	if len(data) < model.NumOfAttributes {
		return nil, errors.PokemonError(fmt.Sprintf("Invalid number of attributes to create a pokemon: Provided: %v. Required: %v", len(data), model.NumOfAttributes))
	}

	result := &model.Pokemon{}
	for i := 0; i < model.NumOfAttributes; i++ {
		switch i {
		case model.Id:
			pokemonId, err := strconv.Atoi(data[model.Id])
			if err != nil {
				return nil, errors.PokemonError(fmt.Sprintf("Error parsing data for pokemon id: Value '%v' can't be parsed to int", data[model.Id]))
			}
			result.Id = pokemonId
		case model.Number:
			pokemonNumber, err := strconv.Atoi(data[model.Number])
			if err != nil {
				return nil, errors.PokemonError(fmt.Sprintf("Error parsing data for pokemon number: Value '%v' can't be parsed to int", data[model.Number]))
			}
			result.Number = pokemonNumber
		case model.Name:
			result.Name = data[model.Name]
		case model.Image:
			result.Image = data[model.Image]
		case model.Abilities:
			abilities := strings.Split(data[model.Abilities], ",")
			result.Abilities = abilities
		}
	}
	return result, nil
}
