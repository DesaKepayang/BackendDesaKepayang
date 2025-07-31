package models

type RTRW struct {
	IDRTRW uint   `gorm:"column:id_rtrw;primaryKey;autoIncrement"`
	RT     string `gorm:"type:varchar(10);not null"`
	RW     string `gorm:"type:varchar(10);not null"`
}

func (RTRW) TableName() string {
	return "rt_rw"
}
