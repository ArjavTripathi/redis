package lexer

import "testing"

func TestReadStream(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string
	}{
		{
			name:  "Short response",
			input: []byte("+OK\r\n"),
			want:  "OK",
		},
		{
			name:  "Longer response",
			input: []byte("+PONG\r\n"),
			want:  "PONG",
		},
		{
			name:  "Sentence response",
			input: []byte("+RECONSTRUCTED STREAM\r\n"),
			want:  "RECONSTRUCTED STREAM",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReadStream(tt.input)
			if got != tt.want {
				t.Errorf("ReadStream() = %q, want %q", got, tt.want)
			}
		})
	}
}
