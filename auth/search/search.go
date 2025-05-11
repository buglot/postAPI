package search

import (
	"errors"
	"net/http"

	"github.com/buglot/postAPI/orm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetNameAndUrl(ctx *gin.Context) {
	search := ctx.Query("name")
	pattern := "%" + search + "%"

	var users []orm.User
	err := orm.Db.
		Where("url LIKE ? OR username LIKE ? OR email LIKE ?", pattern, pattern, pattern).
		Find(&users).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNoContent, gin.H{
			"message": "Not Found :" + search,
		})
		return
	} else if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, users)
}
