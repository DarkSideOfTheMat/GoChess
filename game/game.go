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
	Castling     chess.Castling // Which colors may still castle (e.g. king hasn't moved) starts as 00011000
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
		Castling:     chess.CastleWhiteMask | chess.CastleBlackMask,
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
		Castling:     chess.CastleWhiteMask | chess.CastleBlackMask,
	}
}

func (g *Game) MakeMove(from chess.Square, to chess.Square, promo chess.Piece) error {
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

	if g.Board.State[to].IsColor(g.activePlayer) {
		return fmt.Errorf("cannot capture own piece")
	}
	if piece == chess.Piece(0) || piece == 0 {
		return fmt.Errorf("no piece on square %s", from.ToString())
	}
	// Check the selected piece is the current players color
	if !piece.IsColor(g.activePlayer) {
		return fmt.Errorf("cannot move another player's piece on %s%s", piece.ToString(), from.ToString())
	}

	// Check if the move is a castling move
	isCastling, side, err := IsCastlingMove(piece, from, to)
	if err != nil {
		return err
	}

	if isCastling {
		// check that the player has the right to castle

		// check the conditions for castling are met (no blocking pieces or attacks)

		// Move the rook to new square, allow regular move handling below
		switch side.WithoutColor() {
		case chess.CastleKingSide:
			rookTo, rookFrom := to-1, to+1
			g.Board.State[rookTo] = g.Board.State[rookFrom]
			g.Board.State[rookFrom] = 0

			// flip off all castling for player
			g.Castling &^= chess.CastleMask.WithColor(g.activePlayer)
		case chess.CastleQueenSide:
			rookTo, rookFrom := to+1, to-2
			g.Board.State[rookTo] = g.Board.State[rookFrom]
			g.Board.State[rookFrom] = 0

			// flip off all castling for player
			g.Castling &^= chess.CastleMask.WithColor(g.activePlayer)
		}
	}
	// update castling rights
	if piece.WithoutColor() == chess.KING {
		g.Castling &^= chess.CastleMask.WithColor(g.activePlayer)
	}
	// g.activePlayer / 16 = 1 black, 0 for white
	// rook is moving off queenside home square
	if piece == chess.ROOK && from == chess.Square(0+8*(int(g.activePlayer)/16)) {
		side, _ = CastleFromColorAndSide(g.activePlayer, chess.CastleQueenSide)
		g.Castling &^= side
		// rook is moving off kingside home square
	} else if piece == chess.ROOK && from == chess.Square(7+8*(int(g.activePlayer)/16)) {
		side, _ = CastleFromColorAndSide(g.activePlayer, chess.CastleKingSide)
		g.Castling &^= side
	}
	// check if taking rook on the opposite side home square

	// Update the board state
	g.Board.State[to] = piece
	g.Board.State[from] = 0

	// Record the move and switch active player
	switch g.activePlayer {
	case chess.WHITE:
		// record move
		g.Moves = append(g.Moves, chess.Turn{WhitePly: &ply})
		g.moveIdx = len(g.Moves) - 1
		g.activePlayer = chess.BLACK
	case chess.BLACK:
		// record move
		g.Moves[g.moveIdx].BlackPly = &ply
		g.activePlayer = chess.WHITE
	}
	g.Board.ActiveColor = g.activePlayer
	return nil
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
