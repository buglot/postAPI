package orm

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
	Url      string
	Avatar   string
	RoleID   uint
	Role     Role `gorm:"foreignKey:RoleID"`
}
