package models

type DataPenduduk struct {
	IDPenduduk uint `gorm:"primaryKey;autoIncrement" json:"id_penduduk"`
	IDRTRW     uint `gorm:"not null" json:"id_rtrw"`

	Nama   string `gorm:"type:varchar(100);not null" json:"nama"`
	Agama  string `gorm:"type:varchar(50);not null" json:"agama"`
	Gender string `gorm:"type:varchar(10);not null" json:"gender"`

	RTRW *RTRW `gorm:"foreignKey:IDRTRW;references:IDRTRW" json:"-"`
}

func (DataPenduduk) TableName() string {
	return "data_penduduk"
}
