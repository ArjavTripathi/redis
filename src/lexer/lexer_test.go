package lexer

import (
	"redis/src/commands"
	"testing"
)

//func TestHandleInteger(t *testing.T) {
//	tests := []struct {
//		name   string
//		stream []byte
//		index  *int
//		want   int
//	}{
//		{
//			name:   "positive single digit number",
//			stream: []byte(":1\r\n"),
//			index:  new(1),
//			want:   1,
//		},
//		{
//			name:   "positive double digit number",
//			stream: []byte(":100\r\n"),
//			index:  new(1),
//			want:   100,
//		},
//		{
//			name:   "negative single digit number",
//			stream: []byte(":-1\r\n"),
//			index:  new(1),
//			want:   -1,
//		},
//		{
//			name:   "negative double digit number",
//			stream: []byte(":-100\r\n"),
//			index:  new(1),
//			want:   -100,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got := handleInteger(tt.stream, tt.index)
//			if got != tt.want {
//				t.Errorf("manageArray() = \n%#v, \nwant \n%#v", got, tt.want)
//			}
//		})
//	}
//}

func TestManageArray(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  commands.Token
	}{
		{
			name:  "SET example",
			input: []byte("*1\r\n$3\r\nSET\r\n"),
			want: commands.Token{
				Type: ARRAY,
				Num:  1,
				Array: []commands.Token{
					{
						Type: BULK,
						Str:  "SET",
					},
				},
			},
		},
		{
			name:  "SET full example",
			input: []byte("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nhello\r\n"),
			want: commands.Token{
				Type: ARRAY,
				Num:  3,
				Array: []commands.Token{
					{
						Type: BULK,
						Str:  "SET",
					},
					{
						Type: BULK,
						Str:  "mykey",
					},
					{
						Type: BULK,
						Str:  "hello",
					},
				},
			},
		},
		{
			name:  "SET num example noSign",
			input: []byte("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n:223\r\n"),
			want: commands.Token{
				Type: ARRAY,
				Num:  3,
				Array: []commands.Token{
					{
						Type: BULK,
						Str:  "SET",
					},
					{
						Type: BULK,
						Str:  "mykey",
					},
					{
						Type:   INTEGER,
						Num:    223,
						Symbol: '+',
					},
				},
			},
		},
		{
			name:  "SET num example sign",
			input: []byte("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n:-5\r\n"),
			want: commands.Token{
				Type: ARRAY,
				Num:  3,
				Array: []commands.Token{
					{
						Type: BULK,
						Str:  "SET",
					},
					{
						Type: BULK,
						Str:  "mykey",
					},
					{
						Type:   INTEGER,
						Num:    -5,
						Symbol: '-',
					},
				},
			},
		},
		{
			name:  "SET example null",
			input: []byte("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n_\r\n"),
			want: commands.Token{
				Type: ARRAY,
				Num:  3,
				Array: []commands.Token{
					{
						Type: BULK,
						Str:  "SET",
					},
					{
						Type: BULK,
						Str:  "mykey",
					},
					{
						Type: NULL,
					},
				},
			},
		},
		{
			name:  "SET example boolean",
			input: []byte("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n#t\r\n"),
			want: commands.Token{
				Type: ARRAY,
				Num:  3,
				Array: []commands.Token{
					{
						Type: BULK,
						Str:  "SET",
					},
					{
						Type: BULK,
						Str:  "mykey",
					},
					{
						Type: BOOLEAN,
						Bool: true,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := manageArray(tt.input)
			if !got.Equals(tt.want) {
				t.Errorf("manageArray() = \n%#v, \nwant \n%#v", got, tt.want)
			}
		})
	}
}

//func TestReadManager(t *testing.T) {
//	tests := []struct {
//		name  string
//		input []byte
//		want  string
//	}{
//		{
//			name:  "Short response",
//			input: []byte("+PING\r\n"),
//			want:  "+PONG\r\n",
//		},
//		//{
//		//	name:  "Longer response",
//		//	input: []byte("*1\\r\\n$4\\r\\nPING\\r\\n"),
//		//	want:  "PONG",
//		//},
//		//{
//		//	name:  "Sentence response",
//		//	input: []byte("+RECONSTRUCTED STREAM\r\n"),
//		//	want:  "RECONSTRUCTED STREAM",
//		//},
//		//{
//		//	name:  "Negative int",
//		//	input: []byte(":-123\r\n"),
//		//	want:  "-123",
//		//},
//		//{
//		//	name:  "Positive int",
//		//	input: []byte(":+123\r\n"),
//		//	want:  "123",
//		//},
//		//{
//		//	name:  "Int",
//		//	input: []byte(":123\r\n"),
//		//	want:  "123",
//		//},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got := ReadManager(tt.input)
//			if got != tt.want {
//				t.Errorf("ReadStream() = %q, want %q", got, tt.want)
//			}
//		})
//	}
//}
