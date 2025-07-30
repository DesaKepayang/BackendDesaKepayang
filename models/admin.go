package models

type Admin struct {
	ID       uint   `gorm:"column:id_admin;primaryKey;autoIncrement"`
	Username string `gorm:"type:varchar(100);not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Role     string `gorm:"type:varchar(50);not null"`
}

func (Admin) TableName() string {
	return "admin"
}
