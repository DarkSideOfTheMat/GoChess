package game

import "fmt"

// A chess ply is half a move
type Ply struct {
	piece          Piece // piece that moved, including player color
	start_idx      int
	end_idx        int
	promotion_type Piece // when a pawn promotes we have to denote what it promoted to
}

func (p Ply) Validate() error {
	return nil
}

type Move struct {
	white_ply *Ply
	black_ply *Ply
}

func (m Move) Validate() error {
	if m.white_ply.piece.IsColor(BLACK) || m.black_ply.piece.IsColor(WHITE) {
		return fmt.Errorf("Invalid move! Black cannot move on White's turn")
	}
	if m.black_ply.piece.IsColor(WHITE) {
		return fmt.Errorf("Invalid move! White cannot move on Black's turn!")
	}
	return nil
}
