package controllers

import (
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
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err := a.transaksiService.Store(&newTransaksi)
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success create new transaksi",
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
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	var newTransaksi models.Transaksi

	if err := c.ShouldBindJSON(&newTransaksi); err != nil {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	newTransaksi.ID = id

	err = a.transaksiService.Update(id, newTransaksi)
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success update transaksi",
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
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err = a.transaksiService.Delete(id)

	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success delete transaksi",
	})
}

func (a *transaksiAPI) GetList(c *gin.Context) {

	transaksi, err := a.transaksiService.GetList()
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, transaksi)
}