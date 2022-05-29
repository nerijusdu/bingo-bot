package bingo

import (
	"fmt"
	"strings"
)

type Bingo struct {
	id         int
	cells      map[int]*BingoCell
	channelId  string
	gridSize   int
	lineLength int

	repository *BingoRepository
}

type BingoCell struct {
	id       int
	Text     string
	isMarked bool
}

func NewBingo(channelId string) *Bingo {
	return &Bingo{
		channelId: channelId,
		cells:     map[int](*BingoCell){},
	}
}

func InitBingo(id int, channelId string, cells map[int]*BingoCell, r *BingoRepository) *Bingo {
	bingo := &Bingo{
		id:         id,
		channelId:  channelId,
		cells:      cells,
		repository: r,
	}

	bingo.updateGridSize()

	return bingo
}

func (b *Bingo) AddCell(text string) int {
	i := len(b.cells) + 1
	b.cells[i] = &BingoCell{
		Text:     text,
		isMarked: false,
	}

	b.repository.AddCell(b.id, text, i)

	b.updateGridSize()

	return i
}

func (b *Bingo) RemoveCell(i int) bool {
	cell, ok := b.cells[i]
	if !ok {
		return false
	}

	for i := i; i < len(b.cells); i++ {
		b.cells[i] = b.cells[i+1]
	}

	b.repository.RemoveCell(cell.id)

	b.updateGridSize()
	return true
}

func (b *Bingo) SwitchCells(i1 int, i2 int) bool {
	cell1, ok1 := b.cells[i1]
	cell2, ok2 := b.cells[i2]
	if !ok1 || !ok2 {
		return false
	}

	b.cells[i1] = cell2
	b.cells[i2] = cell1

	b.repository.UpdateCell(cell1.id, i2, cell2.isMarked)
	b.repository.UpdateCell(cell2.id, i1, cell1.isMarked)

	return true
}

func (b *Bingo) MarkCell(i int) bool {
	cell, ok := b.cells[i]
	if !ok {
		return false
	}

	cell.isMarked = true

	b.repository.UpdateCell(cell.id, i, true)

	return true
}

func (b *Bingo) IsCompleted() bool {
	diagonalIncline := 0
	diagonalDecline := 0
	columns := make([]int, b.lineLength)

	for i := 1; i <= b.lineLength; i++ {
		numbersInARow := 0

		for j := 1; j <= b.lineLength; j++ {
			cell, ok := b.cells[(i-1)*b.lineLength+j]
			if !ok || !cell.isMarked {
				continue
			}

			numbersInARow++
			columns[j-1] = columns[j-1] + 1

			if i == j {
				diagonalIncline++
			}

			if i+j == b.lineLength+1 {
				diagonalDecline++
			}

			if numbersInARow == b.lineLength ||
				diagonalIncline == b.lineLength ||
				diagonalDecline == b.lineLength ||
				columns[j-1] == b.lineLength {
				return true
			}
		}
	}

	return false
}

func (b *Bingo) Reset() {
	for _, v := range b.cells {
		v.isMarked = false
	}

	b.repository.ResetBingo(b.id)
}

func (b *Bingo) ToString() string {
	var items []string
	for i := 1; i <= len(b.cells); i++ {
		markedText := ""
		if b.cells[i].isMarked {
			markedText = " :white_check_mark:"
		}
		items = append(items, fmt.Sprintf("%d. %s%s", i, b.cells[i].Text, markedText))
	}

	return strings.Join(items, "\n")
}

func (b *Bingo) updateGridSize() {
	count := len(b.cells)
	switch {
	case count == 1:
		b.gridSize = 1
		b.lineLength = 1
	case count <= 4:
		b.gridSize = 4
		b.lineLength = 2
	case count <= 9:
		b.gridSize = 9
		b.lineLength = 3
	case count <= 16:
		b.gridSize = 16
		b.lineLength = 4
	default:
		b.gridSize = 25
		b.lineLength = 5
	}
}
