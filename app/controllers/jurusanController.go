package controllers

import (
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
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	err := a.jurusanService.Store(&newJurusan)
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
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
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	var newJurusan models.Jurusan

	if err := c.ShouldBindJSON(&newJurusan); err != nil {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	newJurusan.IDJurusan = id

	err = a.jurusanService.Update(id, newJurusan)
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
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
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	err = a.jurusanService.Delete(id)
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success delete jurusan",
	})
}

func (a *jurusanAPI) GetList(c *gin.Context) {

	jurusanList, err := a.jurusanService.GetList()
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
		})
		return
	}

	c.JSON(200, jurusanList)
}
