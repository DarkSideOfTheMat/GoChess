package ui

import (
	game "gochess/game"
)

var pieceSymbols = map[game.Piece]string{
	game.KING: " ♚ ", game.QUEEN: " ♛ ", game.ROOK: " ♜ ",
	game.BISHOP: " ♝ ", game.KNIGHT: " ♞ ", game.PAWN: " ♟ ",
	0: "   ",
}

var fileLabels = [8]string{"a", "b", "c", "d", "e", "f", "g", "h"}

func getPieceSymbol(piece game.Piece) string {
	return pieceSymbols[piece&^game.COLOR_MASK]
}

func (m model) formatCurrentBoard() (string, error) {
	// Header
	return " ", nil
	// compositor := lipgloss.NewCompositor()

	// // pieceFieldBuilder := strings.Builder{}
	// // backgroundFieldBuilder := strings.Builder{}
	// // borderFieldBuilder := strings.Builder{}

	// fieldBuilder := strings.Builder{}
	// plain_style := lipgloss.NewStyle()

	// for rank := 7; rank >= 0; rank-- {
	// 	for file := range 8 {
	// 		piece := m.game.Board.State[rank*8+file]
	// 		cell, err := m.settings.getCellStyle((rank+file)%2 != 0, piece)

	// 		if err != nil {

	// 			return "", fmt.Errorf("Cannot render board! There is an error on rank %v file %s %s", rank+1, fileLabels[file], err)
	// 		}

	// 		fieldBuilder.WriteString(cell.Render(getPieceSymbol(piece)))
	// 	}

	// 	fieldBuilder.WriteString(plain_style.Render("\n"))
	// }

	// board_layer := lipgloss.NewLayer(fieldBuilder.String())

	// canvas := lipgloss.NewCanvas(800, 800)

	// canvas.Compose(board_layer)
	// return canvas.Render(), nil
	// // return m.settings.renderBoard(fieldBuilder.String()), nil
	// // return fieldBuilder.String(), nil
}
