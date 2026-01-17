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

func UpdateProfile(file *multipart.FileHeader, id interface{}, state string, c *gin.Context) (*string, error) {
	var uploadDir string
	if file != nil {
		result, err := DeleteProfile(id, state)
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

func DeleteProfile(id interface{}, state string) (*string, error) {
	var filePath *string
	var err error
	if state == "user" {
		filePath, err = model.GetUserImagePath(id)
	} else {
		filePath, err = model.GetSellerImagePath(id)
	}

	if err != nil {
		return nil, err
	}

	if filePath == nil {
		return nil, nil
	}

	var defaultPath string
	if state == "user" {
		defaultPath = "/assets/img/profile2.jpg"
	} else {
		defaultPath = "/assets/img/profile1.jpg"
	}

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

func UpdateImgProduct(files []*multipart.FileHeader, id interface{}, productId int, c *gin.Context) ([]string, error) {
	var uploadDir []string
	if files != nil {
		var uniqueName string
		var imgPath string

		for i, val := range files {
			uniqueName = fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(val.Filename))
			uploadDir = append(uploadDir, "/assets/img/"+uniqueName)
			imgPath = filepath.Join(os.Getenv("PRODUCT_IMG_UPLOAD_PATH"), uploadDir[i])
			if err := c.SaveUploadedFile(val, imgPath); err != nil {
				return nil, err
			}
		}
	}

	return uploadDir, nil
}

func DeleteImgProduct(id interface{}, productId int) ([]string, error) {
	filePath, err := model.GetProductImagePath(id, productId)

	if err != nil {
		return nil, err
	}

	if filePath == nil {
		return nil, nil
	}

	var fullPath []string

	for i, val := range filePath {
		fullPath = append(fullPath, filepath.Join(os.Getenv("PRODUCT_IMG_UPLOAD_PATH"), val))
		if err := os.Remove(fullPath[i]); err != nil {
			return nil, err
		}
	}

	return fullPath, nil
}
