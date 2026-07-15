package lexer

import (
	"redis/src/commands"
	"testing"
)

func TestHandleBulkStrings(t *testing.T) {
	tests := []struct {
		name   string
		stream []byte
		index  int
		want   commands.Token
	}{
		{
			name:   "simple word",
			stream: []byte("$3\r\nSET\r\n"),
			index:  1,
			want:   commands.Token{Type: BULK, Str: "SET"},
		},
		{
			name:   "longer string",
			stream: []byte("$5\r\nhello\r\n"),
			index:  1,
			want:   commands.Token{Type: BULK, Str: "hello"},
		},
		{
			name:   "single char",
			stream: []byte("$1\r\nX\r\n"),
			index:  1,
			want:   commands.Token{Type: BULK, Str: "X"},
		},
		{
			name:   "string with spaces",
			stream: []byte("$9\r\nhello all\r\n"),
			index:  1,
			want:   commands.Token{Type: BULK, Str: "hello all"},
		},
		{
			name:   "multi-digit length",
			stream: []byte("$11\r\nhello world\r\n"),
			index:  1,
			want:   commands.Token{Type: BULK, Str: "hello world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := tt.index
			got := handleBulkStrings(tt.stream, &idx)
			if !got.Equals(tt.want) {
				t.Errorf("handleBulkStrings() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestHandleInteger(t *testing.T) {
	tests := []struct {
		name   string
		stream []byte
		index  int
		want   commands.Token
	}{
		{
			name:   "positive single digit",
			stream: []byte(":1\r\n"),
			index:  1,
			want:   commands.Token{Type: INTEGER, Num: 1, Symbol: '+'},
		},
		{
			name:   "positive multi digit",
			stream: []byte(":100\r\n"),
			index:  1,
			want:   commands.Token{Type: INTEGER, Num: 100, Symbol: '+'},
		},
		{
			name:   "negative single digit",
			stream: []byte(":-1\r\n"),
			index:  1,
			want:   commands.Token{Type: INTEGER, Num: -1, Symbol: '-'},
		},
		{
			name:   "negative multi digit",
			stream: []byte(":-100\r\n"),
			index:  1,
			want:   commands.Token{Type: INTEGER, Num: -100, Symbol: '-'},
		},
		{
			name:   "zero",
			stream: []byte(":0\r\n"),
			index:  1,
			want:   commands.Token{Type: INTEGER, Num: 0, Symbol: '+'},
		},
		{
			name:   "large positive",
			stream: []byte(":99999\r\n"),
			index:  1,
			want:   commands.Token{Type: INTEGER, Num: 99999, Symbol: '+'},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := tt.index
			got := handleInteger(tt.stream, &idx)
			if !got.Equals(tt.want) {
				t.Errorf("handleInteger() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestHandleNull(t *testing.T) {
	got := handleNull()
	want := commands.Token{Type: NULL}
	if !got.Equals(want) {
		t.Errorf("handleNull() = %#v, want %#v", got, want)
	}
}

func TestHandleBoolean(t *testing.T) {
	tests := []struct {
		name   string
		stream []byte
		index  int
		want   commands.Token
	}{
		{
			name:   "true",
			stream: []byte("#t\r\n"),
			index:  1,
			want:   commands.Token{Type: BOOLEAN, Bool: true},
		},
		{
			name:   "false",
			stream: []byte("#f\r\n"),
			index:  1,
			want:   commands.Token{Type: BOOLEAN, Bool: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := tt.index
			got := handleBoolean(tt.stream, &idx)
			if !got.Equals(tt.want) {
				t.Errorf("handleBoolean() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestManageArray(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  commands.Token
	}{
		{
			name:  "single bulk string",
			input: []byte("*1\r\n$3\r\nSET\r\n"),
			want: commands.Token{
				Type: ARRAY, Num: 1,
				Array: []commands.Token{
					{Type: BULK, Str: "SET"},
				},
			},
		},
		{
			name:  "three bulk strings",
			input: []byte("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nhello\r\n"),
			want: commands.Token{
				Type: ARRAY, Num: 3,
				Array: []commands.Token{
					{Type: BULK, Str: "SET"},
					{Type: BULK, Str: "mykey"},
					{Type: BULK, Str: "hello"},
				},
			},
		},
		{
			name:  "bulk strings with positive integer",
			input: []byte("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n:223\r\n"),
			want: commands.Token{
				Type: ARRAY, Num: 3,
				Array: []commands.Token{
					{Type: BULK, Str: "SET"},
					{Type: BULK, Str: "mykey"},
					{Type: INTEGER, Num: 223, Symbol: '+'},
				},
			},
		},
		{
			name:  "bulk strings with negative integer",
			input: []byte("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n:-5\r\n"),
			want: commands.Token{
				Type: ARRAY, Num: 3,
				Array: []commands.Token{
					{Type: BULK, Str: "SET"},
					{Type: BULK, Str: "mykey"},
					{Type: INTEGER, Num: -5, Symbol: '-'},
				},
			},
		},
		{
			name:  "bulk strings with null",
			input: []byte("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n_\r\n"),
			want: commands.Token{
				Type: ARRAY, Num: 3,
				Array: []commands.Token{
					{Type: BULK, Str: "SET"},
					{Type: BULK, Str: "mykey"},
					{Type: NULL},
				},
			},
		},
		{
			name:  "bulk strings with boolean true",
			input: []byte("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n#t\r\n"),
			want: commands.Token{
				Type: ARRAY, Num: 3,
				Array: []commands.Token{
					{Type: BULK, Str: "SET"},
					{Type: BULK, Str: "mykey"},
					{Type: BOOLEAN, Bool: true},
				},
			},
		},
		{
			name:  "bulk strings with boolean false",
			input: []byte("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n#f\r\n"),
			want: commands.Token{
				Type: ARRAY, Num: 3,
				Array: []commands.Token{
					{Type: BULK, Str: "SET"},
					{Type: BULK, Str: "mykey"},
					{Type: BOOLEAN, Bool: false},
				},
			},
		},
		{
			name:  "zero commands returns error token",
			input: []byte("*0\r\n"),
			want:  commands.ErrorToken("0 commands found in Array"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := manageArray(tt.input)
			if !got.Equals(tt.want) {
				t.Errorf("manageArray() =\n%#v\nwant\n%#v", got, tt.want)
			}
		})
	}
}

func TestReadManager(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  commands.Token
	}{
		{
			name:  "bulk string",
			input: []byte("$5\r\nhello\r\n"),
			want:  commands.Token{Type: BULK, Str: "hello"},
		},
		{
			name:  "positive integer",
			input: []byte(":42\r\n"),
			want:  commands.Token{Type: INTEGER, Num: 42, Symbol: '+'},
		},
		{
			name:  "negative integer",
			input: []byte(":-7\r\n"),
			want:  commands.Token{Type: INTEGER, Num: -7, Symbol: '-'},
		},
		{
			name:  "null",
			input: []byte("_\r\n"),
			want:  commands.Token{Type: NULL},
		},
		{
			name:  "array of bulk strings",
			input: []byte("*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n"),
			want: commands.Token{
				Type: ARRAY, Num: 2,
				Array: []commands.Token{
					{Type: BULK, Str: "GET"},
					{Type: BULK, Str: "foo"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReadManager(tt.input)
			if !got.Equals(tt.want) {
				t.Errorf("ReadManager() =\n%#v\nwant\n%#v", got, tt.want)
			}
		})
	}
}
