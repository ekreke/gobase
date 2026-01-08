package stringutil

import (
	"math/rand"
	"strings"
	"unicode"
)

const (
	DefaultCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?"
	Lowercase      = "abcdefghijklmnopqrstuvwxyz"
	Uppercase      = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits         = "0123456789"
	Symbols        = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

// SnakeCase convert string to snake case
// exp: "HelloWorld" -> "hello_world"
func SnakeCase(str string) string {
	var snakeCase strings.Builder
	for i, r := range str {
		if unicode.IsUpper(r) && i > 0 {
			snakeCase.WriteRune('_')
		}
		snakeCase.WriteRune(unicode.ToLower(r))
	}
	return snakeCase.String()
}

// CamelCase convert string to camel case
// exp: "hello_world" -> "helloWorld"
func CamelCase(str string) string {
	ws := strings.Split(str, "_")
	var b strings.Builder
	wroteFirst := false
	for _, w := range ws {
		if w == "" {
			continue
		}

		if !wroteFirst {
			b.WriteString(strings.ToLower(w))
			wroteFirst = true
			continue
		}

		lowered := strings.ToLower(w)
		runes := []rune(lowered)
		if len(runes) == 0 {
			continue
		}
		runes[0] = unicode.ToUpper(runes[0])
		b.WriteString(string(runes))
	}
	return b.String()
}

// RandomString generator with default charset
func RandomString(length int) string {
	return RandomStringWithCharset(length, DefaultCharset)
}

// RandomStringWithCharset generator with custom charset
func RandomStringWithCharset(length int, charset string) string {
	if length <= 0 {
		return ""
	}
	if charset == "" {
		return ""
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// CamelCaseToSpaceSeparated convert camel case string to space separated string
// exp: "helloWorld" -> "hello World"
func CamelCaseToSpaceSeparated(str string) string {
	var b strings.Builder
	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) {
			b.WriteRune(' ')
		}
		b.WriteRune(r)
	}
	return b.String()
}

// UpperFirst upper first letter of string
func UpperFirst(input string) string {
	if input == "" {
		return ""
	}
	runes := []rune(input)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// LowerFirst lower first letter of string
func LowerFirst(input string) string {
	if input == "" {
		return ""
	}
	runes := []rune(input)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// FormatString format string to fixed length
func FormatString(text string, length int, align string) string {
	if len(text) >= length {
		return text
	}
	if len(text) >= length {
		return text[:length]
	}

	if align == "left" {
		return text + strings.Repeat(" ", length-len(text))
	}
	return strings.Repeat(" ", length-len(text)) + text
}
