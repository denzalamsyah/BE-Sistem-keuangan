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

type SiswaAPI interface {
	AddSiswa(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
	GetList(c *gin.Context)
	History(c *gin.Context)
    GetTotalGenderCount(c *gin.Context)
}

type siswaAPI struct {
	siswaService services.SiswaServices
}

func NewSiswaAPI(siswaRepo services.SiswaServices) *siswaAPI {
	return &siswaAPI{siswaRepo}
}

func (s *siswaAPI) AddSiswa(c *gin.Context) {
    var newSiswa models.Siswa

    err := c.ShouldBind(&newSiswa)
    if err != nil {
		log.Printf("Error: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "invalid request body",
            "error":   err.Error(),
        })
        return
    }

    file, err := c.FormFile("file")
    if err != nil {
        // Jika tidak ada file yang diunggah, gunakan gambar default
        newSiswa.Gambar = "https://res.cloudinary.com/dgvkpzi4p/image/upload/v1706338009/149071_fxemnm.png"
    } else {
        // Jika ada file yang diunggah, upload ke Cloudinary dan dapatkan URL-nya
        imageURL, err := middleware.UploadToCloudinary(file)
        if err != nil {
			log.Printf("Error: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "failed to upload image to Cloudinary",
                "error":   err.Error(),
            })
            return
        }
        newSiswa.Gambar = imageURL
    }
    err = s.siswaService.Store(&newSiswa)
    if err != nil {
		log.Printf("Error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "message": "internal server error",
            "error":   err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "success create new siswa",
        "data":    newSiswa,
    })
}

func (s *siswaAPI) Update(c *gin.Context) {
    siswaID := c.Param("id")

    if siswaID == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "invalid request body",
        })
        return
    }

    id, err := strconv.Atoi(siswaID)
    if err != nil {
		log.Printf("Pesan error: %v", err)

        c.JSON(http.StatusBadRequest, gin.H{
            "message": "invalid request body",
        })
        return
    }

    var existingSiswa models.Siswa
    if err := c.ShouldBind(&existingSiswa); err != nil {
		log.Printf("Pesan error: %v", err)

        c.JSON(http.StatusBadRequest, gin.H{
            "message": "invalid request body",
            "error":   err.Error(),
        })
        return
    }

    // Jika ada file yang diunggah, perbarui gambar siswa
    file, err := c.FormFile("file")
    if err != nil && err != http.ErrMissingFile {
        
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
        existingSiswa.Gambar = imageURL
    }

    // Lakukan pembaruan data siswa di database
    err = s.siswaService.Update(id, existingSiswa)
    if err != nil {
		log.Printf("Pesan error: %v", err)

        c.JSON(http.StatusInternalServerError, gin.H{
            "message": "internal server error",
            "error":   err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "success update siswa",
        "data":    existingSiswa,
    })
}

func (s *siswaAPI) Delete(c *gin.Context) {

	siswaID := c.Param("id")
	
	if siswaID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(siswaID)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	err = s.siswaService.Delete(id)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success delete siswa",
	})
	
}
func (s *siswaAPI) GetByID(c *gin.Context) {

	siswaID, err := strconv.Atoi(c.Param("id"))
	
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	result, err := s.siswaService.GetByID(siswaID)	
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error",
		})
		return
	}

	c.JSON(200, result)
	
}
func (s *siswaAPI) GetList(c *gin.Context) {
    page, err := strconv.Atoi(c.Query("page"))
    if err != nil || page <= 0 {
        page = 1
    }

    pageSize, err := strconv.Atoi(c.Query("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 5000
    }

    result, totalPage, err := s.siswaService.GetList(page, pageSize)
    if err != nil {
		log.Printf("Pesan error: %v", err)

        c.JSON(500, gin.H{
            "message": "internal server error",
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


func (s *siswaAPI) History(c *gin.Context) {

	siswaID, err := strconv.Atoi(c.Param("id"))
	
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	result, err := s.siswaService.HistoryPembayaranSiswa(siswaID)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error",
		})
		return
	}

	c.JSON(200, result)
}

func (api *siswaAPI) GetTotalGenderCount(c *gin.Context) {
    countLakiLaki, countPerempuan, err := api.siswaService.GetTotalGenderCount()
    if err != nil {
        log.Printf("Error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "count_laki_laki": countLakiLaki,
        "count_perempuan": countPerempuan,
    })
}
