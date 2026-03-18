package chess_test

import (
	"testing"

	"gochess/chess"
)

func TestPieceWithColor(t *testing.T) {
	tests := []struct {
		inputPiece chess.Piece
		color      chess.Color
		expected   chess.Piece
	}{
		{chess.KING, chess.WHITE, chess.Piece(0b00001001)},
		{chess.QUEEN, chess.BLACK, chess.Piece(0b00010010)},
		{chess.BISHOP, chess.WHITE, chess.Piece(0b00001100)},
		{chess.KNIGHT, chess.BLACK, chess.Piece(0b00010101)},
	}
	for i, tt := range tests {
		result := tt.inputPiece.WithColor(tt.color)
		if result != tt.expected {
			t.Errorf("%v - Piece(%08b).WithColor(%08b) = %08b, got: %08b", i, tt.inputPiece, tt.color, result, tt.expected)
		}
	}
}

func TestPieceIsColor(t *testing.T) {
	tests := []struct {
		inputPiece chess.Piece
		color      chess.Color
		expected   bool
	}{
		{chess.KING.WithColor(chess.WHITE), chess.WHITE, true},
		{chess.QUEEN.WithColor(chess.BLACK), chess.BLACK, true},
		{chess.KNIGHT.WithColor(chess.WHITE), chess.BLACK, false},
		{chess.BISHOP.WithColor(chess.BLACK), chess.WHITE, false},
		{chess.Piece(0b00010110), chess.Color(0b00010000), true},
		{chess.Piece(0b00001001), chess.Color(0b00001000), true},
		{chess.Piece(0b00011000), chess.Color(0b00001000), false},
		{chess.Piece(0b00010110), chess.Color(0b00011000), false},
	}
	for i, tt := range tests {
		result := tt.inputPiece.IsColor(tt.color)
		if result != tt.expected {
			t.Errorf("%v - Piece(%08b).IsColor(%08b) = %v, want %v", i, tt.inputPiece, tt.color, result, tt.expected)
		}
	}
}

func TestPieceIsPieceSameColor(t *testing.T) {
	tests := []struct {
		inputPiece chess.Piece
		otherPiece chess.Piece
		expected   bool
	}{
		{chess.KING.WithColor(chess.WHITE), chess.PAWN.WithColor(chess.WHITE), true},
		{chess.QUEEN.WithColor(chess.BLACK), chess.KNIGHT.WithColor(chess.BLACK), true},
		{chess.BISHOP.WithColor(chess.WHITE), chess.PAWN.WithColor(chess.BLACK), false},
		{chess.ROOK.WithColor(chess.BLACK), chess.QUEEN.WithColor(chess.WHITE), false},
		{chess.Piece(0b00001011), chess.Piece(0b00001110), true},
		{chess.Piece(0b00010001), chess.Piece(0b00010110), true},
		{chess.Piece(0b00010001), chess.Piece(0b00001110), false},
		{chess.Piece(0b00001010), chess.Piece(0b00010011), false},
	}
	for i, tt := range tests {
		result := tt.inputPiece.IsPieceSameColor(tt.otherPiece)
		if result != tt.expected {
			t.Errorf("%v - Piece(%08b).IsPieceSameColor(%08b) = %v, wanted %v", i, tt.inputPiece, tt.otherPiece, result, tt.expected)
		}
	}
}
