package models

type Admin struct {
	ID       uint   `gorm:"column:id_admin;primaryKey;autoIncrement" json:"id_admin"`
	Username string `gorm:"column:username;type:varchar(100);not null" json:"username"`
	Password string `gorm:"column:password;type:varchar(255);not null" json:"-"`
	Role     string `gorm:"column:role;type:varchar(50);not null" json:"role"`
}

func (Admin) TableName() string {
	return "admin"
}
