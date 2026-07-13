package lexer

import (
	"fmt"
	"strings"
)

func manageArray(stream []byte) string {
	noOfCommands := int(stream[0] - '0')

	if noOfCommands == 0 {
		return ""
	}
	commandCounter := 0
	i := 1
	var returnString []string

	for commandCounter < noOfCommands {
		if stream[i] == '\r' {
			i = i + 2
			if i == len(stream) {
				break
			}
			continue
		}

		if stream[i] == '$' {
			
		}
	}

	return strings.Join(returnString, " ")
}

func ReadManager(stream []byte) string {
	dataType := stream[0]
	var text string

	if dataType == '+' {
		text := fmt.Sprintf("+%s\r\n", string(stream[1:len(stream)-2]))
		return text
	} else if dataType == ':' {
		if stream[1] == '+' {
			text = string(stream[2 : len(stream)-2])
		} else if stream[1] == '-' {
			text = fmt.Sprintf(":-%s\r\n", string(stream[2:len(stream)-2]))
		} else if stream[1] == '*' {
			text = manageArray(stream[2 : len(stream)-1])
		}
		return text
	}
	return "neither detected"
}
