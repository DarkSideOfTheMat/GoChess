package ui

import (
	"gochess/chess"
)

var pieceSymbols = map[chess.Piece]string{
	chess.KING:   "♚",
	chess.QUEEN:  "♛",
	chess.ROOK:   "♜",
	chess.BISHOP: "♝",
	chess.KNIGHT: "♞",
	chess.PAWN:   "♟",
	0:            " ",
}

var fileLabels = [8]string{"a", "b", "c", "d", "e", "f", "g", "h"}

func getPieceSymbol(piece chess.Piece) string {
	if sym, ok := pieceSymbols[piece&^chess.ColorMask]; ok {
		return sym
	}
	return " "
}
