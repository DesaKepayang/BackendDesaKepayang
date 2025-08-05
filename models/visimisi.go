package models

import "gorm.io/gorm"

type VisiMisi struct {
	IDVisiMisi uint   `gorm:"column:id_visimisi;primaryKey;autoIncrement"`
	Visi       string `gorm:"type:varchar(855);not null"`
	Misi       string `gorm:"type:varchar(855);not null"`
}

func (VisiMisi) TableName() string {
	return "visi_misi"
}

// Fungsi untuk drop kolom foto jika masih ada di DB
func DropFotoColumn(db *gorm.DB) {
	db.Migrator().DropColumn(&VisiMisi{}, "foto")
}
