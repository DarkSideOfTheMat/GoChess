package chess

import "testing"

func TestPieceIsColor(t *testing.T) {
	tests := []struct {
		name  string
		piece Piece
		color Color
		want  bool
	}{
		{"white king is white", KING | Piece(WHITE), WHITE, true},
		{"white king is not black", KING | Piece(WHITE), BLACK, false},
		{"black pawn is black", PAWN | Piece(BLACK), BLACK, true},
		{"black pawn is not white", PAWN | Piece(BLACK), WHITE, false},
		{"colorless piece is neither white", KING, WHITE, false},
		{"colorless piece is neither black", KING, BLACK, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.piece.IsColor(tt.color); got != tt.want {
				t.Errorf("Piece(%d).IsColor(%d) = %v, want %v", tt.piece, tt.color, got, tt.want)
			}
		})
	}
}

func TestPieceToColor(t *testing.T) {
	tests := []struct {
		name  string
		piece Piece
		want  Color
	}{
		{"white queen", QUEEN | Piece(WHITE), WHITE},
		{"black rook", ROOK | Piece(BLACK), BLACK},
		{"colorless piece", BISHOP, Color(0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.piece.ToColor(); got != tt.want {
				t.Errorf("Piece(%d).ToColor() = %v, want %v", tt.piece, got, tt.want)
			}
		})
	}
}

func TestPieceWithColor(t *testing.T) {
	tests := []struct {
		name  string
		piece Piece
		color Color
		want  Piece
	}{
		{"king with white", KING, WHITE, KING | Piece(WHITE)},
		{"pawn with black", PAWN, BLACK, PAWN | Piece(BLACK)},
		{"override color", ROOK | Piece(WHITE), BLACK, ROOK | Piece(BLACK)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.piece.WithColor(tt.color); got != tt.want {
				t.Errorf("Piece(%d).WithColor(%d) = %v, want %v", tt.piece, tt.color, got, tt.want)
			}
		})
	}
}

func TestPieceWithoutColor(t *testing.T) {
	tests := []struct {
		name  string
		piece Piece
		want  Piece
	}{
		{"white knight", KNIGHT | Piece(WHITE), KNIGHT},
		{"black bishop", BISHOP | Piece(BLACK), BISHOP},
		{"already colorless", QUEEN, QUEEN},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.piece.WithoutColor(); got != tt.want {
				t.Errorf("Piece(%d).WithoutColor() = %v, want %v", tt.piece, got, tt.want)
			}
		})
	}
}

func TestPieceIsPieceSameColor(t *testing.T) {
	tests := []struct {
		name  string
		a, b  Piece
		want  bool
	}{
		{"both white", KING | Piece(WHITE), PAWN | Piece(WHITE), true},
		{"both black", QUEEN | Piece(BLACK), ROOK | Piece(BLACK), true},
		{"different colors", KNIGHT | Piece(WHITE), BISHOP | Piece(BLACK), false},
		{"both colorless", KING, PAWN, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.IsPieceSameColor(tt.b); got != tt.want {
				t.Errorf("Piece(%d).IsPieceSameColor(%d) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestPieceToString(t *testing.T) {
	tests := []struct {
		name  string
		piece Piece
		want  string
	}{
		{"king", KING | Piece(WHITE), "K"},
		{"queen", QUEEN | Piece(BLACK), "Q"},
		{"rook", ROOK, "R"},
		{"bishop", BISHOP | Piece(WHITE), "B"},
		{"knight", KNIGHT | Piece(BLACK), "N"},
		{"pawn", PAWN, "P"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.piece.ToString(); got != tt.want {
				t.Errorf("Piece(%d).ToString() = %v, want %v", tt.piece, got, tt.want)
			}
		})
	}
}
