package ui

import (
	"fmt"

	game "gochess/game"

	"charm.land/lipgloss/v2"
)

var pieceSymbols = map[game.Piece]string{
	game.KING:   "♚",
	game.QUEEN:  "♛",
	game.ROOK:   "♜",
	game.BISHOP: "♝",
	game.KNIGHT: "♞",
	game.PAWN:   "♟",
	0:           " ",
}

var fileLabels = [8]string{"a", "b", "c", "d", "e", "f", "g", "h"}

func getPieceSymbol(piece game.Piece) string {
	if sym, ok := pieceSymbols[piece&^game.COLOR_MASK]; ok {
		return sym
	}
	return " "
}

func (m model) formatCurrentBoard() (string, error) {
	rankLabelStyle := lipgloss.NewStyle().Width(2)
	fileLabelStyle := lipgloss.NewStyle().Width(5).Align(lipgloss.Center)

	var rows []string
	for rank := 7; rank >= 0; rank-- {
		cells := []string{rankLabelStyle.Render(fmt.Sprintf("%d", rank+1))}
		for file := 0; file < 8; file++ {
			piece := m.game.Board.State[rank*8+file]
			style, err := m.settings.getCellStyle((rank+file)%2 != 0, piece)
			if err != nil {
				return "", err
			}
			cells = append(cells, style.Render(getPieceSymbol(piece)))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Center, cells...))
	}
	fileRow := []string{"  "}
	for _, label := range fileLabels {
		fileRow = append(fileRow, fileLabelStyle.Render(label))
	}
	rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Center, fileRow...))
	return m.settings.renderBoard(lipgloss.JoinVertical(lipgloss.Left, rows...)), nil
}
