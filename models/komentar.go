package models

type Komentar struct {
	IDKomentar uint   `gorm:"column:id_komentar;primaryKey;autoIncrement" json:"id_komentar"`
	Nama       string `gorm:"column:nama;type:varchar(100);not null" json:"nama"`
	Email      string `gorm:"column:email;type:varchar(100);not null" json:"email"`
	NoHP       string `gorm:"column:no_hp;type:varchar(20);not null" json:"no_hp"`
	Komentar   string `gorm:"column:komentar;type:text;not null" json:"komentar"`
}

func (Komentar) TableName() string {
	return "komentar"
}
