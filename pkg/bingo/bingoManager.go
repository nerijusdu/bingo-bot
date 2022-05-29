package bingo

import "restracker/pkg/db"

type BingoManager struct {
	repository *BingoRepository
	cache      map[string]*Bingo
}

func NewBingoManager(database *db.Database) *BingoManager {
	return &BingoManager{
		repository: NewBingoRepository(database),
		cache:      map[string]*Bingo{},
	}
}

func (mgr *BingoManager) Get(channelId string) *Bingo {
	bingo, ok := mgr.cache[channelId]
	if ok {
		return bingo
	}

	bingo = mgr.repository.GetBingo(channelId)
	if bingo != nil {
		mgr.cache[channelId] = bingo
	}

	return bingo
}

func (mgr *BingoManager) Create(channelId string) *Bingo {
	bingo := mgr.repository.CreateBingo(channelId)
	mgr.cache[channelId] = bingo

	return bingo
}

func (mgr *BingoManager) GetOrCreate(channelId string) *Bingo {
	bingo := mgr.Get(channelId)
	if bingo == nil {
		bingo = mgr.Create(channelId)
	}

	return bingo
}
