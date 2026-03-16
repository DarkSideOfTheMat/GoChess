package chess

import "fmt"

// FileNameFromInt returns the file name associated with the current square
//
//	for example, square 0 will return 'a', square 15 will return 'h' etc
func FileNameFromInt(square int) rune {
	r := 'a'
	return rune(int(r) + (square % 8))
}

func (s Square) fileName() rune {
	return FileNameFromInt(int(s))
}

func (s Square) ToString() string {
	rank := s / 8

	return fmt.Sprintf("%c%v", s.fileName(), rank+1)
}

// ParseSquare converts algebraic notation like "e2" to a Square index
func ParseSquare(s string) (Square, error) {
	if len(s) != 2 {
		return 0, fmt.Errorf("invalid square: %s", s)
	}
	file := int(s[0] - 'a')
	rank := int(s[1] - '1')
	if file < 0 || file > 7 || rank < 0 || rank > 7 {
		return 0, fmt.Errorf("invalid square: %s", s)
	}
	return Square(rank*8 + file), nil
}

// Piece and Color Helpers

func (p Piece) ToColor() Color {
	return Color(p &^ PieceMask)
}
