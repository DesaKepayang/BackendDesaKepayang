package models

type StrukturDesa struct {
	IDStruktur uint   `gorm:"column:id_struktur;primaryKey;autoIncrement"`
	Foto       string `gorm:"type:varchar(255);not null"` // URL gambar dari Cloudinary
	FotoID     string `gorm:"type:varchar(255)"`          // Public ID gambar di Cloudinary
	Nama       string `gorm:"type:varchar(255);not null"`
	Jabatan    string `gorm:"type:varchar(255);not null"`
}

func (StrukturDesa) TableName() string {
	return "struktur_desa"
}
