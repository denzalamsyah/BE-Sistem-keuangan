package controllers

import (
	"log"
	"strconv"

	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/services"
	"github.com/gin-gonic/gin"
)

type SppAPI interface {
	AddSPP(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
	GetList(c *gin.Context)
}

type sppAPI struct{
	sppService services.SPPServices
}

func NewSPPAPI(sppRepo services.SPPServices) *sppAPI{
	return &sppAPI{sppRepo}
}

func (s *sppAPI) AddSPP(c *gin.Context) {
	var newPembayaranSPP models.PembayaranSPP

	if err := c.ShouldBindJSON(&newPembayaranSPP); err != nil{
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err := s.sppService.Store(&newPembayaranSPP)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success create new PembayaranSPP",
		"data" : newPembayaranSPP,
	})
}

func (s *sppAPI) Update(c *gin.Context) {

	PembayaranSPPID := c.Param("id")
	
	if PembayaranSPPID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(PembayaranSPPID)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}
	var newPembayaranSPP models.PembayaranSPP

	if err := c.ShouldBindJSON(&newPembayaranSPP); err != nil{
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	newPembayaranSPP.ID = id
	err = s.sppService.Update(id, newPembayaranSPP)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success update PembayaranSPP",
		"data" : newPembayaranSPP,
	})
	
}
func (s *sppAPI) Delete(c *gin.Context) {

	PembayaranSPPID := c.Param("id")
	
	if PembayaranSPPID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			
		})
		return
	}

	id, err := strconv.Atoi(PembayaranSPPID)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err = s.sppService.Delete(id)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success delete PembayaranSPP",
	})
	
}
func (s *sppAPI) GetByID(c *gin.Context) {

	PembayaranSPPID, err := strconv.Atoi(c.Param("id"))
	
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	result, err := s.sppService.GetByID(PembayaranSPPID)	
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
func (s *sppAPI) GetList(c *gin.Context) {
	result, err := s.sppService.GetList()
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