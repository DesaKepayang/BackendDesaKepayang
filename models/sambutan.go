package models

type SambutanKepalaDesa struct {
	ID           uint   `gorm:"column:id_sambutan;primaryKey;autoIncrement"`
	Foto         string `gorm:"column:foto_kepaladesa;type:varchar(255)"`
	KataSambutan string `gorm:"column:kata_sambutan;type:text"`
}

// Optional: Jika ingin nama tabel tetap "sambutan_kepaladesa"
func (SambutanKepalaDesa) TableName() string {
	return "sambutan_kepaladesa"
}
