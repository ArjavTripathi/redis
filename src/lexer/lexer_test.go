package lexer

import "testing"

func TestManageArray(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string
	}{
		{ //"*3$3SET$5mykey$5hello"
			name:  "SET example",
			input: []byte("*1\r\n$3\r\nSET"),
			want:  "SET mykey hello",
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

//func TestReadStream(t *testing.T) {
//	tests := []struct {
//		name  string
//		input []byte
//		want  string
//	}{
//		{
//			name:  "Short response",
//			input: []byte("+OK\r\n"),
//			want:  "OK",
//		},
//		{
//			name:  "Longer response",
//			input: []byte("*1\\r\\n$4\\r\\nPING\\r\\n"),
//			want:  "PONG",
//		},
//		{
//			name:  "Sentence response",
//			input: []byte("+RECONSTRUCTED STREAM\r\n"),
//			want:  "RECONSTRUCTED STREAM",
//		},
//		{
//			name:  "Negative int",
//			input: []byte(":-123\r\n"),
//			want:  "-123",
//		},
//		{
//			name:  "Positive int",
//			input: []byte(":+123\r\n"),
//			want:  "123",
//		},
//		{
//			name:  "Int",
//			input: []byte(":123\r\n"),
//			want:  "123",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got := ReadStream(tt.input)
//			if got != tt.want {
//				t.Errorf("ReadStream() = %q, want %q", got, tt.want)
//			}
//		})
//	}
