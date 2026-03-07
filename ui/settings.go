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
		return lipgloss.NewStyle().Inherit(s.board.cellStyle).Background(bg).Foreground(s.board.lightPieceColor), nil

	case game.BLACK:
		return lipgloss.NewStyle().Inherit(s.board.cellStyle).Background(bg).Foreground(s.board.darkPieceColor), nil

	default:
		if piece == 0 {
			return lipgloss.NewStyle().Inherit(s.board.cellStyle).Background(bg), nil
		}
		return lipgloss.NewStyle().Inherit(s.board.cellStyle), fmt.Errorf("Cannot format a piece with no color %v", piece)
	}
}

func (s TUISettings) renderBoard(input string) string {
	return s.board.boardStyle.Render(input)
}

var headerStyle = lipgloss.NewStyle().
	Bold(true).
	Padding(0, 1)

var footerStyle = lipgloss.NewStyle().
	Padding(0, 1)

var DEFAULT_SETTINGS = TUISettings{
	board: boardSettings{
		boardStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Margin(1).
			BorderForeground(lipgloss.Color("63")),

		cellStyle: lipgloss.NewStyle().
			Width(5).
			Height(1).
			Align(lipgloss.Center),

		lightCellColor:  color.RGBA{R: 240, G: 217, B: 181, A: 255},
		darkCellColor:   color.RGBA{R: 181, G: 136, B: 99, A: 255},
		lightPieceColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		darkPieceColor:  color.RGBA{R: 30, G: 20, B: 10, A: 255},
	},
}
