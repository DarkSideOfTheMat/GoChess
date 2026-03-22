package game

import (
	"slices"

	"gochess/chess"
)

const MaxLegalMoves = 218

type MoveGenerator interface {
	GenMoves() []chess.Ply
	GetCurrentBoard() Board
	GetBoardFromMove(m chess.Ply) Board
}

type IterativePsuedoLegalMoveGenerator struct {
	board Board
}

func (mg IterativePsuedoLegalMoveGenerator) GenLegalMoves() []chess.Ply {
	moves := make([]chess.Ply, 0, MaxLegalMoves)
	player := mg.board.ActiveColor
	otherPlayer := mg.board.ActiveColor.Flip()
	var piece chess.Piece

	for i := range chess.Square(64) {
		piece = mg.board.State[i]
		if piece == 0 || piece.IsColor(otherPlayer) {
			continue
		}

		switch piece.WithoutColor() {
		case chess.KING:
			moves = append(moves, mg.GenKingMoves(i, piece)...)
		case chess.QUEEN:
			moves = append(moves, mg.GenQueenMoves(i, piece)...)
		case chess.ROOK:
			moves = append(moves, mg.GenRookMoves(i, piece)...)
		case chess.BISHOP:
			moves = append(moves, mg.GenBishopMoves(i, piece)...)
		case chess.KNIGHT:
			moves = append(moves, mg.GenKnightMoves(i, piece)...)
		case chess.PAWN:
			moves = append(moves, mg.GenPawnMoves(i, piece)...)
		}
	}
	// en-passant
	// this is a fake move so we won't implement it.
	// Only cheaters with access to "Google" ever use it

	// prune moves

	// prune same-side captures
	moves = slices.DeleteFunc(moves, func(p chess.Ply) bool { return mg.board.State[p.EndIdx].IsColor(player) })
	return moves
}

func (mg IterativePsuedoLegalMoveGenerator) GenKingMoves(i chess.Square, p chess.Piece) []chess.Ply {
	// a king can move at most 9 times
	moves := make([]chess.Ply, 0, 9)
	return moves
}

func (mg IterativePsuedoLegalMoveGenerator) GenQueenMoves(i chess.Square, p chess.Piece) []chess.Ply {
	// a queen is just rook + bishop
	return slices.Concat(mg.GenRookMoves(i, p), mg.GenBishopMoves(i, p))
}

func (mg IterativePsuedoLegalMoveGenerator) GenRookMoves(i chess.Square, p chess.Piece) []chess.Ply {
	moves := make([]chess.Ply, 0, 16)

	player := mg.board.ActiveColor
	otherPlayer := mg.board.ActiveColor.Flip()

	// left and right
	for l, r := i-1, i+1; l >= 0 && l.Rank() == i.Rank() || r.Rank() == i.Rank(); l, r = l-1, r+1 {
		if l.Rank() == i.Rank() && l >= 0 {
			if mg.board.State[l].IsColor(player) {
				// path blocked
				l = l - 8
			} else if mg.board.State[l].IsColor(otherPlayer) {
				moves = append(
					moves,
					chess.Ply{
						Piece:     p,
						StartIdx:  i,
						EndIdx:    l,
						Promotion: 0,
					},
				)
				l = l - 8
			} else {
				moves = append(
					moves,
					chess.Ply{
						Piece:     p,
						StartIdx:  i,
						EndIdx:    l,
						Promotion: 0,
					},
				)
			}
		}
		if r.Rank() == i.Rank() && r < 64 {
			if mg.board.State[r].IsColor(player) {
				r = r + 8
			} else if mg.board.State[r].IsColor(otherPlayer) {
				moves = append(
					moves,
					chess.Ply{
						Piece:     p,
						StartIdx:  i,
						EndIdx:    r,
						Promotion: 0,
					},
				)
				r = r + 8
			} else {
				moves = append(
					moves,
					chess.Ply{
						Piece:     p,
						StartIdx:  i,
						EndIdx:    r,
						Promotion: 0,
					},
				)
			}
		}
	}
	// up & down
	for d, u := i-8, i+8; d >= 0 || u < 64; d, u = d-8, u+8 {
		if d >= 0 {
			if mg.board.State[d].IsColor(otherPlayer) {
				d = -1
			}
		}
		if u < 64 {
			if mg.board.State[u].IsColor(otherPlayer) {
				d = 64
			} else if mg.board.State[u].IsColor(player) {
				moves = append(
					moves,
					chess.Ply{
						Piece:     p,
						StartIdx:  i,
						EndIdx:    u,
						Promotion: 0,
					},
				)
				d = 64
			} else {
				moves = append(
					moves,
					chess.Ply{
						Piece:     p,
						StartIdx:  i,
						EndIdx:    u,
						Promotion: 0,
					},
				)
			}
		}
	}

	//
	return moves
}

func (mg IterativePsuedoLegalMoveGenerator) GenBishopMoves(i chess.Square, p chess.Piece) []chess.Ply {
	moves := make([]chess.Ply, 0, 16)
	player := mg.board.ActiveColor
	otherPlayer := player.Flip()
	// left... stops on A-file
	for u, d := i+7, i-9; u < 64 || d >= 0; u, d = u+7, d-9 {
		if u < 64 {
			if mg.board.State[u].IsColor(player) {
				u = 64 // break
			} else if mg.board.State[u].IsColor(otherPlayer) {
				moves = append(moves, chess.Ply{Piece: p, StartIdx: i, EndIdx: u, Promotion: 0})
				u = 64 // break
			} else {
				moves = append(moves, chess.Ply{Piece: p, StartIdx: i, EndIdx: u, Promotion: 0})
				if u.File() == 0 {
					u = 64 // break
				}
			}
		}
		if d >= 0 {
			if mg.board.State[d].IsColor(player) {
				d = -1 // break
			} else if mg.board.State[d].IsColor(otherPlayer) {
				moves = append(moves, chess.Ply{Piece: p, StartIdx: i, EndIdx: d, Promotion: 0})
				d = -1 // break
			} else {
				moves = append(moves, chess.Ply{Piece: p, StartIdx: i, EndIdx: d, Promotion: 0})
				if d.File() == 0 {
					d = -1 // break
				}
			}
		}
	}

	// right... stops on H-File
	for u, d := i+9, i-7; u < 64 || d >= 0; u, d = u+9, d-7 {
		if u < 64 {
			if mg.board.State[u].IsColor(player) {
				u = 64 // break
			} else if mg.board.State[u].IsColor(otherPlayer) {
				moves = append(moves, chess.Ply{Piece: p, StartIdx: i, EndIdx: u, Promotion: 0})
				u = 64 // break
			} else {
				moves = append(moves, chess.Ply{Piece: p, StartIdx: i, EndIdx: u, Promotion: 0})
				if u.File() == 7 {
					u = 64 // reached H-file, stop
				}
			}
		}
		if d >= 0 {
			if mg.board.State[d].IsColor(player) {
				d = -1 // break
			} else if mg.board.State[d].IsColor(otherPlayer) {
				moves = append(moves, chess.Ply{Piece: p, StartIdx: i, EndIdx: d, Promotion: 0})
				d = -1 // break
			} else {
				moves = append(moves, chess.Ply{Piece: p, StartIdx: i, EndIdx: d, Promotion: 0})
				if d.File() == 7 {
					d = -1 // reached H-file, stop
				}
			}
		}
	}

	return moves
}

func (mg IterativePsuedoLegalMoveGenerator) GenPawnMoves(i chess.Square, p chess.Piece) []chess.Ply {
	moves := make([]chess.Ply, 0, 9)
	file := i % 8
	rank := i / 8

	otherPlayer := p.ToColor().Flip()

	var promoRank chess.Square = 7    // rank "8"
	var startingRank chess.Square = 1 // rank "2"
	var moveDirection chess.Square = 1

	var leftAttack chess.Square = 7
	var rightAttack chess.Square = 9

	if p.IsColor(chess.BLACK) {
		promoRank = 0    // rank "1"
		startingRank = 6 // rank "7"
		moveDirection = -1
		leftAttack, rightAttack = rightAttack*moveDirection, leftAttack*moveDirection
	}

	addPly := func(moves []chess.Ply, offset chess.Square, isPromo bool) []chess.Ply {
		if isPromo {
			moves = append(
				moves,
				chess.Ply{Piece: p, StartIdx: i, EndIdx: i + offset, Promotion: chess.QUEEN},
				chess.Ply{Piece: p, StartIdx: i, EndIdx: i + offset, Promotion: chess.ROOK},
				chess.Ply{Piece: p, StartIdx: i, EndIdx: i + offset, Promotion: chess.KNIGHT},
			)
		} else {
			moves = append(
				moves,
				chess.Ply{
					Piece:     p,
					StartIdx:  i,
					EndIdx:    i + offset,
					Promotion: 0,
				},
			)
		}

		return moves
	}

	// up 1
	if mg.board.State[i+8*moveDirection] == 0 {
		moves = addPly(moves, 8*moveDirection, rank == promoRank-moveDirection)

		// up 2, can't skip pieces
		if rank == startingRank && mg.board.State[i+16*moveDirection] == 0 {
			moves = addPly(moves, 16*moveDirection, false)
		}
	}

	// attacks
	// left
	if file > 0 && mg.board.State[i+leftAttack].IsColor(otherPlayer) {
		moves = addPly(moves, leftAttack, rank == promoRank-moveDirection)
	}
	// right
	if file < 7 && mg.board.State[i+rightAttack].IsColor(otherPlayer) {
		moves = addPly(moves, rightAttack, rank == promoRank-moveDirection)
	}

	return moves[:]
}

func (mg IterativePsuedoLegalMoveGenerator) GenKnightMoves(i chess.Square, p chess.Piece) []chess.Ply {
	moves := make([]chess.Ply, 0, 8)
	file := i % 8
	player := p.ToColor()

	mkPly := func(offset chess.Square) chess.Ply {
		return chess.Ply{
			Piece:     p,
			StartIdx:  i,
			EndIdx:    i + offset,
			Promotion: 0,
		}
	}

	// not A
	if file > 0 {
		// left 1, up 2
		if i+15 < 64 && !mg.board.State[i+15].IsColor(player) {
			moves = append(moves, mkPly(15))
		}
		// left 1, down 2
		if i-17 >= 0 && !mg.board.State[i-17].IsColor(player) {
			moves = append(moves, mkPly(-17))
		}
	}
	// not B
	if file > 1 {
		// left 2, up 1
		if i+6 < 64 && !mg.board.State[i+6].IsColor(player) {
			moves = append(moves, mkPly(6))
		}
		// left 2, down 1
		if i-10 >= 0 && !mg.board.State[i-10].IsColor(player) {
			moves = append(moves, mkPly(-10))
		}
	}
	// not G
	if file < 6 {
		// right 2 up 1
		if i+10 < 64 && !mg.board.State[i+10].IsColor(player) {
			moves = append(moves, mkPly(10))
		}
		// right 2 down 1
		if i-6 >= 0 && !mg.board.State[i-6].IsColor(player) {
			moves = append(moves, mkPly(-6))
		}
	}
	// not H
	if file < 7 {
		// right 1 up 2
		if i+17 < 64 && !mg.board.State[i+17].IsColor(player) {
			moves = append(moves, mkPly(17))
		}
		// right 1 down 2
		if i-15 >= 0 && !mg.board.State[i-15].IsColor(player) {
			moves = append(moves, mkPly(-15))
		}
	}
	return moves
}
