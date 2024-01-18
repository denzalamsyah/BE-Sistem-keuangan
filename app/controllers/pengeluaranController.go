package controllers

import (
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
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}
	err := s.pengeluaranService.Store(&newPengeluaran)
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
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
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	var newPengeluaran models.Pengeluaran

	if err := c.ShouldBindJSON(&newPengeluaran); err != nil {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	newPengeluaran.ID = id

	err = s.pengeluaranService.Update(id, newPengeluaran)
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
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
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	err = s.pengeluaranService.Delete(id)
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
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
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	result, err := s.pengeluaranService.GetByID(id)

	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
		})
		return
	}

	c.JSON(200, result)
}

func (s *pengeluaranAPI) GetList(c *gin.Context) {
	result, err := s.pengeluaranService.GetList()
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
		})
		return
	}
	c.JSON(200, result)
}
