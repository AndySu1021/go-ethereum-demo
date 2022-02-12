package migrations

type Block struct {
	ID         uint
	Number     uint64 `gorm:"index;not null;comment:區塊號碼"`
	Hash       string `gorm:"type:varchar(255);not null;comment:區塊哈希值"`
	ParentHash string `gorm:"type:varchar(255);comment:父區塊哈希值"`
	IsPending  int8   `gorm:"type:tinyint;not null;comment:是否屬於不穩定區塊 1是 2否"`
	Time       uint64 `gorm:"not null;comment:區塊創建時間"`
}

func (Block) TableName() string {
	return "block"
}
