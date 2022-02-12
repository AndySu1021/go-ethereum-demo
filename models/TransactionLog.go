package models

type TransactionLog struct {
	ID          uint   `gorm:"comment:主鍵" json:"id,omitempty"`
	BlockNumber uint64 `gorm:"type:int;not null;comment:區塊號碼" json:"block_num,omitempty"`
	TxHash      string `gorm:"type:varchar(255);not null;comment:交易哈希值" json:"tx_hash,omitempty"`
	Index       uint   `gorm:"type:int;not null;comment:索引" json:"index"`
	Data        string `gorm:"type:varchar(255);not null;comment:數據" json:"data"`
}

func (TransactionLog) TableName() string {
	return "transaction_log"
}
