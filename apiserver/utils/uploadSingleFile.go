package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"os"
	"path/filepath"
)

func SaveSingleFile(c *gin.Context, file *multipart.FileHeader, storePath, suffix string) (string, error) {

	fmt.Println(storePath)
	if !Exists(storePath) {
		err := os.MkdirAll(storePath, 0644)
		if err != nil {
			Logger.Error(err)
			BadResponse(c, "上传文件失败，权限受限")
			return "", err
		}
	}

	newFileName := fmt.Sprintf("%s%s", GenUuid(), suffix)

	err := c.SaveUploadedFile(file, filepath.Join(storePath, newFileName))
	if err != nil {
		Logger.Error(err)
		BadResponse(c, "上传文件失败")
		return "", err
	}

	return newFileName, nil
}
