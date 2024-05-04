package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/services"
	"github.com/gin-gonic/gin"
)

type TransaksiAPI interface {
	AddTransaksi(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetList(c *gin.Context)
	Search(c *gin.Context)
}

type transaksiAPI struct {
	transaksiService services.TransaksiService
}

func NewTransaksiAPI(transaksiRepo services.TransaksiService) *transaksiAPI {
	return &transaksiAPI{transaksiRepo}
}

func (a *transaksiAPI) AddTransaksi(c *gin.Context) {

	var newTransaksi models.Transaksi

	if err := c.ShouldBindJSON(&newTransaksi); err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err := a.transaksiService.Store(&newTransaksi)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"error" : err.Error(),
			"message":   "Gagal menambah data",
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil menambah data",
		"data" : newTransaksi,
	})
}

func (a *transaksiAPI) Update(c *gin.Context) {

	transaksiID := c.Param("id")
	
	if transaksiID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(transaksiID)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	var newTransaksi models.Transaksi

	if err := c.ShouldBindJSON(&newTransaksi); err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	newTransaksi.ID = id

	err = a.transaksiService.Update(id, newTransaksi)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"error" : err.Error(),
			"message":   "Gagal mengubah data",
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil mengubah data",
		"data" : newTransaksi,
	})

}

func (a *transaksiAPI) Delete(c *gin.Context) {

	transaksiID := c.Param("id")

	if transaksiID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(transaksiID)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err = a.transaksiService.Delete(id)

	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"error" : err.Error(),
			"message":   "Gagal menghapus data",
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil menghapus data ",
	})
}

func (a *transaksiAPI) GetList(c *gin.Context) {

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
        page = 1
    }

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 100
    }

	transaksi, totalPage, err := a.transaksiService.GetList(page, pageSize)
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
        "data": transaksi,
        "meta": meta,
    }

	c.JSON(200, response)
}

func (a *transaksiAPI) Search(c *gin.Context){
	nama := c.Query("nama")

	transaksi, err := a.transaksiService.Search(nama)
	if err != nil {
        log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transaksi})
}