package game_test

import (
	"testing"

	"gochess/chess"
	"gochess/game"
)

// These tests verify that Game.Castling state is correctly updated when
// MakeMove is called. They will FAIL until castling revocation logic is
// added to MakeMove.

// requireMove calls MakeMove and verifies the piece moved correctly:
// no error, piece is at destination, source square is empty.
func requireMove(t *testing.T, g *game.Game, from, to chess.Square) {
	t.Helper()
	piece := g.Board.State[from]
	if piece == 0 {
		t.Fatalf("no piece at source square %s before move", from.ToString())
	}
	err := g.MakeMove(from, to, 0)
	if err != nil {
		t.Fatalf("MakeMove(%s→%s) returned error: %v", from.ToString(), to.ToString(), err)
	}
	if g.Board.State[to] != piece {
		t.Fatalf("after MakeMove(%s→%s): expected piece %d at destination, got %d",
			from.ToString(), to.ToString(), piece, g.Board.State[to])
	}
	if g.Board.State[from] != 0 {
		t.Fatalf("after MakeMove(%s→%s): source square should be empty, got %d",
			from.ToString(), to.ToString(), g.Board.State[from])
	}
}

func TestCastlingRevokedOnKingMove(t *testing.T) {
	g := game.NewGame()

	// Clear the path for kings: move pawns out of the way
	requireMove(t, g, 12, 20) // e2→e3
	requireMove(t, g, 52, 36) // e7→e5

	// White king e1→e2
	requireMove(t, g, 4, 12)
	if g.Castling&chess.CastleWhiteMask != 0 {
		t.Error("white castling rights should be cleared after king move")
	}
	if g.Castling&chess.CastleBlackMask == 0 {
		t.Error("black castling rights should be preserved after white king move")
	}

	// Black king e8→e7
	requireMove(t, g, 60, 52)
	if g.Castling&chess.CastleBlackMask != 0 {
		t.Error("black castling rights should be cleared after king move")
	}
}

func TestCastlingRevokedOnKingsideRookMove(t *testing.T) {
	g := game.NewGame()

	// Clear pawns blocking the rooks
	g.Board.State[15] = 0 // remove pawn at h2
	requireMove(t, g, 7, 23) // white rook h1→h3
	if g.Castling&chess.CastleWhiteKingSide != 0 {
		t.Error("white kingside castling right should be revoked after h1 rook move")
	}
	if g.Castling&chess.CastleWhiteQueenSide == 0 {
		t.Error("white queenside castling right should be preserved after h1 rook move")
	}

	g.Board.State[55] = 0 // remove pawn at h7
	requireMove(t, g, 63, 47) // black rook h8→h6
	if g.Castling&chess.CastleBlackKingSide != 0 {
		t.Error("black kingside castling right should be revoked after h8 rook move")
	}
	if g.Castling&chess.CastleBlackQueenSide == 0 {
		t.Error("black queenside castling right should be preserved after h8 rook move")
	}
}

func TestCastlingRevokedOnQueensideRookMove(t *testing.T) {
	g := game.NewGame()

	g.Board.State[8] = 0 // remove pawn at a2
	requireMove(t, g, 0, 16) // white rook a1→a3
	if g.Castling&chess.CastleWhiteQueenSide != 0 {
		t.Error("white queenside castling right should be revoked after a1 rook move")
	}
	if g.Castling&chess.CastleWhiteKingSide == 0 {
		t.Error("white kingside castling right should be preserved after a1 rook move")
	}

	g.Board.State[48] = 0 // remove pawn at a7
	requireMove(t, g, 56, 40) // black rook a8→a6
	if g.Castling&chess.CastleBlackQueenSide != 0 {
		t.Error("black queenside castling right should be revoked after a8 rook move")
	}
	if g.Castling&chess.CastleBlackKingSide == 0 {
		t.Error("black kingside castling right should be preserved after a8 rook move")
	}
}

func TestCastlingRevokedOnRookCaptured(t *testing.T) {
	g := game.NewGame()

	// Setup: place a black bishop at g2 so it can capture white's h1 rook
	g.Board.State[14] = chess.BISHOP | chess.Piece(chess.BLACK) // g2
	g.Board.State[15] = 0                                       // remove white pawn at h2

	// White makes a pawn move to pass the turn
	requireMove(t, g, 9, 17) // b2→b3

	// Black bishop captures white rook at h1
	blackBishop := g.Board.State[14]
	if blackBishop == 0 {
		t.Fatal("expected black bishop at g2 (square 14)")
	}
	err := g.MakeMove(14, 7, 0)
	if err != nil {
		t.Fatalf("black capture of white rook failed: %v", err)
	}
	if g.Board.State[7] != blackBishop {
		t.Fatalf("black bishop should be at h1 after capture, got %d", g.Board.State[7])
	}
	if g.Board.State[14] != 0 {
		t.Fatal("g2 should be empty after bishop moved")
	}
	if g.Castling&chess.CastleWhiteKingSide != 0 {
		t.Error("white kingside castling right should be revoked when h1 rook is captured")
	}
}

func TestCastlingPreservedOnIrrelevantMoves(t *testing.T) {
	g := game.NewGame()
	initialCastling := g.Castling

	requireMove(t, g, 12, 28) // white pawn e2→e4
	if g.Castling != initialCastling {
		t.Errorf("castling rights changed after pawn move: got %d, want %d", g.Castling, initialCastling)
	}

	requireMove(t, g, 52, 36) // black pawn e7→e5
	if g.Castling != initialCastling {
		t.Errorf("castling rights changed after pawn move: got %d, want %d", g.Castling, initialCastling)
	}

	requireMove(t, g, 1, 18) // white knight b1→c3
	if g.Castling != initialCastling {
		t.Errorf("castling rights changed after knight move: got %d, want %d", g.Castling, initialCastling)
	}
}
