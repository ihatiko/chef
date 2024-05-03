package characters

import "example/protoc/characters"

type ICharactersTransport interface {
	characters.CharactersServiceServer
}
