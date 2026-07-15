package commands

import (
	"errors"
	"fmt"
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
	handler := map[string]func(*Store, []Token) (string, error){
		"SET": srv.set,
	}
	return handler[instruction](newStore(), args)
}

func ping(s *Store, args []Token) (string, error) {
	return "+PONG\r\n", nil
}
func (srv *Server) set(s *Store, args []Token) (string, error) {

	if len(args) != 3 {
		return "", errors.New("Wrong number of arguments for the set operation")
	}

	return "what", nil
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
