package game

// This is the main state tracking. We'll have a board and its info
// a game state (IN_PLAY, WHITE_WINS, BLACK_WINS, STALEMATE),
// game clock and increment
// move list

type Game struct {
	board        *Board
	player_clock map[Piece]uint32
	increment    uint32
	moves        []Move
	castling     Piece // Which colors may still castle (e.g. king hasn't moved) starts as 00011000
}
