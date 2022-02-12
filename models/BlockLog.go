package models

type BlockLog struct {
	ID          uint   `json:"id"`
	BlockNumber uint64 `gorm:"not null;comment:區塊號碼" json:"block_num"`
	AvailableAt uint64 `gorm:"not null;comment:可執行時間" json:"available_at"`
}

func (BlockLog) TableName() string {
	return "block_log"
}
