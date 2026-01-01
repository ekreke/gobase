package stringutil

import (
	"testing"
)

func TestSnakeCase(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: ""},
		{name: "singleWordUpper", in: "Hello", want: "hello"},
		{name: "helloWorld", in: "HelloWorld", want: "hello_world"},
		{name: "alreadyLower", in: "hello", want: "hello"},
		{name: "mixed", in: "helloWorldAgain", want: "hello_world_again"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SnakeCase(tt.in)
			if got != tt.want {
				t.Fatalf("SnakeCase(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestCamelCase(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: ""},
		{name: "single", in: "hello", want: "hello"},
		{name: "snake", in: "hello_world", want: "helloWorld"},
		{name: "multiple", in: "hello_world_again", want: "helloWorldAgain"},
		{name: "alreadyCamelNoUnderscore", in: "helloWorld", want: "helloworld"},
		{name: "leadingTrailingUnderscore", in: "_hello_world_", want: "helloWorld"},
		{name: "doubleUnderscore", in: "hello__world", want: "helloWorld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CamelCase(tt.in)
			if got != tt.want {
				t.Fatalf("CamelCase(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestRandomStringWithCharset(t *testing.T) {
	t.Run("nonPositiveLength", func(t *testing.T) {
		if got := RandomStringWithCharset(0, "abc"); got != "" {
			t.Fatalf("expected empty string, got %q", got)
		}
		if got := RandomStringWithCharset(-1, "abc"); got != "" {
			t.Fatalf("expected empty string, got %q", got)
		}
	})

	t.Run("emptyCharset", func(t *testing.T) {
		if got := RandomStringWithCharset(5, ""); got != "" {
			t.Fatalf("expected empty string, got %q", got)
		}
	})

	t.Run("lengthAndAlphabet", func(t *testing.T) {
		charset := "ab"
		length := 64
		got := RandomStringWithCharset(length, charset)
		if len(got) != length {
			t.Fatalf("expected length %d, got %d", length, len(got))
		}
		for i := 0; i < len(got); i++ {
			if got[i] != 'a' && got[i] != 'b' {
				t.Fatalf("unexpected character %q at index %d in %q", got[i], i, got)
			}
		}
	})
}

func TestRandomString(t *testing.T) {
	length := 32
	got := RandomString(length)
	if len(got) != length {
		t.Fatalf("expected length %d, got %d", length, len(got))
	}
	for i := 0; i < len(got); i++ {
		if !stringsContainsByte(DefaultCharset, got[i]) {
			t.Fatalf("unexpected character %q at index %d in %q", got[i], i, got)
		}
	}
}

func stringsContainsByte(s string, b byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return true
		}
	}
	return false
}
