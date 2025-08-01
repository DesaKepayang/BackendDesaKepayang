package models

type RTRW struct {
	IDRTRW uint   `gorm:"primaryKey;autoIncrement" json:"id_rtrw"`
	RT     string `gorm:"type:varchar(10);not null" json:"rt"`
	RW     string `gorm:"type:varchar(10);not null" json:"rw"`
}

func (RTRW) TableName() string {
	return "rt_rw"
}
