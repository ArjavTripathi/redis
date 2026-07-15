package lexer

import (
	"redis/src/commands"
	"strings"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
	SKIP    = '\r'
	NULL    = '_'
	BOOLEAN = '#'
	DOUBLE  = ','
)

func handleBulkStrings(stream []byte, index *int) commands.Token {
	i := 0
	sb := strings.Builder{}
	for stream[*index] != SKIP {
		i = 10*i + int(stream[*index]-'0')
		*index++
	}
	*index += 2
	for j := 0; j < i; j++ {
		sb.WriteByte(stream[*index])
		*index++
	}
	*index += 2
	token := commands.Token{
		Type: BULK,
		Str:  sb.String(),
	}
	return token
}

func handleInteger(stream []byte, index *int) commands.Token {
	toNegate := false
	val := 0
	sym := '+'
	for stream[*index] != SKIP {
		if stream[*index] == '-' {
			toNegate = true
			sym = '-'
			*index++
			continue
		}
		val = val*10 + int(stream[*index]-'0')
		*index++
	}

	if toNegate {
		val *= -1
	}

	inttoken := commands.Token{
		Type:   INTEGER,
		Num:    val,
		Symbol: sym,
	}
	return inttoken

}

func handleNull() commands.Token {
	return commands.Token{
		Type: NULL,
	}
}

func handleBoolean(stream []byte, index *int) commands.Token {
	var sb strings.Builder
	var b bool
	if stream[*index] == 't' {
		sb.WriteString("true")
		b = true
	} else if stream[*index] == 'f' {
		sb.WriteString("false")
		b = false
	}
	*index++
	return commands.Token{
		Type: BOOLEAN,
		Bool: b,
	}
}

func manageArray(stream []byte) commands.Token {
	noOfCommands := int(stream[1] - '0')
	if noOfCommands == 0 {
		return commands.ErrorToken("0 commands found in Array")
	}

	token := commands.Token{
		Type:  ARRAY,
		Num:   noOfCommands,
		Array: make([]commands.Token, 0, noOfCommands),
	}

	commandCounter := 0
	i := 2

	for commandCounter < noOfCommands {
		if i >= len(stream) {
			break
		}
		if stream[i] == SKIP {
			i += 2
			continue
		}
		if stream[i] == BULK {
			commandCounter++
			i++
			bulktoken := handleBulkStrings(stream, &i)
			token.Array = append(token.Array, bulktoken)
		} else if stream[i] == INTEGER {
			commandCounter++
			i++
			inttoken := handleInteger(stream, &i)
			token.Array = append(token.Array, inttoken)
		} else if stream[i] == NULL {
			commandCounter++
			i++
			nullToken := handleNull()
			token.Array = append(token.Array, nullToken)
		} else if stream[i] == BOOLEAN {
			i++
			commandCounter++
			boolToken := handleBoolean(stream, &i)
			token.Array = append(token.Array, boolToken)
		} else if stream[i] == DOUBLE {
			//i++
			//var symbol *byte
			//if stream[i] == '-' {
			//	symbol = &stream[i]
			//	i++
			//}

		}
	}

	return token
}

func ReadManager(srv *commands.Server, stream []byte) commands.Token {
	dataType := stream[0]
	var token commands.Token
	if dataType == BULK {
		token = handleBulkStrings(stream, new(1))
	} else if dataType == INTEGER {
		token = handleInteger(stream, new(1))
	} else if dataType == NULL {
		token = handleNull()
	} else if dataType == ARRAY {
		token = manageArray(stream)
	}

	handler, err := srv.Handler([]commands.Token{token})
	if err != nil {
		return commands.ErrorToken(err.Error())
	}

	return commands.CreateStringToken(handler)
}
