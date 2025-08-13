package models

type InfoDesa struct {
	IDInfo        uint `gorm:"column:id_info;primaryKey;autoIncrement"`
	JumlahKK      int  `gorm:"column:jumlahkk;not null"`
	TotalPenduduk int  `gorm:"column:total_penduduk;not null"`
	LakiLaki      int  `gorm:"column:laki_laki;not null"`
	Perempuan     int  `gorm:"column:perempuan;not null"`
	Islam         int  `gorm:"column:islam;not null"`
	Buddha        int  `gorm:"column:buddha;not null"`
	Kristen       int  `gorm:"column:kristen;not null"`
	Katolik       int  `gorm:"column:katolik;not null"`
	Konghucu      int  `gorm:"column:konghucu;not null"`
	LainLain      int  `gorm:"column:lain_lain;not null"`
}

func (InfoDesa) TableName() string {
	return "info_desa"
}
