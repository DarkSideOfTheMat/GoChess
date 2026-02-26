package game

// TODO: create piece representations and piece literals
// TODO: implement terminal renderer for pieces
type Piece uint8

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
