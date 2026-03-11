package ui

import (
	"fmt"
	"image/color"

	"gochess/game"

	tea "charm.land/bubbletea/v2"
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

func (bs *boardSettings) SetBoardStyle(style lipgloss.Style) {
	bs.boardStyle = style
}

func (bs *boardSettings) SetCellStyle(
	style lipgloss.Style,
	lightCellColor color.Color,
	darkCellColor color.Color,
	lightPieceColor color.Color,
	darkPieceColor color.Color,
) {
	bs.cellStyle = style
	bs.lightCellColor = lightCellColor
	bs.darkCellColor = darkCellColor
	bs.lightPieceColor = lightPieceColor
	bs.darkPieceColor = darkPieceColor
}

// BoardModel will keep track of the board state
// and all the relevant info
//
// including height, width and style
type boardModel struct {
	game     *game.Game
	settings *boardSettings
	style    lipgloss.Style // style of the whole sub widget, use just height and width
}

func (bm *boardModel) Init() tea.Cmd {
	return nil
}

// Update the boardModel when a msg is passed
func (bm *boardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO if necessary
	//
	return bm, nil
}

func (bm *boardModel) Render() string {
	// calculate the heights and widths for the subcomponents
	//

	// Header will contain the White Graveyard
	// Footer will contain the Black Graveyard
	// TODO implement the graveyards
	//

	// header and footer will be 20% of the total height
	boardHeight := bm.style.GetHeight() * 80 / 100
	boardWidth := bm.style.GetWidth()

	minHeaderHeight, minFooterHeight := 5, 5

	footerHeight := max(
		(bm.style.GetHeight()-boardHeight)/2,
		minFooterHeight)
	footerWidth := bm.style.GetWidth()

	headerHeight := max(
		(bm.style.GetHeight()-boardHeight)/2,
		minHeaderHeight)
	headerWidth := bm.style.GetWidth()

	// make adjustments to board to respect min values
	boardHeight = bm.style.GetHeight() - footerHeight - headerHeight

	headerStyle := lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Bottom).
		Height(headerHeight).
		Width(headerWidth)

	footerStyle := lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Top).
		Height(footerHeight).
		Width(footerWidth)

	boardStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Height(boardHeight).
		Width(boardWidth)

	b, _ := bm.formatCurrentBoard()
	board := boardStyle.Render(b)

	WhiteGraveyard := "== Placeholder White Graveyard =="
	BlackGraveyard := "== Placeholder Black Graveyard =="

	header := headerStyle.Render(WhiteGraveyard)
	footer := footerStyle.Render(BlackGraveyard)

	return lipgloss.JoinVertical(lipgloss.Center, header, board, footer)
}

func (bm *boardModel) formatCurrentBoard() (string, error) {
	rankLabelStyle := lipgloss.NewStyle().Width(2)
	fileLabelStyle := lipgloss.NewStyle().Width(5).Align(lipgloss.Center)

	var rows []string
	for rank := 7; rank >= 0; rank-- {
		cells := []string{rankLabelStyle.Render(fmt.Sprintf("%d", rank+1))}
		for file := 0; file < 8; file++ {
			piece := bm.game.Board.State[rank*8+file]
			style, err := bm.getCellStyle((rank+file)%2 != 0, piece)
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
	return lipgloss.JoinVertical(lipgloss.Left, rows...), nil
}

func (bm *boardModel) getCellStyle(isLightCell bool, piece game.Piece) (lipgloss.Style, error) {
	bg := bm.settings.darkCellColor
	if isLightCell {
		bg = bm.settings.lightCellColor
	}

	switch piece.Color() {
	case game.WHITE:
		return lipgloss.NewStyle().Inherit(bm.settings.cellStyle).Background(bg).Foreground(bm.settings.lightPieceColor), nil

	case game.BLACK:
		return lipgloss.NewStyle().Inherit(bm.settings.cellStyle).Background(bg).Foreground(bm.settings.darkPieceColor), nil

	default:
		if piece == 0 {
			return lipgloss.NewStyle().Inherit(bm.settings.cellStyle).Background(bg), nil
		}
		return lipgloss.NewStyle().Inherit(bm.settings.cellStyle), fmt.Errorf("Cannot format a piece with no color %v", piece)
	}
}

func (bm *boardModel) View() tea.View {
	return tea.NewView(bm.Render())
}
