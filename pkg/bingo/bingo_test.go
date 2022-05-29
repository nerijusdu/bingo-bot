package bingo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type cellCompleteTest struct {
	cells    map[int]*BingoCell
	expected bool
	name     string
}

var isCompleteTests = []cellCompleteTest{
	{
		name:     "complete row",
		expected: true,
		cells: map[int]*BingoCell{
			1: {isMarked: true},
			2: {isMarked: true},
			3: {isMarked: true},
			4: {isMarked: false},
			5: {isMarked: false},
			6: {isMarked: false},
			7: {isMarked: false},
			8: {isMarked: false},
			9: {isMarked: false},
		},
	},
	{
		name:     "complete column",
		expected: true,
		cells: map[int]*BingoCell{
			1: {isMarked: true},
			2: {isMarked: false},
			3: {isMarked: false},
			4: {isMarked: true},
			5: {isMarked: false},
			6: {isMarked: false},
			7: {isMarked: true},
			8: {isMarked: false},
			9: {isMarked: false},
		},
	},
	{
		name:     "complete diagonal decline",
		expected: true,
		cells: map[int]*BingoCell{
			1: {isMarked: true},
			2: {isMarked: false},
			3: {isMarked: false},
			4: {isMarked: false},
			5: {isMarked: true},
			6: {isMarked: false},
			7: {isMarked: false},
			8: {isMarked: false},
			9: {isMarked: true},
		},
	},
	{
		name:     "complete diagonal incline",
		expected: true,
		cells: map[int]*BingoCell{
			1: {isMarked: false},
			2: {isMarked: false},
			3: {isMarked: true},
			4: {isMarked: false},
			5: {isMarked: true},
			6: {isMarked: false},
			7: {isMarked: true},
			8: {isMarked: false},
			9: {isMarked: false},
		},
	},
	{
		name:     "another column",
		expected: true,
		cells: map[int]*BingoCell{
			1: {isMarked: false},
			2: {isMarked: false},
			3: {isMarked: true},
			4: {isMarked: false},
			5: {isMarked: false},
			6: {isMarked: true},
			7: {isMarked: false},
			8: {isMarked: false},
			9: {isMarked: true},
		},
	},
	{
		name:     "more items than grid size",
		expected: false,
		cells: map[int]*BingoCell{
			1:  {isMarked: true},
			2:  {isMarked: true},
			3:  {isMarked: true},
			4:  {isMarked: false},
			5:  {isMarked: false},
			6:  {isMarked: false},
			7:  {isMarked: false},
			8:  {isMarked: false},
			9:  {isMarked: false},
			10: {isMarked: false},
		},
	},
}

func TestIsCompletedCompletes(t *testing.T) {
	for _, test := range isCompleteTests {
		bingo := InitBingo(1, "channelId", test.cells, nil)
		assert.Equal(t, test.expected, bingo.IsCompleted(), test.name)
	}
}
