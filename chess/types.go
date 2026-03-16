// Package chess contains shared types and interfaces
package chess

import "fmt"

type GameStatus uint

const (
	Unknown GameStatus = iota
	Disconnected
	InProgress
	Draw
	WhiteWinsCheckmate
	BlackWinsCheckmate
	WhiteWinsClock
	BlackWinsClock
	WhiteWinsForfeit
	BlackWinsForfeit
)

// Square is the current index on the Board
type Square int

// Piece type represents the combination of pieces and colors
type Piece uint8

// Color type represents the color of the piece or player
//
//	uses 4th and 5th bit flag
//	white is b00001000 black is b00010000
type Color uint8

// PieceMask can be used to get a color
// in binary 00000111
const PieceMask Piece = 7

// ColorMask can be used to get a colorless piece type
// in binary 00011000
const ColorMask Piece = 24

// ValidateColor return the masked white or black of the
// also does basic error checking
func (p Piece) ValidateColor() (Piece, error) {
	maskedPiece := p &^ PieceMask
	if maskedPiece != ColorMask && maskedPiece != 0 {
		return maskedPiece, nil
	}
	return maskedPiece, fmt.Errorf("invalid piece color: %v! A piece cannot be both black and white", p)
}

func (p Piece) Color() Piece {
	color, _ := p.ValidateColor()
	return color
}

func (p Piece) IsColor(color Piece) bool {
	color = color &^ PieceMask
	return p&^PieceMask == color
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
	WHITE Color = 1 << 3 << iota
	BLACK
)

type Ply struct {
	Piece     Piece
	StartIdx  Square
	EndIdx    Square
	Promotion Piece
}

type Turn struct {
	WhitePly *Ply
	BlackPly *Ply
}
