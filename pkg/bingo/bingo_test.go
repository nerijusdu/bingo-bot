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
			1: {IsMarked: true},
			2: {IsMarked: true},
			3: {IsMarked: true},
			4: {IsMarked: false},
			5: {IsMarked: false},
			6: {IsMarked: false},
			7: {IsMarked: false},
			8: {IsMarked: false},
			9: {IsMarked: false},
		},
	},
	{
		name:     "complete column",
		expected: true,
		cells: map[int]*BingoCell{
			1: {IsMarked: true},
			2: {IsMarked: false},
			3: {IsMarked: false},
			4: {IsMarked: true},
			5: {IsMarked: false},
			6: {IsMarked: false},
			7: {IsMarked: true},
			8: {IsMarked: false},
			9: {IsMarked: false},
		},
	},
	{
		name:     "complete diagonal decline",
		expected: true,
		cells: map[int]*BingoCell{
			1: {IsMarked: true},
			2: {IsMarked: false},
			3: {IsMarked: false},
			4: {IsMarked: false},
			5: {IsMarked: true},
			6: {IsMarked: false},
			7: {IsMarked: false},
			8: {IsMarked: false},
			9: {IsMarked: true},
		},
	},
	{
		name:     "complete diagonal incline",
		expected: true,
		cells: map[int]*BingoCell{
			1: {IsMarked: false},
			2: {IsMarked: false},
			3: {IsMarked: true},
			4: {IsMarked: false},
			5: {IsMarked: true},
			6: {IsMarked: false},
			7: {IsMarked: true},
			8: {IsMarked: false},
			9: {IsMarked: false},
		},
	},
	{
		name:     "another column",
		expected: true,
		cells: map[int]*BingoCell{
			1: {IsMarked: false},
			2: {IsMarked: false},
			3: {IsMarked: true},
			4: {IsMarked: false},
			5: {IsMarked: false},
			6: {IsMarked: true},
			7: {IsMarked: false},
			8: {IsMarked: false},
			9: {IsMarked: true},
		},
	},
	{
		name:     "more items than grid size",
		expected: false,
		cells: map[int]*BingoCell{
			1:  {IsMarked: true},
			2:  {IsMarked: true},
			3:  {IsMarked: true},
			4:  {IsMarked: false},
			5:  {IsMarked: false},
			6:  {IsMarked: false},
			7:  {IsMarked: false},
			8:  {IsMarked: false},
			9:  {IsMarked: false},
			10: {IsMarked: false},
		},
	},
}

func TestIsCompleted(t *testing.T) {
	for _, test := range isCompleteTests {
		bingo := InitBingo(1, "channelId", test.cells, nil)
		assert.Equal(t, test.expected, bingo.IsCompleted(), test.name)
	}
}
