package models

type Transaction struct {
	ID              uint             `gorm:"comment:主鍵" json:"id,omitempty"`
	BlockNumber     uint64           `gorm:"type:int;not null;comment:區塊號碼" json:"block_num,omitempty"`
	Hash            string           `gorm:"type:varchar(255);not null;comment:交易哈希值" json:"tx_hash,omitempty"`
	From            string           `gorm:"type:varchar(255);not null;comment:來源" json:"from,omitempty"`
	To              string           `gorm:"type:varchar(255);not null;comment:終點" json:"to,omitempty"`
	Nonce           uint64           `gorm:"type:int;not null;comment:區塊隨機值" json:"nonce,omitempty"`
	Data            string           `gorm:"type:varchar(255);not null;comment:數據" json:"data,omitempty"`
	Value           string           `gorm:"type:varchar(255);not null;comment:值" json:"value,omitempty"`
	TransactionLogs []TransactionLog `gorm:"foreignKey:TxHash;references:Hash" json:"logs"`
}

func (Transaction) TableName() string {
	return "transaction"
}
