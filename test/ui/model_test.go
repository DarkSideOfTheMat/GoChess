package ui_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"gochess/chess"
	protocol "gochess/game/protocol"
	session "gochess/game/session"
	"gochess/ui"

	tea "charm.land/bubbletea/v2"
)

func TestHeadlessRender(t *testing.T) {
	var output bytes.Buffer
	var input bytes.Buffer

	input.WriteString("\x03") // ctrl+c

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sess := session.NewGameSession()
	m := ui.NewModel(&sess)

	p := tea.NewProgram(m,
		tea.WithContext(ctx),
		tea.WithInput(&input),
		tea.WithOutput(&output),
		tea.WithWindowSize(80, 40),
		tea.WithoutSignals(),
	)

	if _, err := p.Run(); err != nil {
		t.Fatalf("program error: %v", err)
	}

	rendered := output.String()
	if len(rendered) == 0 {
		t.Fatal("no output rendered")
	}

	t.Logf("Rendered output:\n%s", rendered)
}

func TestMoveE2E4ViaSession(t *testing.T) {
	sess := session.NewGameSession()
	ch := sess.Subscribe()

	// Receive initial state
	initialEvent := <-ch
	initial, ok := initialEvent.(protocol.GameStateEvent)
	if !ok {
		t.Fatalf("expected GameStateEvent, got %T", initialEvent)
	}
	e2 := chess.Square(12) // e2 = rank 1, file 4 = 1*8+4 = 12
	e4 := chess.Square(28) // e4 = rank 3, file 4 = 3*8+4 = 28

	// Verify pawn is on e2, e4 is empty
	if initial.Board.State[e2] == 0 {
		t.Fatal("expected a piece on e2")
	}
	if initial.Board.State[e4] != 0 {
		t.Fatalf("expected e4 to be empty, got %v", initial.Board.State[e4])
	}

	pawn := initial.Board.State[e2]
	t.Logf("Piece on e2: %d (type=%d, color=%d)", pawn, pawn&chess.PieceMask, pawn.ToColor())

	// Send e2 e4
	go func() {
		sess.Send(protocol.MoveMessage{From: e2, To: e4})
	}()

	// Receive updated state
	updatedEvent := <-ch
	updated, ok := updatedEvent.(protocol.GameStateEvent)
	if !ok {
		t.Fatalf("expected GameStateEvent after move, got %T", updatedEvent)
	}

	// Verify pawn moved
	if updated.Board.State[e2] != 0 {
		t.Errorf("e2 should be empty after move, got %v", updated.Board.State[e2])
	}
	if updated.Board.State[e4] != pawn {
		t.Errorf("e4 should have the pawn (%v), got %v", pawn, updated.Board.State[e4])
	}

	t.Logf("Move e2 e4 successful: pawn moved from e2 to e4")
}

func TestParseSquareViaChessPackage(t *testing.T) {
	tests := []struct {
		input    string
		expected chess.Square
		wantErr  bool
	}{
		{"a1", 0, false},
		{"e2", 12, false},
		{"e4", 28, false},
		{"h8", 63, false},
		{"z9", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		sq, err := chess.ParseSquare(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("ParseSquare(%q) expected error, got %d", tt.input, sq)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseSquare(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if sq != tt.expected {
			t.Errorf("ParseSquare(%q) = %d, want %d", tt.input, sq, tt.expected)
		}
	}
}
