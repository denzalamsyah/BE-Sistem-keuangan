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

	pageWidth := 65.0
	pageHeight := 75.0

	// Membuat instance PDF baru dengan ukuran halaman kustom
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "mm",
		Size: gofpdf.SizeType{
			Wd: pageWidth,
			Ht: pageHeight,
		},
	})
	topMargin := 5.0
	pdf.SetTopMargin(topMargin)
	pdf.SetAutoPageBreak(true, topMargin)
	pdf.SetHeaderFunc(func() {
		y := pdf.GetY()
		imageWidth := 8 
		imageHeight := 8 

		xImage := 6
		xText := xImage + imageWidth + 2
		pdf.Image("./app/files/logo.png", float64(xImage), y, float64(imageWidth), float64(imageHeight), false, "", 0, "")
	
		pdf.SetX(float64(xText))
		pdf.SetFont("Arial", "B", 7)
		pdf.Cell(0, 0, "SMA Plus Nurul Iman Leles",)
		pdf.SetFont("Arial", "", 5)
		pdf.Ln(2)
		pdf.SetX(float64(xText))
		pdf.Cell(0, 2, "Kp. Galumpit Kidul RT 005/RW 004 Cipancar Leles")
		pdf.Ln(2)
		pdf.SetX(float64(xText))
		tanggalSekarang := time.Now().Format("02 January 2006") 
		pdf.CellFormat(0, 2, "Garut, "+tanggalSekarang, "0", 1, "", false, 0, "")
		pdf.SetX(float64(xText))
		pdf.Cell(0, 2, "No. Telp: 123456789")
		pdf.Ln(5)

		xStart := 5 
		xEnd := 60.0
		pdf.Line(float64(xStart), 14, float64(xEnd), 14)
		pdf.Ln(3)

    })

    pdf.AddPage()
	pdf.SetFont("Arial", "B", 6)
	pdf.CellFormat(0, 0, "Kwitansi Pembayaran", "0", 1, "C", false, 0, "")
    pdf.Ln(3)
    pdf.SetFont("Arial", "B", 6)
	pdf.SetX(float64(6))
    pdf.CellFormat(0, 4, "Nama: "+result.Nama, "0", 1, "", false, 0, "")
	pdf.SetX(float64(6))
	pdf.CellFormat(0, 4, "NIP : "+result.NIP, "0", 1, "", false, 0, "")
	pdf.SetX(float64(6))
    pdf.CellFormat(0, 4, "Jumlah Bayar: Rp. "+formatNumber(int(result.Jumlah_Bayar)), "0", 1, "", false, 0, "")
	pdf.SetX(float64(6))
    pdf.CellFormat(0, 4, "Tanggal Pembayaran: "+result.Tanggal, "0", 1, "", false, 0, "")
	pdf.SetX(float64(6))
	pdf.CellFormat(0, 4, "Waktu Update : "+result.UpdatedAt, "0", 1, "", false, 0, "")
    pdf.Ln(4)

	pdf.SetFont("Arial", "I", 6)
	pdf.SetX(float64(6))
	pdf.CellFormat(0, 2, "Mengetahui,", "0", 1, "", false, 0, "")
	pdf.SetFont("Arial", "B", 6)
	pdf.SetX(float64(6))
    pdf.CellFormat(0, 3, "*Bendahara Sekolah*", "0", 1, "", false, 0, "")
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

	pageWidth := 65.0
	pageHeight := 75.0

	// Membuat instance PDF baru dengan ukuran halaman kustom
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "mm",
		Size: gofpdf.SizeType{
			Wd: pageWidth,
			Ht: pageHeight,
		},
	})
	topMargin := 5.0
	pdf.SetTopMargin(topMargin)
	pdf.SetAutoPageBreak(true, topMargin)

	pdf.SetHeaderFunc(func() {
	
		y := pdf.GetY()
		imageWidth := 8 
		imageHeight := 8 
		xImage := 6
		xText := xImage + imageWidth + 2
		// Gambar di sebelah kiri teks
		pdf.Image("./app/files/logo.png", float64(xImage), y, float64(imageWidth), float64(imageHeight), false, "", 0, "")
	
		pdf.SetX(float64(xText))
		pdf.SetFont("Arial", "B", 7)
		pdf.Cell(0, 0, "SMA Plus Nurul Iman Leles",)
		pdf.SetFont("Arial", "", 5)
		pdf.Ln(2)
		pdf.SetX(float64(xText))
		pdf.Cell(0, 2, "Kp. Galumpit Kidul RT 005/RW 004 Cipancar Leles")
		pdf.Ln(2)
		pdf.SetX(float64(xText))
		tanggalSekarang := time.Now().Format("02 January 2006") 
		pdf.CellFormat(0, 2, "Garut, "+tanggalSekarang, "0", 1, "", false, 0, "")
		// pdf.Ln(5)
		pdf.SetX(float64(xText))
		pdf.Cell(0, 2, "No. Telp: 123456789")
		pdf.Ln(5)

		xStart := 5 
		xEnd := 60.0
		pdf.Line(float64(xStart), 14, float64(xEnd), 14)
		pdf.Ln(3)

    })

    pdf.AddPage()
	pdf.SetFont("Arial", "B", 6)
	pdf.CellFormat(0, 0, "Kwitansi Pembayaran", "0", 1, "C", false, 0, "")
    pdf.Ln(3)
    pdf.SetFont("Arial", "B", 6)
	pdf.SetX(float64(6))
    pdf.CellFormat(0, 4, "Nama : "+result.Nama, "0", 1, "", false, 0, "")
	pdf.SetX(float64(6))
	pdf.CellFormat(0, 4, "NIP : "+result.NIP, "0", 1, "", false, 0, "")
	pdf.SetX(float64(6))
    pdf.CellFormat(0, 4, "Jumlah Ambil : Rp. "+formatNumber(int(result.JumlahAmbil)), "0", 1, "", false, 0, "")
	pdf.SetX(float64(6))
    pdf.CellFormat(0, 4, "Tanggal Pengambilan : "+result.TanggalAmbil, "0", 1, "", false, 0, "")
	pdf.SetX(float64(6))
	pdf.CellFormat(0, 4, "Waktu Update : "+result.UpdatedAt, "0", 1, "", false, 0, "")
    pdf.Ln(4)

	pdf.SetFont("Arial", "I", 6)
	pdf.SetX(float64(6))
	pdf.CellFormat(0, 2, "Mengetahui,", "0", 1, "", false, 0, "")
	pdf.SetFont("Arial", "B", 6)
	pdf.SetX(float64(6))
    pdf.CellFormat(0, 3, "*Bendahara Sekolah*", "0", 1, "", false, 0, "")

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
