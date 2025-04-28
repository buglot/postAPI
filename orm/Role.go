package orm

import (
	"fmt"

	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name string `gorm:"not null"`
}

func RoleDefault() {
	var role Role
	err := Db.Where("name = ?", "user").First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			role = Role{Name: "user"}
			if err := Db.Create(&role).Error; err != nil {
				fmt.Println("Failed to create role:", err)
				return
			}
		} else {
			fmt.Println("Failed to find role:", err)
			return
		}
	}
	err = Db.Where("name = ?", "admin").First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			role = Role{Name: "admin"}
			if err := Db.Create(&role).Error; err != nil {
				fmt.Println("Failed to create role:", err)
				return
			}
		} else {
			fmt.Println("Failed to find role:", err)
			return
		}
	}
}
