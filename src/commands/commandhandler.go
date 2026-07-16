package commands

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
	NULL    = '_'
	BOOLEAN = '#'
	DOUBLE  = ','
)

type Data struct {
	data string
}

type Store struct {
	mu    sync.RWMutex
	cache map[string]Data
}

type Server struct {
	store *Store
}

func NewServer() *Server {
	return &Server{
		store: &Store{
			cache: make(map[string]Data),
			mu:    sync.RWMutex{},
		},
	}
}

func (srv *Server) Handler(args []Token) (string, error) {

	numTokens := len(args)
	if numTokens == 0 {
		return "", errors.New("Missing tokens in handler")
	}
	instruction := args[0].Str
	handler := map[string]func([]Token) (string, error){
		"SET": srv.set,
		"GET": srv.get,
	}
	return handler[instruction](args)
}

func ping(args []Token) (string, error) {
	return "+PONG\r\n", nil
}

func (srv *Server) get(args []Token) (string, error) {
	if len(args) != 2 {
		return "", errors.New("wrong number of arguments for the get operation")
	}

	if args[1].Type != BULK {
		return "", errors.New("wrong type for the key value")
	}

	srv.store.mu.RLock()
	defer srv.store.mu.RUnlock()

	value, ok := srv.store.cache[args[1].Str]
	if !ok {
		return "", errors.New("key not found")
	}
	return value.data, nil
}

func (srv *Server) set(args []Token) (string, error) {
	if len(args) != 3 {
		return "", errors.New("wrong number of arguments for the set operation")
	}

	if args[1].Type != BULK {
		return "", errors.New("wrong type for the key value")
	}

	srv.store.mu.RLock()
	defer srv.store.mu.RUnlock()

	srv.store.cache[args[1].Str] = Data{data: args[2].getValue()}

	return "success", nil
}

func newStore() *Store {
	return &Store{
		cache: make(map[string]Data),
	}
}

type Token struct {
	Type   int32
	Num    int
	Bool   bool
	Str    string
	Symbol int32
	Array  []Token
}

func (t Token) getValue() string {
	switch t.Type {
	case STRING:
		return t.Str
	case INTEGER:
		return fmt.Sprintf("%d", t.Num)
	case BULK:
		return t.Str
	case BOOLEAN:
		return strconv.FormatBool(t.Bool)
	}
	return ""
}

func (t Token) Equals(other Token) bool {
	if t.Type != other.Type ||
		t.Num != other.Num ||
		t.Bool != other.Bool ||
		t.Str != other.Str ||
		t.Symbol != other.Symbol {
		return false
	}

	if len(t.Array) != len(other.Array) {
		return false
	}

	for i := range t.Array {
		if !t.Array[i].Equals(other.Array[i]) {
			return false
		}
	}

	return true
}

func (t Token) String() string {
	rootString := fmt.Sprintf(
		"Type: %d\n Num: %d\n Bool: %t\n Str: %s\n Symbol: %d\n",
		t.Type, t.Num, t.Bool, t.Str, t.Symbol)

	if len(t.Array) > 0 {
		for i, token := range t.Array {
			formattedString := fmt.Sprintf(
				"Type: %d\n Num: %d\n Bool: %t\n Str: %s\n Symbol: %d\n",
				token.Type, token.Num, token.Bool, token.Str, token.Symbol)

			rootString += fmt.Sprintf("Subarray %d\n", i) + formattedString
		}

	}

	return rootString
}

func ErrorToken(errorMessage string) Token {
	return Token{
		Type: ERROR,
		Str:  errorMessage,
	}
}

func (t Token) array(s *Store, commands []Token) (string, error) {
	return "", nil
}

func CreateStringToken(str string) Token {
	return Token{
		Type: STRING,
		Str:  str,
	}
}
