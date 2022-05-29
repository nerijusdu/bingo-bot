package bingo

import "restracker/pkg/db"

type BingoManager struct {
	repository *BingoRepository
}

func NewBingoManager(database *db.Database) *BingoManager {
	return &BingoManager{
		repository: NewBingoRepository(database),
	}
}

func (mgr *BingoManager) Get(channelId string) *Bingo { // TODO: cache?
	return mgr.repository.GetBingo(channelId)
}

func (mgr *BingoManager) Create(channelId string) *Bingo {
	bingo := mgr.repository.CreateBingo(channelId)

	return bingo
}

func (mgr *BingoManager) GetOrCreate(channelId string) *Bingo {
	bingo := mgr.Get(channelId)
	if bingo == nil {
		bingo = mgr.Create(channelId)
	}

	return bingo
}
