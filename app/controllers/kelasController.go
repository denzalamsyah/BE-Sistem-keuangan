package controllers

import (
	"strconv"

	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/services"
	"github.com/gin-gonic/gin"
)

type KelasAPI interface {
	AddKelas(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetList(c *gin.Context)
}

type kelasAPI struct {
	kelasService services.KelasService
}

func NewKelasAPI(kelasRepo services.KelasService) *kelasAPI {
	return &kelasAPI{kelasRepo}
}

func (a *kelasAPI) AddKelas(c *gin.Context) {

	var newKelas models.Kelas

	if err := c.ShouldBindJSON(&newKelas); err != nil {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err := a.kelasService.Store(&newKelas)
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success create new kelas",
		"data" : newKelas,
	})
}

func (a *kelasAPI) Update(c *gin.Context) {

	kelasID := c.Param("id")

	if kelasID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(kelasID)
	if err != nil {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	var newKelas models.Kelas

	if err := c.ShouldBindJSON(&newKelas); err != nil {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	newKelas.IDKelas = id

	err = a.kelasService.Update(id, newKelas)
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success update kelas",
		"data" : newKelas,
	})
}

func (a *kelasAPI) Delete(c *gin.Context) {

	kelasID := c.Param("id")

	if kelasID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(kelasID)
	if err != nil {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err = a.kelasService.Delete(id)
	if err != nil {
		
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success delete kelas",
	
	})
}

func (a *kelasAPI) GetList(c *gin.Context) {

	kelas, err := a.kelasService.GetList()
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, kelas)
}