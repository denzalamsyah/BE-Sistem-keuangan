package controllers

import (
	"log"
	"strconv"

	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/services"
	"github.com/gin-gonic/gin"
)

type JurusanAPI interface {
	AddJurusan(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetList(c *gin.Context)
	GetTotalJurusanCount(c *gin.Context)
}

type jurusanAPI struct {
	jurusanService services.JurusanService
}

func NewJurusanAPI(jurusanRepo services.JurusanService) *jurusanAPI {
	return &jurusanAPI{
		jurusanService: jurusanRepo,
	}
}

func (a *jurusanAPI) AddJurusan(c *gin.Context) {

	var newJurusan models.Jurusan

	if err := c.ShouldBindJSON(&newJurusan); err != nil {
		log.Printf("Error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err := a.jurusanService.Store(&newJurusan)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success add jurusan",
		"data" : newJurusan,
	})
}

func (a *jurusanAPI) Update(c *gin.Context) {

	jurusanID := c.Param("id")

	if jurusanID == "" {
		
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(jurusanID)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	var newJurusan models.Jurusan

	if err := c.ShouldBindJSON(&newJurusan); err != nil {
		log.Printf("Error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	newJurusan.IDJurusan = id

	err = a.jurusanService.Update(id, newJurusan)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success update jurusan",
		"data" : newJurusan,
	})
}

func (a *jurusanAPI) Delete(c *gin.Context) {

	jurusanID := c.Param("id")

	if jurusanID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(jurusanID)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err = a.jurusanService.Delete(id)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success delete jurusan",
	})
}

func (a *jurusanAPI) GetList(c *gin.Context) {
    page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
        page = 1
    }

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 10
    }
	jurusanList, totalPage, err := a.jurusanService.GetList(page, pageSize)
	if err != nil {
		log.Printf("Error: %v", err)
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
        "data": jurusanList,
        "meta": meta,
    }

    c.JSON(200, response)
}


func (s *jurusanAPI) GetTotalJurusanCount(c *gin.Context) {
	totalJurusan, err := s.jurusanService.GetTotalJurusanCount()
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"totalJurusan": totalJurusan,
	})
}

	