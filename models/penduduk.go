package models

type DataPenduduk struct {
	IDPenduduk uint `gorm:"column:id_penduduk;primaryKey;autoIncrement"`
	IDRTRW     uint `gorm:"column:id_rtrw;not null"`

	Nama   string `gorm:"type:varchar(100);not null"`
	Agama  string `gorm:"type:varchar(50);not null"`
	Gender string `gorm:"type:varchar(10);not null"`

	// HAPUS relasi struct dan ganti dengan ini:
	RTRW RTRW `gorm:"-"`
}

func (DataPenduduk) TableName() string {
	return "data_penduduk"
}
