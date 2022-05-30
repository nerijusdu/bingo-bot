package bingo

import (
	"fmt"
	"strings"
)

type Bingo struct {
	id         int
	Cells      map[int]*BingoCell
	channelId  string
	gridSize   int
	LineLength int

	repository *BingoRepository
}

type BingoCell struct {
	id       int
	Text     string
	IsMarked bool
}

func NewBingo(channelId string) *Bingo {
	return &Bingo{
		channelId: channelId,
		Cells:     map[int](*BingoCell){},
	}
}

func InitBingo(id int, channelId string, cells map[int]*BingoCell, r *BingoRepository) *Bingo {
	bingo := &Bingo{
		id:         id,
		channelId:  channelId,
		Cells:      cells,
		repository: r,
	}

	bingo.updateGridSize()

	return bingo
}

func (b *Bingo) AddCell(text string) int {
	i := len(b.Cells) + 1
	b.Cells[i] = &BingoCell{
		Text:     text,
		IsMarked: false,
	}

	b.repository.AddCell(b.id, text, i)

	b.updateGridSize()

	return i
}

func (b *Bingo) RemoveCell(i int) bool {
	cell, ok := b.Cells[i]
	if !ok {
		return false
	}

	for i := i; i < len(b.Cells); i++ {
		b.Cells[i] = b.Cells[i+1]
	}

	b.repository.RemoveCell(cell.id)

	b.updateGridSize()
	return true
}

func (b *Bingo) SwitchCells(i1 int, i2 int) bool {
	cell1, ok1 := b.Cells[i1]
	cell2, ok2 := b.Cells[i2]
	if !ok1 || !ok2 {
		return false
	}

	b.Cells[i1] = cell2
	b.Cells[i2] = cell1

	b.repository.UpdateCell(cell1.id, i2, cell2.IsMarked)
	b.repository.UpdateCell(cell2.id, i1, cell1.IsMarked)

	return true
}

func (b *Bingo) MarkCell(i int) bool {
	cell, ok := b.Cells[i]
	if !ok {
		return false
	}

	cell.IsMarked = true

	b.repository.UpdateCell(cell.id, i, true)

	return true
}

func (b *Bingo) IsCompleted() bool {
	diagonalIncline := 0
	diagonalDecline := 0
	columns := make([]int, b.LineLength)

	for i := 1; i <= b.LineLength; i++ {
		numbersInARow := 0

		for j := 1; j <= b.LineLength; j++ {
			cell, ok := b.Cells[(i-1)*b.LineLength+j]
			if !ok || !cell.IsMarked {
				continue
			}

			numbersInARow++
			columns[j-1] = columns[j-1] + 1

			if i == j {
				diagonalIncline++
			}

			if i+j == b.LineLength+1 {
				diagonalDecline++
			}

			if numbersInARow == b.LineLength ||
				diagonalIncline == b.LineLength ||
				diagonalDecline == b.LineLength ||
				columns[j-1] == b.LineLength {
				return true
			}
		}
	}

	return false
}

func (b *Bingo) Reset() {
	for _, v := range b.Cells {
		v.IsMarked = false
	}

	b.repository.ResetBingo(b.id)
}

func (b *Bingo) ToString() string {
	var items []string
	for i := 1; i <= len(b.Cells); i++ {
		markedText := ""
		if b.Cells[i].IsMarked {
			markedText = " :white_check_mark:"
		}
		items = append(items, fmt.Sprintf("%d. %s%s", i, b.Cells[i].Text, markedText))
	}

	return strings.Join(items, "\n")
}

func (b *Bingo) updateGridSize() {
	count := len(b.Cells)
	switch {
	case count == 1:
		b.gridSize = 1
		b.LineLength = 1
	case count <= 4:
		b.gridSize = 4
		b.LineLength = 2
	case count <= 9:
		b.gridSize = 9
		b.LineLength = 3
	case count <= 16:
		b.gridSize = 16
		b.LineLength = 4
	default:
		b.gridSize = 25
		b.LineLength = 5
	}
}
