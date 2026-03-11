package game

const startingPositionFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// This is the main state tracking. We'll have a board and its info
// a game state (IN_PLAY, WHITE_WINS, BLACK_WINS, STALEMATE),
// game clock and increment
// move list

// TODO: implement the clock
type Clock struct{}

type Game struct {
	Board    *Board
	Clock    *Clock
	Moves    []Move
	moveIdx  int
	Castling Piece // Which colors may still castle (e.g. king hasn't moved) starts as 00011000
}

func NewGame() *Game {
	new_game_board := LoadFromFEN(startingPositionFen)

	// TODO: implement the clock
	clock := Clock{}
	return &Game{
		//
		Board:    &new_game_board,
		Clock:    &clock,
		Moves:    make([]Move, 0, 8850), // longest possible chess game is ~8,849.5 moves
		Castling: WHITE | BLACK,
	}
}

func LoadGameFromFen(fen FENCode) *Game {
	board := LoadFromFEN(fen)
	clock := Clock{}
	moves := make([]Move, 0, 8850)

	return &Game{
		Board:    &board,
		Clock:    &clock,
		Moves:    moves,
		moveIdx:  0,
		Castling: COLOR_MASK,
	}
}

func (g Game) MakeMove(start_idx int, end_idx int) error {
	return nil
}
