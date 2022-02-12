package migrations

type TransactionLog struct {
	ID          uint
	BlockNumber uint64 `gorm:"index;type:int;not null;comment:區塊號碼"`
	TxHash      string `gorm:"index;type:varchar(255);not null;comment:交易哈希值"`
	Index       uint   `gorm:"type:int;not null;comment:索引"`
	Data        string `gorm:"type:varchar(255);not null;comment:數據"`
}

func (TransactionLog) TableName() string {
	return "transaction_log"
}
