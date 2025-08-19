package service

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/gin-gonic/gin"
)

func UpdateProfile(file *multipart.FileHeader, id interface{}, c *gin.Context) (*string, error) {
	var uploadDir string
	if file != nil {
		result, err := DeleteProfile(id)
		if err != nil {
			return nil, err
		}

		if result == nil {
			return nil, nil
		}

		uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
		uploadDir = "/assets/img/profile/" + uniqueName
		imgPath := filepath.Join(os.Getenv("PROFILE_UPLOAD_PATH"), uploadDir)
		if err := c.SaveUploadedFile(file, imgPath); err != nil {
			return nil, err
		}
	}

	return &uploadDir, nil
}

func DeleteProfile(id interface{}) (*string, error) {
	filePath, err := model.GetUserImagePath(id)

	if err != nil {
		return nil, err
	}

	if filePath == nil {
		return nil, nil
	}

	var defaultPath = "/assets/img/profile2.jpg"

	if *filePath == defaultPath {
		var is_default = "default"
		return &is_default, nil
	}

	fullPath := filepath.Join(os.Getenv("PROFILE_UPLOAD_PATH"), *filePath)
	if err := os.Remove(fullPath); err != nil {
		return nil, err
	}

	return &fullPath, nil
}
