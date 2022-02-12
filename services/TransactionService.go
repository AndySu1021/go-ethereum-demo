package services

import (
	"go-ethereum-demo/databases"
	"go-ethereum-demo/models"
	"gorm.io/gorm"
)

func GetTransactionByHash(txHash string) (tx *models.Transaction, err error)  {
	err = databases.MySqlClient.Preload("TransactionLogs", func(db *gorm.DB) *gorm.DB {
		return db.Select("tx_hash", "index", "data")
	}).Select("hash", "from", "to", "nonce", "data", "value").Where("hash = ?", txHash).First(&tx).Error

	for i := 0; i < len(tx.TransactionLogs); i++ {
		tx.TransactionLogs[i].TxHash = ""
	}

	return
}
