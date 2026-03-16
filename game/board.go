package game

import (
	"strings"

	"gochess/chess"
)

// FEN strings represent a chess position
// describing each rank starting with the last rank separated by /.
// PNBRQK for white pieces and pnbrqk for black
//
// example: 1B6/2n5/p1N1P2R/P1K3N1/4Pk2/1Q2p2p/6nP/1B4R1 w - - 0 1
type FENCode string

// LoadFromFEN populates a Board from a FEN string.
func LoadFromFEN(fen FENCode) Board {
	var board_state BoardState
	fields := strings.SplitN(string(fen), " ", 3)

	ranks := strings.Split(fields[0], "/")
	for rankIdx, rankStr := range ranks {
		fileIdx := 0
		for _, ch := range rankStr {
			if ch >= '1' && ch <= '8' {
				fileIdx += int(ch - '0')
				continue
			}
			board_state[(7-rankIdx)*8+fileIdx] = fenCharToPiece(ch)
			fileIdx++
		}
	}

	active_color := chess.WHITE // may need a better way to
	if len(fields) >= 2 && fields[1] == "b" {
		active_color = chess.BLACK
	}

	return Board{board_state, active_color}
}

func fenCharToPiece(ch rune) chess.Piece {
	color := chess.WHITE
	if ch >= 'a' && ch <= 'z' {
		color = chess.BLACK
		ch -= 32
	}
	switch ch {
	case 'K':
		return chess.Piece(uint8(color) | uint8(chess.KING))
	case 'Q':
		return chess.Piece(uint8(color) | uint8(chess.QUEEN))
	case 'R':
		return chess.Piece(uint8(color) | uint8(chess.ROOK))
	case 'B':
		return chess.Piece(uint8(color) | uint8(chess.BISHOP))
	case 'N':
		return chess.Piece(uint8(color) | uint8(chess.KNIGHT))
	case 'P':
		return chess.Piece(uint8(color) | uint8(chess.PAWN))
	}
	return 0
}

// the board state will be a single array, the indexes will look like:
//
// -+ ------------------------------ +
// 8| 56, 57, 58, 59, 60, 61, 62, 63 |
// 7| 48, 49, 50, 51, 52, 53, 54, 55,|
// 6| 40, 41, 42, 43, 44, 45, 46, 47,|
// 5| 32, 33, 34, 35, 36, 37, 38, 39,|
// 4| 24, 25, 26, 27, 28, 29, 30, 31,|
// 3| 16, 17, 18, 19, 20, 21, 22, 23,|
// 2|  8,  9, 10, 11, 12, 13, 14, 15,|
// 1|  0,  1,  2,  3,  4,  5,  6,  7,|
// -+ -------------------------------+
// .   a   b   c   d   e   f   g   h
type BoardState [64]chess.Piece

// Board represents the current state of a game of chess.
// It holds the piece positions and the active player's color.
type Board struct {
	State       BoardState
	ActiveColor chess.Color // chess.WHITE or chess.BLACK, also can be last moved piece
}

// Validate returns true if every non-empty square has a valid piece type and color.
func (b *Board) Validate() bool {
	const typeMask chess.Piece = 0b00000111
	const colorMask chess.Piece = 0b00011000
	for _, p := range b.State {
		if p == 0 {
			continue
		}
		t := p & typeMask
		c := p.ToColor()
		if t < 1 || t > 6 {
			return false
		}
		if c != chess.WHITE && c != chess.BLACK {
			return false
		}
	}
	return true
}

// BitBoards are a representation of the current board state...
// each bit represents a position on the board corresponding to the
// matching tile...
//
// # Collections of bitboards can be used to perform calculations
//
// -+ ------------------------------ +
// 8| 56, 57, 58, 59, 60, 61, 62, 63 |
// 7| 48, 49, 50, 51, 52, 53, 54, 55,|
// 6| 40, 41, 42, 43, 44, 45, 46, 47,|
// 5| 32, 33, 34, 35, 36, 37, 38, 39,|
// 4| 24, 25, 26, 27, 28, 29, 30, 31,|
// 3| 16, 17, 18, 19, 20, 21, 22, 23,|
// 2|  8,  9, 10, 11, 12, 13, 14, 15,|
// 1|  0,  1,  2,  3,  4,  5,  6,  7,|
// -+ -------------------------------+
// .   a   b   c   d   e   f   g   h
type BitBoard uint64

func ValidKingMoves(startIdx int) BitBoard {
	var b uint64
	var i uint64 = 1 << startIdx

	notA := uint64(0xFEFEFEFEFEFEFEFE)
	notH := uint64(0x7F7F7F7F7F7F7F7F)

	b |= i << 8        // 1 up
	b |= i >> 8        // 1 down
	b |= i & notH << 1 // 1 right
	b |= i & notA >> 1 // 1 left

	return BitBoard(b)
}

// Return a BitBoard of the possible valid knights moves
// starting on the index...
// note: these are just the possbilities, doesn't respect pins or checks
func ValidKnightsMoves(startIdx int) BitBoard {
	var b uint64
	var i uint64 = 1 << startIdx

	// hex cheatsheet
	// 0 0000  1 0001
	// 2 0010  3 0011
	// 4 0100  5 0101
	// 6 0110  7 0111
	// 8 1000  9 1001
	// A 1010  B 1011
	// C 1100  D 1101
	// E 1110  F 1111
	notA := uint64(0xFEFEFEFEFEFEFEFE)
	notB := uint64(0xFDFDFDFDFDFDFDFD)
	notG := uint64(0xBFBFBFBFBFBFBFBF)
	notH := uint64(0x7F7F7F7F7F7F7F7F)
	b |= (i & notA) << 15        // +2 ranks, -1 file
	b |= (i & notH) << 17        // +2 ranks, +1 file
	b |= (i & notA & notB) << 6  // +1 rank, -2 files
	b |= (i & notH & notG) << 10 // +1 rank, +2 files

	b |= (i & notH) >> 15        // -2 ranks, +1 file
	b |= (i & notA) >> 17        // -2 ranks, -1 file
	b |= (i & notH & notG) >> 6  // -1 rank, +2 files
	b |= (i & notA & notB) >> 10 // -1 rank, -2 files

	return BitBoard(b)
}

// Return a BitBoard of sliding Rook moves
// starting on the index
// note: these don't respect pins, checks or blocking pieces
func ValidRookMoves(startIdx int) BitBoard {
	var b int64
	// TODO implement this
	// fileA := uint64(0x0101010101010101)
	return BitBoard(b)
}
