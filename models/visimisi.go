package models

type VisiMisi struct {
	IDVisiMisi uint   `gorm:"column:id_visimisi;primaryKey;autoIncrement"`
	Visi       string `gorm:"type:varchar(855);not null"`
	Misi       string `gorm:"type:varchar(855);not null"`
}

func (VisiMisi) TableName() string {
	return "visi_misi"
}
