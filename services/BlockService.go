package services

import (
	"go-ethereum-demo/databases"
	"go-ethereum-demo/models"
	"go-ethereum-demo/requests"
	"gorm.io/gorm"
)

const defaultLimit = 10

func GetBlocksList(params requests.GetBlocksListParam) (blocksList []*models.Block, err error) {
	tx := databases.MySqlClient.Select("id", "number", "hash", "time", "parent_hash", "is_pending")

	limit := defaultLimit
	if params.Limit != 0 {
		limit = params.Limit
	}

	err = tx.Order("id desc").Limit(limit).Find(&blocksList).Error

	return
}

type BlockDetail struct {
	Hash         string   `json:"block_hash"`
	Number       uint64   `json:"block_num"`
	Time         uint64   `json:"block_time"`
	ParentHash   string   `json:"parent_hash"`
	IsPending    int8     `json:"is_pending"`
	Transactions []string `json:"transactions"`
}

func GetBlockById(id string) (blockDetail *BlockDetail, err error)  {
	var block models.Block
	var txHashes []string
	err = databases.MySqlClient.Preload("Transactions", func(db *gorm.DB) *gorm.DB {
		return db.Select("block_number", "hash")
	}).Select("number", "hash", "time", "parent_hash", "is_pending").First(&block, id).Error

	for _, transaction := range block.Transactions {
		txHashes = append(txHashes, transaction.Hash)
	}

	blockDetail = &BlockDetail{
		Hash: block.Hash,
		Number: block.Number,
		Time: block.Time,
		ParentHash: block.ParentHash,
		IsPending: block.IsPending,
		Transactions: txHashes,
	}

	return
}
