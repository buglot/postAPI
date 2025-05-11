package post

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/buglot/postAPI/lib"
	"github.com/buglot/postAPI/orm"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Uploads(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}
	ext := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + ext
	dst := filepath.Join("uploads", newFileName)
	err = ctx.SaveUploadedFile(file, dst)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Upload successful",
		"filename": newFileName,
		"url":      newFileName,
	})
}

type EditUp struct {
	Url string `json:"url"`
}

func UploadProfile(ctx *gin.Context) {
	userid := lib.AnyToUInt(ctx.MustGet("userID"))
	file, err := ctx.FormFile("image")
	jsonData := ctx.PostForm("data")
	var input EditUp
	if err := json.Unmarshal([]byte(jsonData), &input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON: " + err.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}
	ext := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + ext
	dst := filepath.Join("img", "profile", newFileName)
	err = ctx.SaveUploadedFile(file, dst)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file"})
		return
	}
	var user orm.User
	orm.Db.Where("id = ?", userid).First(&user)
	user.Avatar = input.Url + "/img/public/" + newFileName
	orm.Db.Save(&user)
	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Upload successful",
		"filename": newFileName,
		"url":      newFileName,
	})
}
