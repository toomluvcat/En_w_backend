package handler

import (
	"Render/app/conect"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

func UploadToCld(ctx *gin.Context, file multipart.File, fileName string) (string, error) {
	PublicID := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	uploadResult, err := conect.CLD.Upload.Upload(ctx.Request.Context(), file, uploader.UploadParams{
		Folder:   "my_items",
		PublicID: PublicID})
	fmt.Printf("result: %s \n",uploadResult)
	if err != nil {
		return "", err
	}
	return uploadResult.SecureURL, nil
}

func DeleteCld(ctx *gin.Context, url string) error {
	publicID := "my_items/"+strings.TrimSuffix(filepath.Base(url),filepath.Ext(url))
	_, err := conect.CLD.Upload.Destroy(ctx.Request.Context(), uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete Cloudinary image: %w", err)
	}
	return nil
}
