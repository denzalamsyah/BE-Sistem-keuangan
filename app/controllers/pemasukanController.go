package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/services"
	"github.com/gin-gonic/gin"
)

type PemasukanAPI interface {
	FindAll(c *gin.Context)
	TotalKeuangan(c *gin.Context)
	Add(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
	GetList(c *gin.Context)
	SearchAll(c *gin.Context)
	Search(c *gin.Context)
}

type pemasukanAPI struct {
	pemasukanService services.PemasukanService
}

func NewPemasukanAPI(pemasukanRepo services.PemasukanService) *pemasukanAPI {
	return &pemasukanAPI{pemasukanRepo}
}


func (s *pemasukanAPI) FindAll(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
        page = 1
    }

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 1000
    }

	result,totalPage, err := s.pemasukanService.FindAll(page, pageSize)
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

func (s *pemasukanAPI) Add(c *gin.Context) {
	var newPemasukan models.Pemasukanlainnya

	if err := c.ShouldBindJSON(&newPemasukan); err != nil {
		log.Printf("Pesan error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}
	err := s.pemasukanService.Store(&newPemasukan)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" :err.Error(),
			"error": "Gagal menambah data"  ,
		})
		return
	}
	c.JSON(200, gin.H{
		"message" : "Berhasil menambah data",
		"data" : newPemasukan,
	})

}

func (s *pemasukanAPI) Update(c *gin.Context) {
	pemasukanID := c.Param("id")

	if pemasukanID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(pemasukanID)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	var newPemasukan models.Pemasukanlainnya

	if err := c.ShouldBindJSON(&newPemasukan); err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	newPemasukan.ID = id

	err = s.pemasukanService.Update(id, newPemasukan)
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
		"data" : newPemasukan,
	})
}

func (s *pemasukanAPI) Delete(c *gin.Context) {
	pemasukanID := c.Param("id")

	if pemasukanID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(pemasukanID)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err = s.pemasukanService.Delete(id)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" :    err.Error(),
			"error":   "Gagal menghapus data",
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil menghapus data",
	})
}

func (s *pemasukanAPI) GetByID(c *gin.Context) {
	pemasukanID := c.Param("id")

	if pemasukanID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(pemasukanID)

	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	result, err := s.pemasukanService.GetByID(id)

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

func (s *pemasukanAPI) GetList(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
        page = 1
    }

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 1000
    }
	result, totalPage, err := s.pemasukanService.GetList(page, pageSize)
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

func (s *pemasukanAPI) TotalKeuangan(c *gin.Context) {
	saldo, pengeluaran, pemasukan, err := s.pemasukanService.TotalKeuangan()
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	response := models.Total{
		Saldo:    saldo,
		Pengeluaran: pengeluaran,
		Pemasukan: pemasukan,
	}
	c.JSON(200, response)
}

func (s *pemasukanAPI) SearchAll(c *gin.Context){
	nama := c.Query("nama")
	tanggal := c.Query("tanggal")

	pemasukan, err := s.pemasukanService.SearchAll(nama, tanggal)
	if err != nil {
        log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": pemasukan})

}

func (s *pemasukanAPI) Search(c *gin.Context){
	nama := c.Query("nama")
	tanggal := c.Query("tanggal")

	pemasukan, err := s.pemasukanService.Search(nama, tanggal)
	if err != nil {
        log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": pemasukan})

}
