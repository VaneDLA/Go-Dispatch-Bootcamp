package errors

type PokemonError string

var (
	ErrNotFound           = PokemonError("pokemon not found")
	ErrEmptyData          = PokemonError("data is empty")
	ErrDataNotInitialized = PokemonError("data not initialized")
	ErrPokemonExists      = PokemonError("pokemon already exists")
)

func (e PokemonError) Error() string {
	return string(e)
}
