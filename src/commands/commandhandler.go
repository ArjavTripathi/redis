package commands

func Handler(instruction string, args []Token) string {
	handler := map[string]func([]Token) string{
		"PING": ping,
		"SET":  set,
	}
	return handler[instruction](args)
}

type Token struct {
	Type   int32
	Num    int
	Bool   bool
	Str    string
	Symbol byte
	Array  []Token
}

func ping(args []Token) string {
	return "+PONG\r\n"
}
func set(args []Token) string {
	return "what"
}
