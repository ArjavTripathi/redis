package lexer

import "fmt"

func ReadStream(stream []byte) string {
	data_type := stream[0]

	if data_type == '+' {
		text := string(stream[1 : len(stream)-2])
		fmt.Println(text)
		return text
	}
	return "hello"
}
