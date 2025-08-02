package models

type JumlahKK struct {
	IDJumlahKK uint `gorm:"column:id_jumlahkk;primaryKey;autoIncrement"`
	JumlahKK   int  `gorm:"column:jumlahkk;not null"`
}

func (JumlahKK) TableName() string {
	return "jumlah_kk"
}
