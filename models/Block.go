package models

type Block struct {
	ID           uint          `gorm:"comment:主鍵" json:"id,omitempty"`
	Number       uint64        `gorm:"not null;comment:區塊號碼" json:"block_num,omitempty"`
	Hash         string        `gorm:"type:varchar(255);not null;comment:區塊哈希值" json:"block_hash,omitempty"`
	ParentHash   string        `gorm:"type:varchar(255);comment:父區塊哈希值" json:"parent_hash"`
	IsPending    int8          `gorm:"type:tinyint;not null;comment:是否屬於不穩定區塊 1是 2否" json:"is_pending"`
	Time         uint64        `gorm:"not null;comment:區塊創建時間" json:"block_time,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:BlockNumber;references:Number" json:"transactions,omitempty"`
}

func (Block) TableName() string {
	return "block"
}
