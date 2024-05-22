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

type PengeluaranAPI interface {
	Add(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
	GetList(c *gin.Context)
	Search(c *gin.Context)
	DownloadLaporan(c *gin.Context)
}

type pengeluaranAPI struct {
	pengeluaranService services.PengeluaranService
}

func NewPengeluaranAPI(pengeluaranRepo services.PengeluaranService) *pengeluaranAPI {
	return &pengeluaranAPI{pengeluaranRepo}
}

func (s *pengeluaranAPI) Add(c *gin.Context) {
	var newPengeluaran models.Pengeluaran

	if err := c.ShouldBindJSON(&newPengeluaran); err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}
	err := s.pengeluaranService.Store(&newPengeluaran)
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
		"data" : newPengeluaran,
	})

}

func (s *pengeluaranAPI) Update(c *gin.Context) {
	PengeluaranID := c.Param("id")

	if PengeluaranID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(PengeluaranID)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	var newPengeluaran models.Pengeluaran

	if err := c.ShouldBindJSON(&newPengeluaran); err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	newPengeluaran.ID = id

	err = s.pengeluaranService.Update(id, newPengeluaran)
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
		"data" : newPengeluaran,
	})
}

func (s *pengeluaranAPI) Delete(c *gin.Context) {
	PengeluaranID := c.Param("id")

	if PengeluaranID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(PengeluaranID)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err = s.pengeluaranService.Delete(id)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"error" : err.Error(),
			"message":   "Gagal menghapus data",
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil menghapus data",
	})
}

func (s *pengeluaranAPI) GetByID(c *gin.Context) {
	PengeluaranID := c.Param("id")

	if PengeluaranID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
			
		})
		return
	}

	id, err := strconv.Atoi(PengeluaranID)

	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	result, err := s.pengeluaranService.GetByID(id)

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

func (s *pengeluaranAPI) GetList(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
        page = 1
    }

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 1000
    }
	result, totalPage, err := s.pengeluaranService.GetList(page, pageSize)
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

func (s *pengeluaranAPI) Search(c *gin.Context){
	nama := c.Query("nama")
	tanggal := c.Query("tanggal")

	pengeluaran, total, err := s.pengeluaranService.Search(nama, tanggal)
	if err != nil {
        log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": pengeluaran, "total":total})
}

func (s *pengeluaranAPI) DownloadLaporan(c *gin.Context){
	tanggal := c.Query("tanggal")
	nama := c.Query("nama")

	pemasukan, total, err := s.pengeluaranService.Search(nama, tanggal)
	if err != nil {
        log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.SetHeaderFunc(func() {
		// Simpan posisi Y saat ini untuk digunakan pada akhir teks
		y := pdf.GetY()
		// Hitung lebar gambar
		imageWidth := 20 // Ubah sesuai kebutuhan
		imageHeight := 20 // Ubah sesuai kebutuhan
		// Tentukan posisi X untuk gambar dan teks
		xImage := 10
		xText := xImage + imageWidth + 5
		// Gambar di sebelah kiri teks
		pdf.Image("./app/files/logo.png", float64(xImage), y, float64(imageWidth), float64(imageHeight), false, "", 0, "")
	
		pdf.SetX(float64(xText))
		pdf.SetFont("Times", "B", 14)
		pdf.Cell(0, 0, "SMA Plus Nurul Iman Leles",)
		pdf.SetFont("Times", "", 12)
		pdf.Ln(2)
		// pdf.SetX(float64(xText))
		// pdf.Cell(0, 10, "SMA Plus Nurul Iman Leles")
		// pdf.Ln(5)
		pdf.SetX(float64(xText))
		pdf.Cell(0, 10, "Kp. Galumpit Kidul RT 005/RW 004 Des. Cipancar Kec. Leles Garut Jawa Barat")
		pdf.Ln(5)
		pdf.SetX(float64(xText))
		tanggalSekarang := time.Now().Format("02 January 2006") 
		pdf.CellFormat(0, 10, "Garut, "+tanggalSekarang, "0", 1, "", false, 0, "")
		// pdf.Ln(5)
		pdf.SetX(float64(xText))
		pdf.Cell(0, 0, "No. Telp: 123456789")
		pdf.Ln(5)

		xStart := 10 
		xEnd := 200 
		pdf.Line(float64(xStart), 33, float64(xEnd), 33)
		pdf.Ln(5)

    })


	pdf.AddPage()
	pdf.SetFont("Times", "B", 12)
	pdf.CellFormat(0, 10, "Laporan Pengeluaran : "+pemasukan[0].Tanggal, "0", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, "Jumlah Pemasukan : Rp. "+formatNumber(int(total)), "0", 1, "", false, 0, "")
	pdf.Ln(5)
	pdf.SetFont("Times", "BU", 12)
	pdf.CellFormat(0, 10, "RINCIAN BIAYA PENGELUARAN", "0", 1, "C", false, 0, "")
	pdf.Ln(3)

		xStart := 10 
		xEnd := 200 
		tableWidth := xEnd - xStart
		numColumns := 3
		columnWidth := float64(tableWidth) / float64(numColumns)
		// pdf.Ln(-1)
		pdf.SetFont("Times", "B", 12)
		pdf.CellFormat(columnWidth, 10, "Nama Transaksi", "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth, 10, "Tanggal", "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidth, 10, "Jumlah", "1", 0, "C", false, 0, "")
		pdf.Ln(-1) // Pindah ke baris baru
		
		for _, data := range pemasukan {
			pdf.SetFont("Times", "", 12)
			pdf.CellFormat(columnWidth, 8, data.Nama, "1", 0, "C", false, 0, "")
			pdf.CellFormat(columnWidth, 8, data.Tanggal, "1", 0, "C", false, 0, "")
			pdf.CellFormat(columnWidth, 8, "Rp. " +formatNumber(int(data.Jumlah)), "1", 0, "C", false, 0, "")
			pdf.Ln(-1)	
		}
	pdf.Ln(15)
	pdf.SetFont("Times", "I", 11)
	pdf.SetX(float64(150))
	pdf.CellFormat(0, 2, "Mengetahui,", "0", 1, "", false, 0, "")
	pdf.Ln(25)
	pdf.SetFont("Times", "BU", 11)
	pdf.SetX(float64(150))
    pdf.CellFormat(0, 3, "Kepala Sekolah", "0", 1, "", false, 0, "")

		

    fileName := "pemasukan.pdf"
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