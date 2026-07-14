package lexer

import "testing"

func TestManageArray(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string
	}{
		{
			name:  "SET example",
			input: []byte("*1\r\n$3\r\nSET\r\n"),
			want:  "SET",
		},
		{
			name:  "SET full example",
			input: []byte("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nhello\r\n"),
			want:  "SET mykey hello",
		},
		{
			name:  "SET num example noSign",
			input: []byte("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n:3\r\n"),
			want:  "SET mykey 3",
		},
		{
			name:  "SET num example sign",
			input: []byte("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n:-5\r\n"),
			want:  "SET mykey -5",
		},
		{
			name:  "SET example null",
			input: []byte("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n_\r\n"),
			want:  "SET mykey NULL",
		},
		{
			name:  "SET example boolean",
			input: []byte("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n#t\r\n"),
			want:  "SET mykey true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := manageArray(tt.input)
			if got != tt.want {
				t.Errorf("manageArray() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReadManager(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string
	}{
		{
			name:  "Short response",
			input: []byte("+PING\r\n"),
			want:  "+PONG\r\n",
		},
		//{
		//	name:  "Longer response",
		//	input: []byte("*1\\r\\n$4\\r\\nPING\\r\\n"),
		//	want:  "PONG",
		//},
		//{
		//	name:  "Sentence response",
		//	input: []byte("+RECONSTRUCTED STREAM\r\n"),
		//	want:  "RECONSTRUCTED STREAM",
		//},
		//{
		//	name:  "Negative int",
		//	input: []byte(":-123\r\n"),
		//	want:  "-123",
		//},
		//{
		//	name:  "Positive int",
		//	input: []byte(":+123\r\n"),
		//	want:  "123",
		//},
		//{
		//	name:  "Int",
		//	input: []byte(":123\r\n"),
		//	want:  "123",
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReadManager(tt.input)
			if got != tt.want {
				t.Errorf("ReadStream() = %q, want %q", got, tt.want)
			}
		})
	}
}
