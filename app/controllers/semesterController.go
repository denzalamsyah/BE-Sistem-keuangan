package controllers

import (
	"log"
	"strconv"

	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/services"
	"github.com/gin-gonic/gin"
)

type SemesterAPI interface {
	AddSemester(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
	GetList(c *gin.Context)
}

type semesterAPI struct{
	semesterService services.SemesterServices
}

func NewSemesterAPI(semesterRepo services.SemesterServices) *semesterAPI{
	return &semesterAPI{semesterRepo}
}

func (s *semesterAPI) AddSemester(c *gin.Context) {
	var newPembayaranSemester models.PembayaranSemester

	if err := c.ShouldBindJSON(&newPembayaranSemester); err != nil{
		log.Printf("Error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err := s.semesterService.Store(&newPembayaranSemester)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success create new PembayaranSemester",
		"data" : newPembayaranSemester,
	})
}

func (s *semesterAPI) Update(c *gin.Context) {

	PembayaranSemesterID := c.Param("id")
	
	if PembayaranSemesterID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(PembayaranSemesterID)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}
	var newPembayaranSemester models.PembayaranSemester

	if err := c.ShouldBindJSON(&newPembayaranSemester); err != nil{
		log.Printf("Error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	newPembayaranSemester.ID = id
	err = s.semesterService.Update(id, newPembayaranSemester)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success update PembayaranSemester",
		"data" : newPembayaranSemester,
	})
	
}
func (s *semesterAPI) Delete(c *gin.Context) {

	PembayaranSemesterID := c.Param("id")
	
	if PembayaranSemesterID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(PembayaranSemesterID)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err = s.semesterService.Delete(id)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success delete PembayaranSemester",
	})
	
}
func (s *semesterAPI) GetByID(c *gin.Context) {

	PembayaranSemesterID, err := strconv.Atoi(c.Param("id"))
	
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	result, err := s.semesterService.GetByID(PembayaranSemesterID)	
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, result)
	
}
func (s *semesterAPI) GetList(c *gin.Context) {
	result, err := s.semesterService.GetList()
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, result)
}
