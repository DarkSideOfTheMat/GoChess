package game_test

import (
	"testing"

	"gochess/chess"
	"gochess/game"
)

// These tests verify that Game.Castling state is correctly updated when
// MakeMove is called. They will FAIL until castling revocation logic is
// added to MakeMove.

func TestCastlingRevokedOnKingMove(t *testing.T) {
	g := game.NewGame()
	var err error
	startingMoves := []struct {
		from chess.Square
		to   chess.Square
	}{
		{chess.Square(12), chess.Square(20)}, // e2 e3
		{chess.Square(52), chess.Square(36)}, // e7 e5
	}

	for _, move := range startingMoves {
		err = g.MakeMove(move.from, move.to, 0)
		if err != nil {
			t.Fatalf("board setup failed: %v", err)
		}
	}

	// White king e1→e2
	err = g.MakeMove(chess.Square(4), chess.Square(12), 0)
	if err != nil {
		t.Fatalf("white king move failed: %v", err)
	}
	if g.Castling&chess.CastleWhiteMask != 0 {
		t.Error("white castling rights should be cleared after king move")
	}
	if g.Castling&chess.CastleBlackMask == 0 {
		t.Error("black castling rights should be preserved after white king move")
	}

	// Black king e8→e7
	err = g.MakeMove(chess.Square(60), chess.Square(52), 0)
	if err != nil {
		t.Fatalf("black king move failed: %v", err)
	}
	if g.Castling&chess.CastleBlackMask != 0 {
		t.Error("black castling rights should be cleared after king move")
	}
}

func TestCastlingRevokedOnKingsideRookMove(t *testing.T) {
	g := game.NewGame()

	// White rook h1→h3 (square 7→23)
	// Need to clear the pawn at h2 first for a legal path
	g.Board.State[15] = 0 // remove pawn at h2
	err := g.MakeMove(chess.Square(7), chess.Square(23), 0)
	if err != nil {
		t.Fatalf("white rook move failed: %v", err)
	}
	if g.Castling&chess.CastleWhiteKingSide != 0 {
		t.Error("white kingside castling right should be revoked after h1 rook move")
	}

	// Black rook h8→h6 (square 63→47)
	g.Board.State[55] = 0 // remove pawn at h7
	err = g.MakeMove(chess.Square(63), chess.Square(47), 0)
	if err != nil {
		t.Fatalf("black rook move failed: %v", err)
	}
	if g.Castling&chess.CastleBlackKingSide != 0 {
		t.Error("black kingside castling right should be revoked after h8 rook move")
	}
}

func TestCastlingRevokedOnQueensideRookMove(t *testing.T) {
	g := game.NewGame()

	// White rook a1→a3 (square 0→16)
	g.Board.State[8] = 0 // remove pawn at a2
	err := g.MakeMove(chess.Square(0), chess.Square(16), 0)
	if err != nil {
		t.Fatalf("white rook move failed: %v", err)
	}
	if g.Castling&chess.CastleWhiteQueenSide != 0 {
		t.Error("white queenside castling right should be revoked after a1 rook move")
	}

	// Black rook a8→a6 (square 56→40)
	g.Board.State[48] = 0 // remove pawn at a7
	err = g.MakeMove(chess.Square(56), chess.Square(40), 0)
	if err != nil {
		t.Fatalf("black rook move failed: %v", err)
	}
	if g.Castling&chess.CastleBlackQueenSide != 0 {
		t.Error("black queenside castling right should be revoked after a8 rook move")
	}
}

func TestCastlingRevokedOnRookCaptured(t *testing.T) {
	g := game.NewGame()

	// Setup: place a black bishop where it can capture white's h1 rook
	// Clear path and place black bishop at g2 (square 14)
	g.Board.State[14] = chess.BISHOP | chess.Piece(chess.BLACK)

	// Black bishop captures white rook at h1 (square 7)
	// We need it to be black's turn, so make a white move first
	g.Board.State[12] = 0                                    // remove white pawn at e2
	err := g.MakeMove(chess.Square(12), chess.Square(20), 0) // move placeholder e2→e3 (but e2 is empty, use another)
	if err != nil {
		// Try a simple pawn move instead
		err = g.MakeMove(chess.Square(9), chess.Square(17), 0) // b2→b3
	}

	// Now black captures h1 rook
	err = g.MakeMove(chess.Square(14), chess.Square(7), 0)
	if err != nil {
		t.Fatalf("black capture of white rook failed: %v", err)
	}
	if g.Castling&chess.CastleWhiteKingSide != 0 {
		t.Error("white kingside castling right should be revoked when h1 rook is captured")
	}
}

func TestCastlingPreservedOnIrrelevantMoves(t *testing.T) {
	g := game.NewGame()
	initialCastling := g.Castling

	// White pawn e2→e4 (square 12→28)
	err := g.MakeMove(chess.Square(12), chess.Square(28), 0)
	if err != nil {
		t.Fatalf("white pawn move failed: %v", err)
	}
	if g.Castling != initialCastling {
		t.Errorf("castling rights changed after pawn move: got %v, want %v", g.Castling, initialCastling)
	}

	// Black pawn e7→e5 (square 52→36)
	err = g.MakeMove(chess.Square(52), chess.Square(36), 0)
	if err != nil {
		t.Fatalf("black pawn move failed: %v", err)
	}
	if g.Castling != initialCastling {
		t.Errorf("castling rights changed after pawn move: got %v, want %v", g.Castling, initialCastling)
	}

	// White knight b1→c3 (square 1→18)
	err = g.MakeMove(chess.Square(1), chess.Square(18), 0)
	if err != nil {
		t.Fatalf("white knight move failed: %v", err)
	}
	if g.Castling != initialCastling {
		t.Errorf("castling rights changed after knight move: got %v, want %v", g.Castling, initialCastling)
	}
}
