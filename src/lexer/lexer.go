package lexer

import (
	"fmt"
	"redis/src/commands"
	"strconv"
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

func manageArray(stream []byte) string {
	noOfCommands := int(stream[1] - '0')
	if noOfCommands == 0 {
		return ""
	}

	token := commands.Token{
		Type:  ARRAY,
		Num:   noOfCommands,
		Array: make([]commands.Token, noOfCommands),
	}

	commandCounter := 0
	i := 2
	var returnString []string

	for commandCounter < noOfCommands {
		if i == len(stream) {
			break
		}
		if stream[i] == SKIP {
			i += 2
			continue
		}
		if stream[i] == BULK {
			commandCounter++
			i++
			var sizeString strings.Builder
			for stream[i] != '\r' {
				sizeString.WriteByte(stream[i])
				i++
			}
			i += 2
			size, err := strconv.Atoi(sizeString.String())
			if err != nil {
				return fmt.Sprintln(err)
			}
			var sb strings.Builder
			for j := i; j < size+i; j++ {
				sb.WriteByte(stream[j])
			}
			returnString = append(returnString, sb.String())
			i += size
			bulktoken := commands.Token{
				Type: BULK,
				Str:  sb.String(),
			}
			token.Array = append(token.Array, bulktoken)
		} else if stream[i] == INTEGER {
			commandCounter++
			i++
			var sb strings.Builder
			var sym byte
			for stream[i] != '\r' {
				if stream[i] == '-' || stream[i] == '+' {
					sym = stream[i]
				}
				sb.WriteByte(stream[i])
				i++
			}
			val, err := strconv.Atoi(sb.String())
			if err != nil {
				return fmt.Sprintln(err)
			}

			returnString = append(returnString, sb.String())

			inttoken := commands.Token{
				Type:   INTEGER,
				Num:    val,
				Symbol: sym,
			}
			token.Array = append(token.Array, inttoken)
			token.Array = append(token.Array, inttoken)
		} else if stream[i] == NULL {
			commandCounter++
			returnString = append(returnString, "NULL")
			i++
		} else if stream[i] == BOOLEAN {
			i++
			var sb strings.Builder
			var b bool
			if stream[i] == 't' {
				sb.WriteString("true")
				b = true
			} else if stream[i] == 'f' {
				sb.WriteString("false")
				b = false
			}

			boolToken := commands.Token{
				Type: BOOLEAN,
				Bool: b,
			}
			token.Array = append(token.Array, boolToken)
			returnString = append(returnString, sb.String())
			i++
		} else if stream[i] == DOUBLE {
			//i++
			//var symbol *byte
			//if stream[i] == '-' {
			//	symbol = &stream[i]
			//	i++
			//}

		}
	}

	return strings.Join(returnString, " ")
}

func ReadManager(stream []byte) string {
	dataType := stream[0]
	var text string

	if dataType == STRING {
		text := string(stream[1 : len(stream)-2])
		if text == "PING" {
			return commands.Handler("PING", make([]commands.Token, 0))
		}
		return fmt.Sprintf("+%s\r\n", text)
	} else if dataType == INTEGER {
		if stream[1] == '+' {
			text = string(stream[2 : len(stream)-2])
		} else if stream[1] == ERROR {
			text = fmt.Sprintf(":-%s\r\n", string(stream[2:len(stream)-2]))
		} else if stream[1] == ARRAY {
			text = manageArray(stream[2 : len(stream)-1])
		}
		return text
	}
	return "_\r\n"
}
