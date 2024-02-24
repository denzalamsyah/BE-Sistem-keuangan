package middleware

import (
	"context"
	"log"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/denzalamsyah/simak/app/initializers"
	"github.com/gin-gonic/gin"
)

func FileUploadMiddleware() gin.HandlerFunc{
	return func (ctx *gin.Context) {
		file, header, err := ctx.Request.FormFile("file")
		if err != nil {
			ctx.AbortWithStatusJSON(400, gin.H{
				"message" : "invalid request body",
			})
			return
		}
		defer file.Close()

		ctx.Set("filePateh", header.Filename)
		ctx.Set("file", file)

		ctx.Next()
	}
}


func UploadToCloudinary(file *multipart.FileHeader) (string, error) {
    // Open the file
    fileContent, err := file.Open()
    if err != nil {
        return "", err
    }
    defer fileContent.Close()

    ctx := context.Background()
    cld, err := initializers.SetUpCloudinary()
    if err != nil {
        return "", err
    }
    uploadParams := uploader.UploadParams{}
    result, err := cld.Upload.Upload(ctx, fileContent, uploadParams)
    if err != nil {
        return "", err
    }

	log.Printf("Uploaded image to Cloudinary. URL: %s",result.SecureURL)
    return result.SecureURL, nil
}