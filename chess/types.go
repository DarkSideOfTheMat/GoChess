// Package chess contains shared types and interfaces
package chess

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

const (
	WhiteKingHome Square = 4  // e1
	BlackKingHome Square = 60 // e8
)

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

// Castling takes the front 3 bits of the piece mask
// it should only be used for tracking what kind of castiling is availble
//
// for convience we will have uncolored and the colored variation can be calcualted with
// Side<<(Color>>2)
type Castling uint8

const (
	CastleKingSide = 1 << iota
	CastleQueenSide
	CastleWhiteKingSide
	CastleWhiteQueenSide
	CastleBlackKingSide
	CastleBlackQueenSide
)

const (
	CastleMask Castling = 0b11 << (2 * iota)
	CastleWhiteMask
	CastleBlackMask
)

// BoardState represents the 64-square board as a flat array.
//
// -+ ------------------------------ +
// 8| 56, 57, 58, 59, 60, 61, 62, 63 |
// 7| 48, 49, 50, 51, 52, 53, 54, 55,|
// 6| 40, 41, 42, 43, 44, 45, 46, 47,|
// 5| 32, 33, 34, 35, 36, 37, 38, 39,|
// 4| 24, 25, 26, 27, 28, 29, 30, 31,|
// 3| 16, 17, 18, 19, 20, 21, 22, 23,|
// 2|  8,  9, 10, 11, 12, 13, 14, 15,|
// 1|  0,  1,  2,  3,  4,  5,  6,  7,|
// -+ -------------------------------+
// .   a   b   c   d   e   f   g   h
type BoardState [64]Piece

// Board represents the current state of a chess position.
type Board struct {
	State           BoardState
	ActiveColor     Color
	Castling        Castling
	EnpassantTarget Square // square that was moved through by pawn moving 2 squares, -1 otherwise
}

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
