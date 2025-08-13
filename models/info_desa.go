package models

type InfoDesa struct {
	IDInfo    uint   `gorm:"column:id_info;primaryKey;autoIncrement" json:"id_info"`
	Indikator string `gorm:"type:varchar(255);not null" json:"indikator"`
	Jumlah    int    `gorm:"not null" json:"jumlah"`
}

func (InfoDesa) TableName() string {
	return "info_desa"
}
