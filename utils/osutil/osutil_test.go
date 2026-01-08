package osutil

import (
	"os"
	"strings"
	"testing"
)

func TestCallSuccess(t *testing.T) {
	fn := func(a int, b string) (int, string) {
		return a + 1, b + "x"
	}
	results, err := Call(fn, 1, "a")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 return values, got %d", len(results))
	}
	if got, want := results[0].Interface(), 2; got != want {
		t.Fatalf("first return value = %v, want %v", got, want)
	}
	if got, want := results[1].Interface(), "ax"; got != want {
		t.Fatalf("second return value = %v, want %v", got, want)
	}
}

func TestCallNonFunction(t *testing.T) {
	_, err := Call(123)
	if err == nil {
		t.Fatalf("expected error when calling non-function, got nil")
	}
}

func TestCallParamCountMismatch(t *testing.T) {
	fn := func(a int) {}
	_, err := Call(fn)
	if err == nil {
		t.Fatalf("expected error on parameter count mismatch, got nil")
	}
}

func TestGetProcessName(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"/usr/local/bin/myproc"}
	got := GetProcessName()
	if got != "myproc" {
		t.Fatalf("GetProcessName() = %q, want %q", got, "myproc")
	}
}

func TestGetCurrentGoroutineIDFromStack(t *testing.T) {
	id := GetCurrentGoroutineIDFromStack()
	if id == "" {
		t.Fatalf("expected non-empty goroutine id")
	}
	for _, r := range id {
		if !strings.ContainsRune("0123456789", r) {
			t.Fatalf("expected goroutine id to be numeric, got %q", id)
		}
	}
}

