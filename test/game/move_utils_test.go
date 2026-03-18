package game_test

import (
	"testing"

	"gochess/chess"
	"gochess/game"
)

func TestCastleFromColorAndSide(t *testing.T) {
	t.Run("happy paths", func(t *testing.T) {
		tests := []struct {
			name     string
			color    chess.Color
			side     chess.Castling
			expected chess.Castling
		}{
			{"white kingside", chess.WHITE, chess.CastleKingSide, chess.CastleWhiteKingSide},
			{"white queenside", chess.WHITE, chess.CastleQueenSide, chess.CastleWhiteQueenSide},
			{"black kingside", chess.BLACK, chess.CastleKingSide, chess.CastleBlackKingSide},
			{"black queenside", chess.BLACK, chess.CastleQueenSide, chess.CastleBlackQueenSide},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := game.CastleFromColorAndSide(tt.color, tt.side)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got != tt.expected {
					t.Errorf("got %d, want %d", got, tt.expected)
				}
			})
		}
	})

	t.Run("error cases", func(t *testing.T) {
		tests := []struct {
			name  string
			color chess.Color
			side  chess.Castling
		}{
			{"both color bits set", chess.Color(chess.WHITE | chess.BLACK), chess.CastleKingSide},
			{"no color bits set", chess.Color(0), chess.CastleKingSide},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := game.CastleFromColorAndSide(tt.color, tt.side)
				if err == nil {
					t.Errorf("expected error for %s, got nil", tt.name)
				}
			})
		}

		t.Run("already colored side", func(t *testing.T) {
			// CastleWhiteKingSide (4) has bits outside CastleMask, should error
			_, err := game.CastleFromColorAndSide(chess.WHITE, chess.Castling(chess.CastleWhiteKingSide))
			if err == nil {
				t.Error("expected error for already colored side, got nil")
			}
		})
	})
}

func TestIsCastlingMove(t *testing.T) {
	t.Run("castling detected", func(t *testing.T) {
		tests := []struct {
			name     string
			piece    chess.Piece
			from     chess.Square
			to       chess.Square
			expected chess.Castling
		}{
			{"white kingside", chess.KING | chess.Piece(chess.WHITE), 4, 6, chess.CastleWhiteKingSide},
			{"white queenside", chess.KING | chess.Piece(chess.WHITE), 4, 2, chess.CastleWhiteQueenSide},
			{"black kingside", chess.KING | chess.Piece(chess.BLACK), 60, 62, chess.CastleBlackKingSide},
			{"black queenside", chess.KING | chess.Piece(chess.BLACK), 60, 58, chess.CastleBlackQueenSide},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				isCastling, castling, err := game.IsCastlingMove(tt.piece, tt.from, tt.to)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !isCastling {
					t.Fatal("expected castling move, got false")
				}
				if castling != tt.expected {
					t.Errorf("got castling flag %d, want %d", castling, tt.expected)
				}
			})
		}
	})

	t.Run("not castling", func(t *testing.T) {
		tests := []struct {
			name  string
			piece chess.Piece
			from  chess.Square
			to    chess.Square
		}{
			{"king moves one square", chess.KING | chess.Piece(chess.WHITE), 4, 5},
			{"pawn move", chess.PAWN | chess.Piece(chess.WHITE), 12, 28},
			{"queen move", chess.QUEEN | chess.Piece(chess.WHITE), 3, 5},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				isCastling, castling, err := game.IsCastlingMove(tt.piece, tt.from, tt.to)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if isCastling {
					t.Error("expected non-castling move, got true")
				}
				if castling != 0 {
					t.Errorf("expected castling flag 0, got %d", castling)
				}
			})
		}
	})
}
