package models

type StrukturDesa struct {
	IDStruktur uint   `gorm:"column:id_struktur;primaryKey;autoIncrement"`
	Foto       string `gorm:"type:varchar(255);not null"`
	Nama       string `gorm:"type:varchar(255);not null"`
	Jabatan    string `gorm:"type:varchar(255);not null"`
}

func (StrukturDesa) TableName() string {
	return "struktur_desa"
}
