// Package game is the main body of the chess program
package game

import (
	"fmt"

	"gochess/chess"
)

const startingPositionFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// This is the main state tracking. We'll have a board and its info
// a game state (IN_PLAY, chess.WHITE_WINS, chess.BLACK_WINS, STALEMATE),
// game clock and increment
// move list

// Clock represents the remaining time in the game
// TODO: implement the clock
type Clock struct{}

func (c *Clock) stop(player chess.Color) {
	// do nothing... will implement in the future
}

type Game struct {
	Board        *Board
	Clock        *Clock
	activePlayer chess.Color
	Moves        []chess.Turn
	moveIdx      int
	Castling     chess.Color // Which colors may still castle (e.g. king hasn't moved) starts as 00011000
	Status       chess.GameStatus
}

func NewGame() *Game {
	newGameBoard := LoadFromFEN(startingPositionFen)

	// TODO: implement the clock
	clock := Clock{}
	return &Game{
		Board:        &newGameBoard,
		Clock:        &clock,
		activePlayer: chess.WHITE,
		Moves:        make([]chess.Turn, 0, 8850), // longest possible chess game is ~8,849.5 moves
		Castling:     chess.WHITE | chess.BLACK,
	}
}

func LoadGameFromFen(fen FENCode) *Game {
	board := LoadFromFEN(fen)
	clock := Clock{}
	moves := make([]chess.Turn, 0, 8850)

	return &Game{
		Board:        &board,
		Clock:        &clock,
		activePlayer: board.ActiveColor,
		Moves:        moves,
		moveIdx:      0,
		Castling:     chess.WHITE | chess.BLACK,
	}
}

func (g *Game) MakeMove(from chess.Square, to chess.Square, promo chess.Piece) {
	// stop clock immediately for the current player
	g.Clock.stop(g.activePlayer)

	// TODO: validate the move is legal

	piece := g.Board.State[from]
	ply := chess.Ply{
		Piece:     piece,
		StartIdx:  from,
		EndIdx:    to,
		Promotion: promo,
	}

	// Update the board state
	g.Board.State[to] = piece
	g.Board.State[from] = 0

	// Record the move and switch active player
	switch g.activePlayer {
	case chess.WHITE:
		g.Moves = append(g.Moves, chess.Turn{WhitePly: &ply})
		g.moveIdx = len(g.Moves) - 1
		g.activePlayer = chess.BLACK
	case chess.BLACK:
		g.Moves[g.moveIdx].BlackPly = &ply
		g.activePlayer = chess.WHITE
	}
	g.Board.ActiveColor = g.activePlayer
}

func (g *Game) GetLastMove() *chess.Ply {
	if len(g.Moves) == 0 {
		return nil
	}
	lastTurn := g.Moves[g.moveIdx]
	if lastTurn.BlackPly != nil {
		return lastTurn.BlackPly
	}
	return lastTurn.WhitePly
}

func (g *Game) GetLegalMoves() []chess.Ply {
	switch g.Status {
	case chess.InProgress, chess.Disconnected, chess.Unknown:
		// TODO: define the possible legal moves

		// note: there the theorethical max is 218
		// this is a placeholder until I implement a function
		// to calculate the possible legal moves.
		return make([]chess.Ply, 218)
	}
	return nil
}

func (g *Game) PlayerResigns(player chess.Color) error {
	switch player {
	case chess.WHITE:
		g.Status = chess.BlackWinsForfeit
	case chess.BLACK:
		g.Status = chess.WhiteWinsForfeit
	default:
		//
		return fmt.Errorf("only a player can resign %v is not a player", player)
	}
	// End game gracefully
	g.endGame()
	return nil
}

// endGame gracefully cleans up the game
func (g *Game) endGame() {}
