package game

import "fmt"

// TODO: create piece representations and piece literals
// TODO: implement terminal renderer for pieces
type Piece uint8

// 00000111
const PIECE_MASK = 7

// 00011000
const COLOR_MASK = 24

// return the masked white or black of the
// also does basic error checking
func (p Piece) Color() (Piece, error) {
	masked_p := p &^ PIECE_MASK
	if masked_p != COLOR_MASK && masked_p != 0 {
		return masked_p, nil
	}
	return masked_p, fmt.Errorf("Invalid piece color: %v! A piece cannot be both black and white.", p)
}

func (p Piece) IsColor(col Piece) bool {
	p_col, _ := p.Color() // intentionally ignoring this error since we want it to panic
	return p_col == col
}

// chess pieces
// .  King,     Queen,     Rook,   Bishop,   Knight,     Pawn
// 00000001, 00000010, 00000011, 00000100, 00000101, 00000110
const (
	KING Piece = 1 + iota
	QUEEN
	ROOK
	BISHOP
	KNIGHT
	PAWN
)

// colors
// .  White,    Black
// 00001000, 00010000
const (
	WHITE Piece = 1 << 3 << iota
	BLACK
)
