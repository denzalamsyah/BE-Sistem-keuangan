package controllers

import (
	"log"
	"net/http"
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
	GetTotalKelasCount(c *gin.Context)
	Search(c *gin.Context)
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
		log.Printf("Pesan error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err := a.kelasService.Store(&newKelas)
	if err != nil {
		log.Printf("Pesan error: %v", err)
		c.JSON(500, gin.H{
			"message" : err.Error(),
			"error": "Gagal menambah data"   ,
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil menambah data",
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
		log.Printf("Pesan error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	var newKelas models.Kelas

	if err := c.ShouldBindJSON(&newKelas); err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	newKelas.IDKelas = id

	err = a.kelasService.Update(id, newKelas)
	if err != nil {
		log.Printf("Pesan error: %v", err)
		c.JSON(500, gin.H{
			"message" : err.Error(),
			"error":  "Gagal mengubah data" ,
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil mengubah data",
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
		log.Printf("Pesan error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err = a.kelasService.Delete(id)
	if err != nil {
		log.Printf("Pesan error: %v", err)
		c.JSON(500, gin.H{
			"message" :  err.Error(),
			"error": "Gagal mengubah data" ,
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil menghapus data",
	
	})
}

func (a *kelasAPI) GetList(c *gin.Context) {

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
        page = 1
    }

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 30
    }

	kelas, totalPage, err := a.kelasService.GetList(page, pageSize)
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
        "data": kelas,
        "meta": meta,
    }

	c.JSON(200, response)
}

func (s *kelasAPI) GetTotalKelasCount(c *gin.Context){
	totalKelas, err := s.kelasService.GetTotalKelasCount()
	if err != nil {
		log.Printf("Pesan error: %v", err)
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
	}

	c.JSON(200,gin.H{
		"totalKelas": totalKelas,
	})

}

func (s *kelasAPI) Search(c *gin.Context){
	nama := c.Query("nama")

	kelas, err := s.kelasService.Search(nama)
	if err != nil {
        log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kelas})
}
