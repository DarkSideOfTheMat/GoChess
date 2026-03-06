package ui

import (
	"fmt"
	game "gochess/game"
	"image/color"

	"charm.land/lipgloss/v2"
)

type boardSettings struct {
	boardStyle      lipgloss.Style
	cellStyle       lipgloss.Style
	lightCellColor  color.Color
	darkCellColor   color.Color
	lightPieceColor color.Color
	darkPieceColor  color.Color
}

type TUISettings struct {
	board boardSettings
}

// darkCell property
func (s TUISettings) getCellStyle(isLightCell bool, piece game.Piece) (lipgloss.Style, error) {
	bg := s.board.darkCellColor
	if isLightCell {
		bg = s.board.lightCellColor
	}

	switch piece.Color() {
	case game.WHITE:
		return lipgloss.NewStyle().Inherit(s.board.cellStyle).Background(bg).Foreground(s.board.lightCellColor), nil

	case game.BLACK:
		return lipgloss.NewStyle().Inherit(s.board.cellStyle).Background(bg).Foreground(s.board.darkCellColor), nil

	default:
		if piece == 0 {
			return lipgloss.NewStyle().Inherit(s.board.cellStyle).Background(bg), nil
		}
		return lipgloss.NewStyle().Inherit(s.board.cellStyle), fmt.Errorf("Cannot format a piece with no color %v", piece)

		// defaut:
		// 	return bs.cellStyle, fmt.Errorf("Cannot format a piece with no color %v", piece)
	}
}

func (s TUISettings) renderBoard(input string) string {
	return s.board.boardStyle.Render(input)
}

var DEFAULT_SETTINGS = TUISettings{
	board: boardSettings{
		boardStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Margin(4).
			BorderForeground(lipgloss.Color("63")),

		cellStyle: lipgloss.NewStyle().
			// Width(3).
			// Height(3).
			Align(lipgloss.Center),
		lightCellColor:  lipgloss.White,
		darkCellColor:   lipgloss.Blue,
		lightPieceColor: lipgloss.BrightWhite,
		darkPieceColor:  lipgloss.Black,
	},
}
