package formatutil

import (
	"strings"
	"testing"
)

func TestProgressBarZeroTotal(t *testing.T) {
	got := ProgressBar("task", 0, 0)
	if !strings.Contains(got, "task") {
		t.Fatalf("expected output to contain name, got %q", got)
	}
	if !strings.Contains(got, "0/0") {
		t.Fatalf("expected output to contain progress fraction 0/0, got %q", got)
	}
	if strings.Contains(got, "█") {
		t.Fatalf("expected no progress blocks for zero total, got %q", got)
	}
}

func TestProgressBarFullProgress(t *testing.T) {
	got := ProgressBar("task", 100, 100)
	if !strings.Contains(got, "100/100") {
		t.Fatalf("expected output to contain progress fraction 100/100, got %q", got)
	}
	if !strings.Contains(got, "100%") {
		t.Fatalf("expected output to contain 100%%, got %q", got)
	}
	if c := strings.Count(got, "█"); c != 50 {
		t.Fatalf("expected 50 progress blocks, got %d in %q", c, got)
	}
}

func TestProgressBarDoesNotExceedBarLength(t *testing.T) {
	got := ProgressBar("task", 150, 100)
	if c := strings.Count(got, "█"); c != 50 {
		t.Fatalf("expected capped progress blocks at 50, got %d in %q", c, got)
	}
}

func TestProgressBarPartialProgress(t *testing.T) {
	got := ProgressBar("task", 25, 100)
	if !strings.Contains(got, "25/100") {
		t.Fatalf("expected output to contain progress fraction 25/100, got %q", got)
	}
	if c := strings.Count(got, "█"); c != 12 {
		t.Fatalf("expected 12 progress blocks for 25%% of 50, got %d in %q", c, got)
	}
}
