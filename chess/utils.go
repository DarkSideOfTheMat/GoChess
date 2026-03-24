package chess

import "fmt"

// Square Helpers

func SquareAbs(x Square) Square {
	if x < 0 {
		return x * -1
	}
	return x
}

// File returns the 0 indexed file left to right
func (s Square) File() Square {
	return s % 8
}

// Rank returns the 0 indexed rank
func (s Square) Rank() Square {
	return s / 8
}

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

func (p Piece) ToString() string {
	pieceLabelMap := map[Piece]string{
		KING:   "K",
		QUEEN:  "Q",
		ROOK:   "R",
		BISHOP: "B",
		KNIGHT: "N",
		PAWN:   "P",
	}
	return pieceLabelMap[p.WithoutColor()]
}

func (p Piece) ToColor() Color {
	return Color(p &^ PieceMask)
}

func (p Piece) WithColor(color Color) Piece {
	return (p & PieceMask) | Piece(color)
}

func (p Piece) WithoutColor() Piece {
	return (p &^ ColorMask)
}

func (p Piece) IsColor(color Color) bool {
	return (p&^PieceMask)^Piece(color) == 0
}

func (p Piece) IsPieceSameColor(other Piece) bool {
	return p&^PieceMask == other&^PieceMask
}

// Color Utils and Helpers

// Flip returns the color of the other player, or flip to an empty mask
func (c Color) Flip() Color {
	return c ^ ColorMask.ToColor()
}

// Castling Utils and Helpers

func (c Castling) WithColor(color Color) Castling {
	return Castling(uint8(c) << (uint8(color) >> 2))
}

func (c Castling) WithoutColor() Castling {
	for c&CastleMask == 0 && c > 0 {
		c >>= 2
	}
	return c
}
