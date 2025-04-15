package orm

import (
	"github.com/buglot/postAPI/orm"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
	Avatar   string
	RoleID   orm.Role `gorm:"references:RoleID"`
}
