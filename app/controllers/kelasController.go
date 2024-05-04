package controllers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

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
		
		if strings.Contains(err.Error(), "foreign key constraint") {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Tidak bisa mengubah kode kelas",
				"error":   "Gagal mengubah data",
			})
			return
		}

		existingKelas, err := a.kelasService.GetKode(newKelas.KodeKelas)
        if err != nil {
            log.Printf("Error checking Kelas: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   err.Error(),
                "message":   "Gagal menambah data",
            })
            return
        }
        
        if existingKelas.KodeKelas == newKelas.KodeKelas{
            c.JSON(http.StatusBadRequest, gin.H{
                "message":   "Kode kelas sudah ada",
                "error":   "Gagal menambah data",
            })
            return
        }
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil menambah data",
		"data" : newKelas,
	})
}

func (a *kelasAPI) Update(c *gin.Context) {
	kodeKelas := c.Param("kode")
    if kodeKelas == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "invalid request body" ,
			"error" : "kosong",
        })
        return
    }

    kode, err := strconv.Atoi(kodeKelas)
    if err != nil {
		log.Printf("Pesan error: %v", err)

        c.JSON(http.StatusBadRequest, gin.H{
            "message": "invalid request body",
            "error" : err.Error(),
        })
        return
    }

    var newKelas models.Kelas
    if err := c.ShouldBind(&newKelas); err != nil {
		log.Printf("Encode error: %v", err)

        c.JSON(http.StatusBadRequest, gin.H{
            "message": "invalid request body",
            "error" : err.Error(),
        })
        return
    }

    err = a.kelasService.Update(kode, newKelas)
    if err != nil {
		log.Printf("Update error: %v", err)

        existingKelas, err := a.kelasService.GetKode(newKelas.KodeKelas)
        if err != nil {
            log.Printf("Error checking NISN: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   err.Error(),
            })
            return
        }
        
        if existingKelas.KodeKelas == newKelas.KodeKelas{
            c.JSON(http.StatusBadRequest, gin.H{
                "message":   "Kode kelas sudah ada",
                "error":   "Gagal mengubah data",
            })
            return
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Berhasil mengubah Kelas",
        "data":    newKelas,
    })
}

func (a *kelasAPI) Delete(c *gin.Context) {

	kelasKode := c.Param("kode")

	if kelasKode == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	kode, err := strconv.Atoi(kelasKode)
	if err != nil {
		log.Printf("Pesan error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err = a.kelasService.Delete(kode)
	if err != nil {
		log.Printf("Pesan error: %v", err)
		c.JSON(500, gin.H{
			"error" :  err.Error(),
			"message": "Gagal mengubah data" ,
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
