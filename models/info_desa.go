package models

type InfoDesa struct {
	IDInfo        uint `gorm:"column:id_info;primaryKey;autoIncrement" json:"id_info"`
	JumlahKK      int  `gorm:"column:jumlahkk;not null" json:"jumlahkk"`
	TotalPenduduk int  `gorm:"column:total_penduduk;not null" json:"total_penduduk"`
	LakiLaki      int  `gorm:"column:laki_laki;not null" json:"laki_laki"`
	Perempuan     int  `gorm:"column:perempuan;not null" json:"perempuan"`
	Islam         int  `gorm:"column:islam;not null" json:"islam"`
	Buddha        int  `gorm:"column:buddha;not null" json:"buddha"`
	Kristen       int  `gorm:"column:kristen;not null" json:"kristen"`
	Katolik       int  `gorm:"column:katolik;not null" json:"katolik"`
	Konghucu      int  `gorm:"column:konghucu;not null" json:"konghucu"`
	LainLain      int  `gorm:"column:lain_lain;not null" json:"lain_lain"`
}

func (InfoDesa) TableName() string {
	return "info_desa"
}
