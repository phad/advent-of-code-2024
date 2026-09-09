package aoc

import "testing"

func TestPointString(t *testing.T) {
	if got := (Point{3, 4}).String(); got != "(3,4)" {
		t.Errorf("Point.String = %q, want (3,4)", got)
	}
}

func TestNewRuneGrid(t *testing.T) {
	g, err := NewRuneGrid([]string{"AB", "CD"}, true)
	if err != nil {
		t.Fatalf("NewRuneGrid err=%v", err)
	}
	if g.W != 2 || g.H != 2 {
		t.Errorf("dims = %dx%d, want 2x2", g.W, g.H)
	}
	if v, ok := g.At(Point{1, 0}); !ok || v != 'B' {
		t.Errorf("At(1,0) = %q,%v want 'B',true", v, ok)
	}
	if _, ok := g.At(Point{2, 0}); ok {
		t.Error("At(2,0) in bounds, want out of bounds")
	}
	if p, ok := g.Find('D'); !ok || p != (Point{1, 1}) {
		t.Errorf("Find('D') = %v,%v want (1,1),true", p, ok)
	}
}

func TestNewRuneGridNonSquare(t *testing.T) {
	if _, err := NewRuneGrid([]string{"ABC", "DE"}, false); err == nil {
		t.Error("ragged grid: want error, got nil")
	}
	if _, err := NewRuneGrid([]string{"ABC"}, true); err == nil {
		t.Error("non-square with wantSquare: want error, got nil")
	}
}

func TestNewDigitGrid(t *testing.T) {
	g, err := NewDigitGrid([]string{"12", "34"}, true)
	if err != nil {
		t.Fatalf("NewDigitGrid err=%v", err)
	}
	if v, ok := g.At(Point{0, 1}); !ok || v != 3 {
		t.Errorf("At(0,1) = %d,%v want 3,true", v, ok)
	}
}

func TestGridSetSwapCount(t *testing.T) {
	g := NewGrid[bool](3, 2)
	if !g.Set(Point{1, 1}, true) {
		t.Error("Set in bounds returned false")
	}
	if g.Set(Point{5, 5}, true) {
		t.Error("Set out of bounds returned true")
	}
	if n := g.Count(true); n != 1 {
		t.Errorf("Count(true) = %d, want 1", n)
	}
	g.Set(Point{0, 0}, true)
	g.Swap(Point{0, 0}, Point{2, 0})
	if v, _ := g.At(Point{2, 0}); !v {
		t.Error("after swap, (2,0) should be true")
	}
	if v, _ := g.At(Point{0, 0}); v {
		t.Error("after swap, (0,0) should be false")
	}
}

func TestGridString(t *testing.T) {
	g, _ := NewRuneGrid([]string{"AB", "CD"}, true)
	want := "width:2 height:2\nAB\nCD\n"
	if got := g.String(); got != want {
		t.Errorf("String = %q, want %q", got, want)
	}
	if got := g.Highlight('A'); got != "width:2 height:2\nA.\n..\n" {
		t.Errorf("Highlight('A') = %q", got)
	}
}
