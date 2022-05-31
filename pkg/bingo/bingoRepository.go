package bingo

import "restracker/pkg/db"

type BingoRepository struct {
	db *db.Database
}

func NewBingoRepository(database *db.Database) *BingoRepository {
	return &BingoRepository{db: database}
}

func (r *BingoRepository) GetBingo(channelId string) *Bingo {
	row := r.db.Query("SELECT id, channelId FROM Bingo WHERE channelId = ? LIMIT 1", channelId)
	defer row.Close()
	if row == nil {
		return nil
	}

	var id int
	var chId string
	row.Scan(&id, &chId)

	return InitBingo(id, chId, r.GetBingoCells(id), r)
}

func (r *BingoRepository) GetBingoCells(bingoId int) map[int]*BingoCell {
	row := r.db.Query("SELECT id, `index`, text, isMarked FROM BingoItems WHERE bingoId = ?", bingoId)
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

func (r *BingoRepository) CreateBingo(channelId string) *Bingo {
	r.db.Exec("INSERT INTO Bingo (channelId) VALUES (?)", channelId)

	return r.GetBingo(channelId)
}

func (r *BingoRepository) AddCell(bingoId int, text string, index int) int {
	res := r.db.Exec("INSERT INTO BingoItems (bingoId, text, `index`) VALUES (?, ?, ?)", bingoId, text, index)
	id, _ := res.LastInsertId()
	return int(id)
}

func (r *BingoRepository) RemoveCell(bingoId, index, cellId int) {
	r.db.Exec("DELETE FROM BingoItems WHERE id = ?", cellId)
	r.db.Exec("UPDATE BingoItems SET `index` = `index` - 1 WHERE  `index` > ? AND bingoId = ?", index, bingoId)
}

func (r *BingoRepository) UpdateCell(id int, index int, isMarked bool) {
	r.db.Exec("UPDATE BingoItems SET `index` = ?, isMarked = ? WHERE id = ?", index, isMarked, id)
}

func (r *BingoRepository) ResetBingo(id int) {
	r.db.Exec("UPDATE BingoItems SET isMarked = 0 WHERE bingoId = ?", id)
}
