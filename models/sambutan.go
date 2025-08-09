package models

type SambutanKepalaDesa struct {
	ID             uint   `gorm:"column:id_sambutan;primaryKey;autoIncrement"`
	Foto           string `gorm:"column:foto_kepaladesa;type:varchar(255)"`
	FotoID         string `gorm:"column:foto_id;type:varchar(255)"` // Tambahan untuk Cloudinary
	KataSambutan   string `gorm:"column:kata_sambutan;type:text"`
	NamaKepalaDesa string `gorm:"column:nama_kepaladesa;type:varchar(255)"`
}

func (SambutanKepalaDesa) TableName() string {
	return "sambutan_kepaladesa"
}
