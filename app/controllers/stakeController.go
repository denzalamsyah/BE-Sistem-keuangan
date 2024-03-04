package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/denzalamsyah/simak/app/middleware"
	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/services"
	"github.com/gin-gonic/gin"
)

type StakeAPI interface {
	AddStake(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
	GetList(c *gin.Context)
}

type stakeAPI struct{
	stakeService services.StakeholderServices
}

func NewStakeAPI(stakeRepo services.StakeholderServices) *stakeAPI{
	return &stakeAPI{stakeRepo}
}

func (s *stakeAPI) AddStake(c *gin.Context){
	var stake models.Stakeholder

	if err := c.ShouldBind(&stake); err != nil{
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	file, err := c.FormFile("file")
    if err != nil {
        // Jika tidak ada file yang diunggah, gunakan gambar default
        stake.Gambar = "https://res.cloudinary.com/dgvkpzi4p/image/upload/v1706338009/149071_fxemnm.png"
    } else {
        // Jika ada file yang diunggah, upload ke Cloudinary dan dapatkan URL-nya
        imageURL, err := middleware.UploadToCloudinary(file)
        if err != nil {
		log.Printf("Pesan error: %v", err)

            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "failed to upload image to Cloudinary",
                "error":   err.Error(),
            })
            return
        }
        stake.Gambar = imageURL
    }

	err = s.stakeService.Store(&stake)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success create new stakeholder",
		"data" : stake,
	})	

}

func (s *stakeAPI) Update(c *gin.Context){

	stakeID := c.Param("id")

	if stakeID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",

		})
		return
	}
	
	id, err := strconv.Atoi(stakeID)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	var newStake models.Stakeholder

	if err := c.ShouldBind(&newStake); err != nil{
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// newStake.ID = id
	file, err := c.FormFile("file")
    if err != nil && err != http.ErrMissingFile {
		log.Printf("Pesan error: %v", err)

        c.JSON(http.StatusBadRequest, gin.H{
            "message": "failed to get image from form-data",
            "error":   err.Error(),
        })
        return
    }

	if file != nil {
		imageURL, err := middleware.UploadToCloudinary(file)
		if err != nil {
		log.Printf("Pesan error: %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to upload image to Cloudinary",
				"error":   err.Error(),
			})
			return
		}
		newStake.Gambar = imageURL
	}

	err = s.stakeService.Update(id, newStake)

	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success update stakeholder",
		"data" : newStake,
	})
}

func (s *stakeAPI) Delete(c *gin.Context){
	stakeID := c.Param("id")
	
	if stakeID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(stakeID)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err = s.stakeService.Delete(id)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success delete stake",
	})
}

func (s *stakeAPI) GetByID(c *gin.Context){
	stakeID, err := strconv.Atoi(c.Param("id"))
	if stakeID == 0 {
		c.JSON(400, gin.H{
			"message" : "data notfound",
			"error":   err.Error(),
		})
		return
	}
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	result, err := s.stakeService.GetByID(stakeID)	
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, result)
}

func (s *stakeAPI) GetList(c *gin.Context){
	page, err := strconv.Atoi(c.Query("page"))
    if err != nil || page <= 0 {
        page = 1
    }

    pageSize, err := strconv.Atoi(c.Query("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 100
    }

	result, totalPage, err := s.stakeService.GetList(page, pageSize)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}
	meta := gin.H{
        "current_page": page,
        "total_pages":  totalPage,
    }

    response := gin.H{
        "data": result,
        "meta": meta,
    }

    c.JSON(200, response)
}