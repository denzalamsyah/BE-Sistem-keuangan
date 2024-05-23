package controllers

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/services"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

type KasAPI interface {
	Store(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetList(c *gin.Context)
	Search(c *gin.Context)
	DownloadPembayaranKas(c *gin.Context)
	DownloadPengambilanKas(c *gin.Context)
}

type kasAPI struct {
	kasService services.KasServices
}

func NewKasAPI(kasRepo services.KasServices) *kasAPI {
	return &kasAPI{
		kasService: kasRepo,
	}
}

func (a *kasAPI) Store(c *gin.Context) {

	var newKas models.KasGuru

	if err := c.ShouldBindJSON(&newKas); err != nil {
		log.Printf("Error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err := a.kasService.Store(&newKas)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil menambah data",
		"data" : newKas,
	})
}

func (a *kasAPI) Update(c *gin.Context) {

	kasID := c.Param("id")

	if kasID == "" {
		
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(kasID)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	var newKas models.KasGuru

	if err := c.ShouldBindJSON(&newKas); err != nil {
		log.Printf("Error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	newKas.ID = id

	err = a.kasService.Update(id, newKas)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil mengubah data",
		"data" : newKas,
	})
}

func (a *kasAPI) Delete(c *gin.Context) {

	kasID := c.Param("id")

	if kasID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(kasID)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err = a.kasService.Delete(id)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(500, gin.H{
			"message" : "internal server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil menghapus data",
	})
}

func (a *kasAPI) GetList(c *gin.Context) {
    page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
        page = 1
    }

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 10
    }
	kasList, totalPage, err := a.kasService.GetList(page, pageSize)
	if err != nil {
		log.Printf("Error: %v", err)
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
        "data": kasList,
        "meta": meta,
    }

    c.JSON(200, response)
}

func (s *kasAPI) Search(c *gin.Context){
	nama := c.Query("nama")
	tanggal := c.Query("tanggal")


	kas, err := s.kasService.Search(nama, tanggal)
	if err != nil {
        log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kas})
}

func (s *kasAPI) DownloadPembayaranKas(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	log.Printf("ID: %v", id)
	
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	result, err := s.kasService.GetByID(id)
	log.Printf("data: %v", result)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error",
            "error" : err.Error(),
		})
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.SetHeaderFunc(func() {
		y := pdf.GetY()
		imageWidth := 20
		imageHeight := 20

		xImage := 10
		xText := xImage + imageWidth + 5

		pdf.Image("./app/files/logo.png", float64(xImage), y, float64(imageWidth), float64(imageHeight), false, "", 0, "")

		pdf.SetX(float64(xText))
		pdf.SetFont("Times", "B", 12)
		pdf.Cell(0, 0, "SMA Plus Nurul Iman Leles")
		pdf.SetFont("Times", "", 10)
		pdf.Ln(2)
		pdf.SetX(float64(xText))
		pdf.Cell(0, 10, "Kp. Galumpit Kidul RT 005/RW 004 Des. Cipancar Kec. Leles Garut Jawa Barat")
		pdf.Ln(5)
		pdf.SetX(float64(xText))
		tanggalSekarang := time.Now().Format("02 January 2006")
		pdf.CellFormat(0, 10, "Garut, "+tanggalSekarang, "0", 1, "", false, 0, "")
		pdf.SetX(float64(xText))
		pdf.Cell(0, 0, "No. Telp: 123456789")
		pdf.Ln(5)

		xStart := 10
		xEnd := 200
		pdf.Line(float64(xStart), 33, float64(xEnd), 33)
		pdf.Ln(5)
	})

	leftMargin := 10.0
	rightMargin := 110.0
	lineHeight := 6.0

    pdf.AddPage()
	pdf.SetFont("Times", "BU", 11)
	pdf.CellFormat(0, 8, "KWITANSI PEMBAYARAN KAS", "0", 1, "C", false, 0, "")
    pdf.Ln(0)
    pdf.SetFont("Times", "", 10)
	pdf.SetX(leftMargin)
    pdf.CellFormat(0, lineHeight, "Nama: "+result.Nama, "0", 1, "", false, 0, "")
	pdf.SetX(leftMargin)
	pdf.CellFormat(0, lineHeight, "NIP : "+result.NIP, "0", 1, "", false, 0, "")
	pdf.SetX(leftMargin)
    pdf.CellFormat(0, lineHeight, "Jumlah Bayar: Rp. "+formatNumber(int(result.Jumlah_Bayar)), "0", 1, "", false, 0, "")
	pdf.SetX(leftMargin)
    pdf.CellFormat(0, lineHeight, "Tanggal Pembayaran: "+result.Tanggal, "0", 1, "", false, 0, "")
	pdf.SetX(leftMargin)
    pdf.Ln(6)

	pdf.SetY(pdf.GetY() - 5*lineHeight)
	pdf.SetX(rightMargin)
	pdf.CellFormat(0, lineHeight, "Waktu Update : "+result.UpdatedAt, "0", 1, "", false, 0, "")
    pdf.Ln(4)
	pdf.SetFont("Times", "I", 10)
	pdf.SetX(rightMargin)
	pdf.CellFormat(0, lineHeight, "Mengetahui,", "", 1, "", false, 0, "")
	pdf.SetFont("Times", "B", 10)
	pdf.SetX(rightMargin)
	pdf.CellFormat(0, 25, "(Bendahara Sekolah)", "", 1, "", false, 0, "")
	
	xStart := 10
	xEnd := 200
	pdf.Line(float64(xStart), 80, float64(xEnd), 80)

    // Simpan file PDF
    fileName := "kas.pdf"
    err = pdf.OutputFileAndClose("./app/files/" + fileName)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    defer func() {
        if err := os.Remove("./app/files/" + fileName); err != nil {
            log.Printf("Gagal menghapus file: %v", err)
        }
    }()


    c.File("./app/files/" + fileName)

}

func (s *kasAPI) DownloadPengambilanKas(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	log.Printf("ID: %v", id)
	
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	result, err := s.kasService.GetAmbilByID(id)
	log.Printf("data: %v", result)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error",
            "error" : err.Error(),
		})
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.SetHeaderFunc(func() {
		y := pdf.GetY()
		imageWidth := 20
		imageHeight := 20

		xImage := 10
		xText := xImage + imageWidth + 5

		pdf.Image("./app/files/logo.png", float64(xImage), y, float64(imageWidth), float64(imageHeight), false, "", 0, "")

		pdf.SetX(float64(xText))
		pdf.SetFont("Times", "B", 12)
		pdf.Cell(0, 0, "SMA Plus Nurul Iman Leles")
		pdf.SetFont("Times", "", 10)
		pdf.Ln(2)
		pdf.SetX(float64(xText))
		pdf.Cell(0, 10, "Kp. Galumpit Kidul RT 005/RW 004 Des. Cipancar Kec. Leles Garut Jawa Barat")
		pdf.Ln(5)
		pdf.SetX(float64(xText))
		tanggalSekarang := time.Now().Format("02 January 2006")
		pdf.CellFormat(0, 10, "Garut, "+tanggalSekarang, "0", 1, "", false, 0, "")
		pdf.SetX(float64(xText))
		pdf.Cell(0, 0, "No. Telp: 123456789")
		pdf.Ln(5)

		xStart := 10
		xEnd := 200
		pdf.Line(float64(xStart), 33, float64(xEnd), 33)
		pdf.Ln(5)
	})

	leftMargin := 10.0
	rightMargin := 110.0
	lineHeight := 6.0
    pdf.AddPage()
	pdf.SetFont("Times", "BU", 11)
	pdf.CellFormat(0, 8, "KWITANSI PENGAMBILAN KAS", "0", 1, "C", false, 0, "")
	pdf.Ln(0)
    pdf.SetFont("Times", "", 10)
	pdf.SetX(leftMargin)
    pdf.CellFormat(0, lineHeight, "Nama : "+result.Nama, "0", 1, "", false, 0, "")
	pdf.SetX(leftMargin)
	pdf.CellFormat(0, lineHeight, "NIP : "+result.NIP, "0", 1, "", false, 0, "")
	pdf.SetX(leftMargin)
    pdf.CellFormat(0, lineHeight, "Jumlah Ambil : Rp. "+formatNumber(int(result.JumlahAmbil)), "0", 1, "", false, 0, "")
	pdf.SetX(leftMargin)
    pdf.CellFormat(0, lineHeight, "Tanggal Pengambilan : "+result.TanggalAmbil, "0", 1, "", false, 0, "")
    pdf.Ln(6)

	pdf.SetY(pdf.GetY() - 5*lineHeight)
	pdf.SetX(rightMargin)
	pdf.CellFormat(0, lineHeight, "Waktu Update : "+result.UpdatedAt, "0", 1, "", false, 0, "")
    pdf.Ln(4)
	pdf.SetFont("Times", "I", 10)
	pdf.SetX(rightMargin)
	pdf.CellFormat(0, lineHeight, "Mengetahui,", "", 1, "", false, 0, "")
	pdf.SetFont("Times", "B", 10)
	pdf.SetX(rightMargin)
	pdf.CellFormat(0, 25, "(Bendahara Sekolah)", "", 1, "", false, 0, "")
	
	xStart := 10
	xEnd := 200
	pdf.Line(float64(xStart), 80, float64(xEnd), 80)

    fileName := "kas.pdf"
    err = pdf.OutputFileAndClose("./app/files/" + fileName)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    defer func() {
        if err := os.Remove("./app/files/" + fileName); err != nil {
            log.Printf("Gagal menghapus file: %v", err)
        }
    }()


    c.File("./app/files/" + fileName)

}
