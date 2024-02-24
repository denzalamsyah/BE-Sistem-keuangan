package initializers

import (
	"github.com/cloudinary/cloudinary-go/v2"
)

func  SetUpCloudinary() (*cloudinary.Cloudinary, error){
	cld, err := cloudinary.NewFromURL("cloudinary://366452479239627:G8gEsIiIeD5Y22h09yZ-uJ7tIrM@dgvkpzi4p")
	if err != nil {
		return nil, err
	}

	return cld, nil
}