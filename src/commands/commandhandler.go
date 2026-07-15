package commands

import (
	"fmt"
	"sync"
)

func Handler(instruction string, args []Token) string {
	handler := map[string]func(*Store, []Token) string{
		"PING": ping,
		"SET":  set,
	}
	return handler[instruction](newStore(), args)
}

type Data struct {
	data string
}

type Store struct {
	mu    sync.RWMutex
	cache map[string]Data
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
		Type: '-',
		Str:  errorMessage,
	}
}

func ping(s *Store, args []Token) string {
	return "+PONG\r\n"
}
func set(s *Store, args []Token) string {

	return "what"
}
