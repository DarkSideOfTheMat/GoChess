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

// Piece and Color Helpers

func (p Piece) ToColor() Color {
	return Color(p &^ PieceMask)
}
