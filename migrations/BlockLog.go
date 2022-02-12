package migrations

type BlockLog struct {
	ID          uint
	BlockNumber uint64 `gorm:"not null;comment:區塊號碼"`
	AvailableAt uint64 `gorm:"not null;comment:可執行時間"`
}

func (BlockLog) TableName() string {
	return "block_log"
}
