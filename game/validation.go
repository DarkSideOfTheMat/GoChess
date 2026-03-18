package game

import (
	"fmt"

	"gochess/chess"
)

// ValidatePieceColor returns the color bits of the piece, or an error if the piece
// has both color flags set (which is invalid).
func ValidatePieceColor(p chess.Piece) (chess.Piece, error) {
	maskedPiece := p &^ chess.PieceMask
	if maskedPiece != chess.ColorMask && maskedPiece != 0 {
		return maskedPiece, nil
	}
	return maskedPiece, fmt.Errorf("invalid piece color: %v! A piece cannot be both black and white", p)
}

// PieceColorBits returns the color bits of a piece, ignoring validation errors.
func PieceColorBits(p chess.Piece) chess.Piece {
	color, _ := ValidatePieceColor(p)
	return color
}
