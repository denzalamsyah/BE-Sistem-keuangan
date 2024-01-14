package controllers

import (
	"strconv"

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
}

type siswaAPI struct {
	siswaService services.SiswaServices
}

func NewSiswaAPI(siswaRepo services.SiswaServices) *siswaAPI {
	return &siswaAPI{siswaRepo}
}

func (s *siswaAPI) AddSiswa(c *gin.Context) {
	var newSiswa models.Siswa

	if err := c.ShouldBindJSON(&newSiswa); err != nil{
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	err := s.siswaService.Store(&newSiswa)
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success create new siswa",
		"data" : newSiswa,
	})
}

func (s *siswaAPI) Update(c *gin.Context) {

	siswaID := c.Param("id")
	
	if siswaID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(siswaID)
	if err != nil {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}
	var newSiswa models.Siswa

	if err := c.ShouldBindJSON(&newSiswa); err != nil{
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	newSiswa.ID = id
	err = s.siswaService.Update(id, newSiswa)
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success update siswa",
		"data" : newSiswa,
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
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	err = s.siswaService.Delete(id)
	if err != nil {
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
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	result, err := s.siswaService.GetByID(siswaID)	
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
		})
		return
	}

	c.JSON(200, result)
	
}
func (s *siswaAPI) GetList(c *gin.Context) {
	result, err := s.siswaService.GetList()
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
		})
		return
	}
	c.JSON(200, result)
}