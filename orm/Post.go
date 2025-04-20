package orm

import (
	"fmt"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Url          string `gorm:"type:varchar(255);not null"`
	UserID       uint
	User         User    `gorm:"foreignKey:UserID"`
	Image        []Image `gorm:"foreignKey:PostID"`
	Message      string
	AccessID     uint
	Access       Access `gorm:"foreignKey:AccessID"`
	TypeofPostID uint
	TypeofPost   TypeofPost `gorm:"foreignKey:TypeofPostID"`
}
type Access struct {
	gorm.Model
	Name string `gorm:"type:enum('public','follow','private');default:'public'"`
}
type TypeofPost struct {
	gorm.Model
	Name string `grom:"type:varchar(30);not null"`
}
type Image struct {
	gorm.Model
	Url    string `gorm:"type:varchar(255);not null"`
	PostID uint   // foreign key to Post
}

func TypeofPostDefault() {
	var typepost TypeofPost
	err := Db.Where("name = ?", "public").First(&typepost).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			typepost = TypeofPost{Name: "public"}
			if err := Db.Create(&typepost).Error; err != nil {
				fmt.Println("Failed to create role:", err)
				return
			}
		} else {
			fmt.Println("Failed to find role:", err)
			return
		}
	}
	err = Db.Where("name = ?", "follow").First(&typepost).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			typepost = TypeofPost{Name: "follow"}
			if err := Db.Create(&typepost).Error; err != nil {
				fmt.Println("Failed to create role:", err)
				return
			}
		} else {
			fmt.Println("Failed to find role:", err)
			return
		}
	}
}
func AccessDefault() {
	var access Access
	err := Db.Where("name = ?", "public").First(&access).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			access = Access{Name: "public"}
			if err := Db.Create(&access).Error; err != nil {
				fmt.Println("Failed to create role:", err)
				return
			}
		} else {
			fmt.Println("Failed to find role:", err)
			return
		}
	}
	err = Db.Where("name = ?", "follow").First(&access).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			access = Access{Name: "follow"}
			if err := Db.Create(&access).Error; err != nil {
				fmt.Println("Failed to create role:", err)
				return
			}
		} else {
			fmt.Println("Failed to find role:", err)
			return
		}
	}
}

func AccessAndTypePostDefault() {
	TypeofPostDefault()
	AccessDefault()
}
