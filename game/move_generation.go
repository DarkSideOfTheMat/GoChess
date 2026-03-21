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
	moves := make([]chess.Ply, 0, 64)
	return moves
}

func (mg IterativePsuedoLegalMoveGenerator) GenBishopMoves(i chess.Square, p chess.Piece) []chess.Ply {
	moves := make([]chess.Ply, 0, 64)
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
		if i+15 < 64 {
			moves = append(moves, mkPly(15))
		}
		// left 1, down 2
		if i-17 >= 0 {
			moves = append(moves, mkPly(-17))
		}
	}
	// not B
	if file > 1 {
		// left 2, up 1
		if i+6 < 64 {
			moves = append(moves, mkPly(6))
		}
		// left 2, down 1
		if i-10 >= 0 {
			moves = append(moves, mkPly(-10))
		}
	}
	// not G
	if file < 6 {
		// right 2 up 1
		if i+10 < 64 {
			moves = append(moves, mkPly(10))
		}
		// right 2 down 1
		if i-6 >= 0 {
			moves = append(moves, mkPly(-6))
		}
	}
	// not H
	if file < 7 {
		// right 1 up 2
		if i+17 < 64 {
			moves = append(moves, mkPly(17))
		}
		// right 1 down 2
		if i-15 >= 0 {
			moves = append(moves, mkPly(-15))
		}
	}
	return moves
}
