package bingo

import (
	"bingobot/pkg/db"
	"fmt"
)

type BingoRepository struct {
	db *db.Database
}

func NewBingoRepository(database *db.Database) *BingoRepository {
	return &BingoRepository{db: database}
}

func (r *BingoRepository) GetBingo(channelId string) *Bingo {
	row, err := r.db.Query("SELECT id, channelId FROM Bingo WHERE channelId = ? LIMIT 1", channelId)
	if err != nil || row == nil {
		return nil
	}

	defer row.Close()

	var id int
	var chId string
	hasResult := row.Next()
	row.Scan(&id, &chId)
	if !hasResult || id == 0 {
		fmt.Println("No bingo found for channel: " + channelId)
		return nil
	}

	return InitBingo(id, chId, r.GetBingoCells(id), r)
}

func (r *BingoRepository) GetBingoCells(bingoId int) map[int]*BingoCell {
	row, err := r.db.Query("SELECT id, `index`, text, isMarked FROM BingoItems WHERE bingoId = ?", bingoId)
	if err != nil || row == nil {
		return nil
	}

	defer row.Close()

	cells := map[int]*BingoCell{}
	for row.Next() {
		var id int
		var index int
		var text string
		var isMarked bool
		row.Scan(&id, &index, &text, &isMarked)
		cells[index] = &BingoCell{
			id:       id,
			Text:     text,
			IsMarked: isMarked,
		}
	}

	return cells
}

func (r *BingoRepository) CreateBingo(channelId string) (*Bingo, error) {
	_, err := r.db.Insert("INSERT INTO Bingo (channelId) VALUES (?)", channelId)
	if err != nil {
		return nil, err
	}

	return r.GetBingo(channelId), nil
}

func (r *BingoRepository) AddCell(bingoId int, text string, index int) (int, error) {
	return r.db.Insert("INSERT INTO BingoItems (bingoId, text, `index`) VALUES (?, ?, ?)", bingoId, text, index)
}

func (r *BingoRepository) RemoveCell(bingoId, index, cellId int) error {
	err := r.db.Exec("DELETE FROM BingoItems WHERE id = ?", cellId)
	if err != nil {
		return err
	}

	return r.db.Exec("UPDATE BingoItems SET `index` = `index` - 1 WHERE  `index` > ? AND bingoId = ?", index, bingoId)
}

func (r *BingoRepository) UpdateCell(id int, index int, isMarked bool) error {
	return r.db.Exec("UPDATE BingoItems SET `index` = ?, isMarked = ? WHERE id = ?", index, isMarked, id)
}

func (r *BingoRepository) ResetBingo(id int) error {
	return r.db.Exec("UPDATE BingoItems SET isMarked = 0 WHERE bingoId = ?", id)
}
