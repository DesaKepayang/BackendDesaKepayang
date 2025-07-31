package models

type DataPenduduk struct {
	IDPenduduk uint `gorm:"column:id_penduduk;primaryKey;autoIncrement"`
	IDRTRW     uint `gorm:"column:id_rtrw;not null"`

	Nama   string `gorm:"type:varchar(100);not null"`
	Agama  string `gorm:"type:varchar(50);not null"`
	Gender string `gorm:"type:varchar(10);not null"`

	// Relasi dengan tabel rt_rw
	RTRW RTRW `gorm:"foreignKey:IDRTRW;references:IDRTRW;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (DataPenduduk) TableName() string {
	return "data_penduduk"
}
