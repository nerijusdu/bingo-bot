package visual

import (
	"bingobot/pkg/bingo"
	"io"

	"github.com/fogleman/gg"
)

func GenerateImage(b *bingo.Bingo, w io.Writer) error {
	width := 300
	height := 300
	cellHeight := float64(height / b.LineLength)
	cellWidth := float64(width / b.LineLength)

	dc := gg.NewContext(width, height)

	for i := 1; i <= b.LineLength; i++ {
		for j := 1; j <= b.LineLength; j++ {
			cell, ok := b.Cells[(i-1)*b.LineLength+j]
			if !ok {
				continue
			}

			drawCell(dc, cellHeight, cellWidth, i, j, cell)
		}
	}

	return dc.EncodePNG(w)
}

func drawCell(dc *gg.Context, cellHeight, cellWidth float64, i, j int, cell *bingo.BingoCell) {
	x1 := float64(j-1) * cellHeight
	y1 := float64(i-1) * cellWidth
	x2 := float64(j) * cellHeight
	y2 := float64(i) * cellWidth

	dc.SetRGB255(38, 70, 83)
	dc.DrawRectangle(x1, y1, x2, y2)
	dc.Stroke()

	dc.SetRGB255(233, 196, 106)
	if cell.IsMarked {
		dc.SetRGB255(42, 157, 143)
	}
	dc.DrawRectangle(x1, y1, x2, y2)
	dc.Fill()

	dc.SetRGB255(38, 70, 83)
	dc.DrawStringWrapped(
		cell.Text,
		x1+cellWidth/2,
		y1+cellHeight/2,
		0.5, 0.5,
		cellWidth,
		1,
		gg.AlignCenter,
	)
}
