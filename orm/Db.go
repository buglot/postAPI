package orm

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB
var err error

func InitDB() {
	dsn := os.Getenv("DSN_MYSQL")
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	Db.AutoMigrate(&User{}, &Post{}, &Image{}, &Follow{}, &LikePost{}, &Comment{})
}
