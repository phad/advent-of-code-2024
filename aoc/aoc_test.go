package aoc

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReadLines(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "in.txt")
	if err := os.WriteFile(f, []byte("a\nbb\nccc\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	got, err := ReadLines(f)
	if err != nil {
		t.Fatalf("ReadLines err=%v", err)
	}
	want := []string{"a", "bb", "ccc"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ReadLines = %v, want %v", got, want)
	}
}

func TestReadLinesMissingFile(t *testing.T) {
	if _, err := ReadLines(filepath.Join(t.TempDir(), "nope.txt")); err == nil {
		t.Error("ReadLines on missing file: want error, got nil")
	}
}

func TestMustParseInt(t *testing.T) {
	if got := MustParseInt("42"); got != 42 {
		t.Errorf("MustParseInt(\"42\") = %d, want 42", got)
	}
	if got := MustParseInt("-7"); got != -7 {
		t.Errorf("MustParseInt(\"-7\") = %d, want -7", got)
	}
}
