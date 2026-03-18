package game

import (
	"fmt"

	"gochess/chess"
)

func CastleFromColorAndSide(color chess.Color, side chess.Castling) (chess.Castling, error) {
	if side&^chess.CastleMask != 0 {
		return 0, fmt.Errorf("must choose an uncolored side to castle with, got: %b", side)
	}
	if color != chess.WHITE && color != chess.BLACK {
		return 0, fmt.Errorf("invalid color choice, got: %b", color)
	}
	return chess.Castling(uint8(side) << (uint8(color) >> 2)), nil
}

func IsCastlingMove(piece chess.Piece, from chess.Square, to chess.Square) (bool, chess.Castling, error) {
	delta := to - from
	// player is attempting to IsCastlingMove
	// note, golang doesn't have a builtin int absolute value func
	if piece&chess.PieceMask == chess.KING && (delta == 2 || delta == -2) {
		var side chess.Castling
		if to > from {
			side = chess.CastleKingSide
		} else {
			side = chess.CastleQueenSide
		}
		castling, err := CastleFromColorAndSide(piece.ToColor(), side)
		return true, castling, err
	}
	return false, 0, nil
}

// Board and Square Helpers

func StartingSquareByFileIdx(square chess.Square, color chess.Color) chess.Square {
	file := square % 8

	// black is 16 so we get
	return file + chess.Square(8*7*(uint8(color)/16))
}
