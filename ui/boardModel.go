package ui

import (
	"fmt"
	"image/color"

	"gochess/chess"

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
	state    [64]chess.Piece
	settings *boardSettings
	style       lipgloss.Style
	boardStyle  lipgloss.Style
	headerStyle lipgloss.Style
	footerStyle lipgloss.Style
}

func newBoardModel(state [64]chess.Piece, settings *boardSettings) boardModel {
	return boardModel{
		state:       state,
		settings:    settings,
		style:       lipgloss.NewStyle(), // overall style of the board model
		boardStyle:  lipgloss.NewStyle(), // style of the subcomponent
		headerStyle: lipgloss.NewStyle(), // style of the header
		footerStyle: lipgloss.NewStyle(), // style fo the footer
	}
}

func (bm *boardModel) Init() tea.Cmd {
	return tea.RequestWindowSize
}

// Update the boardModel when a msg is passed
func (bm *boardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		bm.style.Height(msg.Height)
		bm.style.Width(msg.Width)

		// this way height and width calculations respect other style considerations
		// like borders and padding
		height, width := bm.style.GetHeight(), bm.style.GetWidth()
		// header and footer will be 20% of the total height
		boardHeight := height * 80 / 100
		boardWidth := width

		minHeaderHeight, minFooterHeight := 5, 5

		footerHeight := max(
			(height-boardHeight)/2,
			minFooterHeight)
		footerWidth := width

		headerHeight := max(
			(height-boardHeight)/2,
			minHeaderHeight)
		headerWidth := width

		// make adjustments to board to respect min values
		boardHeight = height - footerHeight - headerHeight

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

		bm.boardStyle = boardStyle
		bm.headerStyle = headerStyle
		bm.footerStyle = footerStyle
	}
	return bm, nil
}

func (bm *boardModel) Render() string {
	// Header will contain the White Graveyard
	// Footer will contain the Black Graveyard
	// TODO implement the graveyards

	b, _ := bm.formatCurrentBoard()
	board := bm.boardStyle.Render(b)

	WhiteGraveyard := "== Placeholder White Graveyard =="
	BlackGraveyard := "== Placeholder Black Graveyard =="

	header := bm.headerStyle.Render(WhiteGraveyard)
	footer := bm.footerStyle.Render(BlackGraveyard)

	return lipgloss.JoinVertical(lipgloss.Center, header, board, footer)
}

func (bm *boardModel) formatCurrentBoard() (string, error) {
	rankLabelStyle := lipgloss.NewStyle().Width(2)
	fileLabelStyle := lipgloss.NewStyle().Width(5).Align(lipgloss.Center)

	var rows []string
	for rank := 7; rank >= 0; rank-- {
		cells := []string{rankLabelStyle.Render(fmt.Sprintf("%d", rank+1))}
		for file := 0; file < 8; file++ {
			piece := bm.state[rank*8+file]
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

func (bm *boardModel) getCellStyle(isLightCell bool, piece chess.Piece) (lipgloss.Style, error) {
	bg := bm.settings.darkCellColor
	if isLightCell {
		bg = bm.settings.lightCellColor
	}

	switch piece.ToColor() {
	case chess.WHITE:
		return lipgloss.NewStyle().Inherit(bm.settings.cellStyle).Background(bg).Foreground(bm.settings.lightPieceColor), nil

	case chess.BLACK:
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
