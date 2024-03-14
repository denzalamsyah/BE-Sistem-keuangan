package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/services"
	"github.com/gin-gonic/gin"
)

type PengeluaranAPI interface {
	Add(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
	GetList(c *gin.Context)
	Search(c *gin.Context)
}

type pengeluaranAPI struct {
	pengeluaranService services.PengeluaranService
}

func NewPengeluaranAPI(pengeluaranRepo services.PengeluaranService) *pengeluaranAPI {
	return &pengeluaranAPI{pengeluaranRepo}
}

func (s *pengeluaranAPI) Add(c *gin.Context) {
	var newPengeluaran models.Pengeluaran

	if err := c.ShouldBindJSON(&newPengeluaran); err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}
	err := s.pengeluaranService.Store(&newPengeluaran)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message" : "success create new Pengeluaran",
		"data" : newPengeluaran,
	})

}

func (s *pengeluaranAPI) Update(c *gin.Context) {
	PengeluaranID := c.Param("id")

	if PengeluaranID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(PengeluaranID)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	var newPengeluaran models.Pengeluaran

	if err := c.ShouldBindJSON(&newPengeluaran); err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	newPengeluaran.ID = id

	err = s.pengeluaranService.Update(id, newPengeluaran)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success update Pengeluaran",
		"data" : newPengeluaran,
	})
}

func (s *pengeluaranAPI) Delete(c *gin.Context) {
	PengeluaranID := c.Param("id")

	if PengeluaranID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(PengeluaranID)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err = s.pengeluaranService.Delete(id)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success delete Pengeluaran",
	})
}

func (s *pengeluaranAPI) GetByID(c *gin.Context) {
	PengeluaranID := c.Param("id")

	if PengeluaranID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			
		})
		return
	}

	id, err := strconv.Atoi(PengeluaranID)

	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	result, err := s.pengeluaranService.GetByID(id)

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

func (s *pengeluaranAPI) GetList(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
        page = 1
    }

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 1000
    }
	result, totalPage, err := s.pengeluaranService.GetList(page, pageSize)
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

func (s *pengeluaranAPI) Search(c *gin.Context){
	nama := c.Query("nama")
	tanggal := c.Query("tanggal")

	pengeluaran, err := s.pengeluaranService.Search(nama, tanggal)
	if err != nil {
        log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": pengeluaran})
}