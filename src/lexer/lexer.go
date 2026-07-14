package lexer

import (
	"fmt"
	"strconv"
	"strings"
)

func manageArray(stream []byte) string {
	stringRep := string(stream)
	noOfCommands := int(stream[1] - '0')
	fmt.Println(stringRep)

	if noOfCommands == 0 {
		return ""
	}
	commandCounter := 0
	i := 2
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
			i = i + size
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
