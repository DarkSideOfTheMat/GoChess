package ui

import (
	"bytes"
	"context"
	"testing"
	"time"

	game "gochess/game"

	tea "charm.land/bubbletea/v2"
)

func TestHeadlessRender(t *testing.T) {
	var output bytes.Buffer
	var input bytes.Buffer

	// Send "q" to quit after initial render
	input.WriteString("q")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	m := LoadTeaModelFromFen(game.FENCode("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"))

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
