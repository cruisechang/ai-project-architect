package buildinfo

import (
	"testing"
)

func TestStringDefault(t *testing.T) {
	got := String()
	want := "version=dev commit=none date=unknown"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestStringCustomValues(t *testing.T) {
	origV, origC, origD := Version, Commit, BuildDate
	t.Cleanup(func() { Version, Commit, BuildDate = origV, origC, origD })

	Version = "1.2.3"
	Commit = "abc123"
	BuildDate = "2026-01-01"

	got := String()
	want := "version=1.2.3 commit=abc123 date=2026-01-01"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
