package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/services"
	"github.com/gin-gonic/gin"
)

type KasAPI interface {
	Store(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetList(c *gin.Context)
	Search(c *gin.Context)
}

type kasAPI struct {
	kasService services.KasServices
}

func NewKasAPI(kasRepo services.KasServices) *kasAPI {
	return &kasAPI{
		kasService: kasRepo,
	}
}

func (a *kasAPI) Store(c *gin.Context) {

	var newKas models.KasGuru

	if err := c.ShouldBindJSON(&newKas); err != nil {
		log.Printf("Error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err := a.kasService.Store(&newKas)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil menambah data",
		"data" : newKas,
	})
}

func (a *kasAPI) Update(c *gin.Context) {

	kasID := c.Param("id")

	if kasID == "" {
		
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(kasID)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	var newKas models.KasGuru

	if err := c.ShouldBindJSON(&newKas); err != nil {
		log.Printf("Error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	newKas.ID = id

	err = a.kasService.Update(id, newKas)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil mengubah data",
		"data" : newKas,
	})
}

func (a *kasAPI) Delete(c *gin.Context) {

	kasID := c.Param("id")

	if kasID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(kasID)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err = a.kasService.Delete(id)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil menghapus data",
	})
}

func (a *kasAPI) GetList(c *gin.Context) {
    page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
        page = 1
    }

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 10
    }
	kasList, totalPage, err := a.kasService.GetList(page, pageSize)
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
        "data": kasList,
        "meta": meta,
    }

    c.JSON(200, response)
}

func (s *kasAPI) Search(c *gin.Context){
	nama := c.Query("nama")
	tanggal := c.Query("tanggal")


	kas, err := s.kasService.Search(nama, tanggal)
	if err != nil {
        log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kas})
}

	