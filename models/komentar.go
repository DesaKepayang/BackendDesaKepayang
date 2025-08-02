package models

type Komentar struct {
	IDKomentar uint   `gorm:"column:id_komentar;primaryKey;autoIncrement"`
	Nama       string `gorm:"column:nama;type:varchar(100);not null"`
	Email      string `gorm:"column:email;type:varchar(100);not null"`
	NoHP       string `gorm:"column:no_hp;type:varchar(20);not null"`
	Komentar   string `gorm:"column:komentar;type:text;not null"`
}

func (Komentar) TableName() string {
	return "komentar"
}
