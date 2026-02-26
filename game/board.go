package game

import (
	"fmt"
	"strings"
)

// Boardstate

// TODO: create a representation of the board with piece positions
// TODO: add a renderer interface, implement the interface for printing positions in the terminal

// the board will be a single array, the indexes will look like:
//
//   - ------------------------------ +
//
// 8| 56, 57, 58, 59, 60, 61, 62, 63 |
// 7| 48, 49, 50, 51, 52, 53, 54, 55,|
// 6| 40, 41, 42, 43, 44, 45, 46, 47,|
// 5| 32, 33, 34, 35, 36, 37, 38, 39,|
// 4| 24, 25, 26, 27, 28, 29, 30, 31,|
// 3| 16, 17, 18, 19, 20, 21, 22, 23,|
// 2|  8,  9, 10, 11, 12, 13, 14, 15,|
// 1|  0,  1,  2,  3,  4,  5,  6,  7,|
//   - -------------------------------+
//
// .    a   b   c   d   e   f   g   h
// ]
type BoardState [64]Piece

// Board represents the current state of a game of chess
// it has a board state which has pieces and positions
// it has the current player turn
type Board interface {
	state() BoardState
}

// FEN strings represent a chess position
// describing each rank starting with the last rank separated by /.
// PNBRQK for white pieces and pnbrqk for black
//
// example: 1B6/2n5/p1N1P2R/P1K3N1/4Pk2/1Q2p2p/6nP/1B4R1 w - - 0 1
type FENCode string

// LoadFromFEN populates a BoardState from a FEN string.
func LoadFromFEN(fen FENCode, board_state *BoardState) {
	fields := strings.SplitN(string(fen), " ", 2)
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
}

func fenCharToPiece(ch rune) Piece {
	color := WHITE
	if ch >= 'a' && ch <= 'z' {
		color = BLACK
		ch -= 32
	}
	switch ch {
	case 'K':
		return color | KING
	case 'Q':
		return color | QUEEN
	case 'R':
		return color | ROOK
	case 'B':
		return color | BISHOP
	case 'N':
		return color | KNIGHT
	case 'P':
		return color | PAWN
	}
	return 0
}

// Validate returns true if every non-empty square has a valid piece type and color.
func (b *BoardState) Validate() bool {
	const typeMask Piece = 0b00000111
	const colorMask Piece = 0b00011000
	for _, p := range b {
		if p == 0 {
			continue
		}
		t := p & typeMask
		c := p & colorMask
		if t < 1 || t > 6 {
			return false
		}
		if c != WHITE && c != BLACK {
			return false
		}
	}
	return true
}

var pieceSymbols = map[Piece]string{
	WHITE | KING: "♔", WHITE | QUEEN: "♕", WHITE | ROOK: "♖",
	WHITE | BISHOP: "♗", WHITE | KNIGHT: "♘", WHITE | PAWN: "♙",
	BLACK | KING: "♚", BLACK | QUEEN: "♛", BLACK | ROOK: "♜",
	BLACK | BISHOP: "♝", BLACK | KNIGHT: "♞", BLACK | PAWN: "♟",
}

// Render prints the board state to stdout.
func (b *BoardState) Render() {

	for rank := 7; rank >= 0; rank-- {
		fmt.Printf("%d ", rank+1)
		for file := range 8 {
			p := b[rank*8+file]
			if sym, ok := pieceSymbols[p]; ok {
				fmt.Printf("%s ", sym)
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Println("  a b c d e f g h")
}
