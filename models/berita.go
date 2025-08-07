package models

type Berita struct {
	IDBerita  uint   `gorm:"column:id_berita;primaryKey;autoIncrement"`
	Foto      string `gorm:"type:varchar(255);not null"`
	Judul     string `gorm:"type:varchar(255);not null"`
	Deskripsi string `gorm:"type:text;not null"`
	Tanggal   string `gorm:"type:varchar(100);not null"` // Tambahan kolom tanggal manual
}

func (Berita) TableName() string {
	return "berita"
}
