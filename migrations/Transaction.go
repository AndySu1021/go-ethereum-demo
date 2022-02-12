package migrations

type Transaction struct {
	ID          uint
	BlockNumber uint64 `gorm:"index;type:int;not null;comment:區塊號碼"`
	Hash        string `gorm:"index;type:varchar(255);not null;comment:交易哈希值"`
	From        string `gorm:"type:varchar(255);not null;comment:來源"`
	To          string `gorm:"type:varchar(255);not null;comment:終點"`
	Nonce       uint64 `gorm:"type:int;not null;comment:區塊隨機值"`
	Data        string `gorm:"type:varchar(255);not null;comment:數據"`
	Value       string `gorm:"type:varchar(255);not null;comment:值"`
}

func (Transaction) TableName() string {
	return "transaction"
}
