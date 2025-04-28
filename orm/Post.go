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
	LikePost     []LikePost `gorm:"foreignKey:PostID"`
	Comment      []Comment  `gorm:"foreignKey:PostID"`
}
type LikePost struct {
	gorm.Model
	UserID uint
	User   User `gorm:"foreignKey:UserID"`
	PostID uint
}
type Access struct {
	gorm.Model
	Name string `gorm:"type:enum('public','follow','private');default:'public';unique;not null"`
}
type TypeofPost struct {
	gorm.Model
	Name string `gorm:"type:enum('daily','shop');default:'daily';not null;unique"`
}
type Image struct {
	gorm.Model
	Url    string `gorm:"type:varchar(255);not null"`
	PostID uint
}
type Comment struct {
	gorm.Model
	UserID  uint
	User    User `gorm:"foreignKey:UserID"`
	PostID  uint
	Comment string `gorm:"not null"`
}

func TypeofPostDefault() {
	names := []string{"daily", "shop"}

	for _, name := range names {
		var typepost TypeofPost
		err := Db.Where("name = ?", name).First(&typepost).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				newType := TypeofPost{Name: name}
				if err := Db.Create(&newType).Error; err != nil {
					fmt.Printf("Failed to create typeofpost '%s': %v\n", name, err)
				} else {
					fmt.Printf("Created typeofpost: %s\n", name)
				}
			} else {
				fmt.Printf("Failed to find typeofpost '%s': %v\n", name, err)
			}
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
