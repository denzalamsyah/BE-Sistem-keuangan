package controllers

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/services"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

type SemesterAPI interface {
	AddSemester(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
	GetList(c *gin.Context)
	GetListByCategory(c *gin.Context)
	Search(c *gin.Context)
	DownloadPembayaranSiswa(c *gin.Context)
	GetLunasByNIP(ctx *gin.Context)
	DownloadReportSiswa(c *gin.Context)
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
			"message" : err.Error(),
			"error":   "Periksa kembali inputan anda",
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil menambah data",
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
			"message" : err.Error(),
			"error":   "Periksa kembali inputan anda",
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil mengubah data",
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
			"error" : err.Error(),
			"message":   "Gagal menghapus data",
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil menhapus data",
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
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
        page = 1
    }

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 10000
    }
	result, totalPage, err := s.semesterService.GetList(page, pageSize)
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
        "data": result,
        "meta": meta,
    }

    c.JSON(200, response)
}

func (s *semesterAPI) GetListByCategory(c *gin.Context){
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
        page = 1
    }

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 10000
    }

	category := c.Query("kategori")
	result, totalPage, err := s.semesterService.GetListByCategory(page, pageSize, category)
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
        "data": result,
        "meta": meta,
    }

    c.JSON(200, response)
}

func (s *semesterAPI) Search(c *gin.Context){
	siswa := c.Query("siswa")
	tahunAjar := c.Query("tahunAjar")
	transaksi := c.Query("transaksi")
	semester := c.Query("semester")
	tanggal := c.Query("tanggal")
	nisn := c.Query("nisn")
	kategori := c.Query("kategori")


	// penerima := c.Query("penerima")

	pembayaran, err := s.semesterService.Search(siswa, tahunAjar, transaksi, semester, tanggal, nisn, kategori)
	if err != nil {
        log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": pembayaran})

}

func (s *semesterAPI) DownloadPembayaranSiswa(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	result, err := s.semesterService.GetByID(id)
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
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(0, 0, "SMA Plus Nurul Iman Leles",)
		pdf.SetFont("Arial", "", 10)
		pdf.Ln(2)
		// pdf.SetX(float64(xText))
		// pdf.Cell(0, 10, "SMA Plus Nurul Iman Leles")
		// pdf.Ln(5)
		pdf.SetX(float64(xText))
		pdf.Cell(0, 10, "Kp. Galumpit Kidul RT 005/RW 004 Des. Cipancar Kec. Leles Garut Jawa Barat")
		pdf.Ln(5)
		pdf.SetX(float64(xText))
		pdf.CellFormat(0, 10, "Garut, "+result.Kelas, "0", 1, "", false, 0, "")
		// pdf.Ln(5)
		pdf.SetX(float64(xText))
		pdf.Cell(0, 0, "No. Telp: 123456789")
		pdf.Ln(5)

		xStart := 10 
		xEnd := 200 
		pdf.Line(float64(xStart), 33, float64(xEnd), 33)
		pdf.Ln(10)

    })

    pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 0, "Kwetansi Pembayaran", "0", 1, "C", false, 0, "")
    pdf.Ln(10)
    pdf.SetFont("Arial", "B", 12)
    pdf.CellFormat(0, 10, "Nama : "+result.Siswa, "0", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, "NISN : "+strconv.Itoa(int(result.NISN)), "0", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, "Nama Pembayaran : "+result.Transaksi, "0", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, "Semester : "+result.Bulan, "0", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, "Kelas : "+result.Kelas, "0", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, "Jumlah Bayar : Rp. "+strconv.Itoa(result.Jumlah), "0", 1, "", false, 0, "")
    pdf.CellFormat(0, 10, "Tanggal Pembayaran : "+result.Tanggal, "0", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, "Status : "+result.Status, "0", 1, "", false, 0, "")
    pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
    pdf.CellFormat(0, 10, "Bendahara Sekolah,", "0", 1, "", false, 0, "")
    pdf.Ln(20)

	pdf.SetFont("Arial", "I", 12)
    pdf.CellFormat(0, 10, "Sambas S.Pd.", "0", 1, "", false, 0, "")
    // Simpan file PDF
    fileName := "semester.pdf"
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

func (c *semesterAPI) GetLunasByNIP(ctx *gin.Context) {
	nisn, err := strconv.Atoi(ctx.Param("nisn"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid NIP"})
		return
	}

	pembayaran, err := c.semesterService.GetLunasByNISN(nisn)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get data"})
		return
	}

	ctx.JSON(http.StatusOK, pembayaran)
}

func (s *semesterAPI) DownloadReportSiswa(c *gin.Context) {
	nisn, err := strconv.Atoi(c.Param("nisn"))

	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	result, err := s.semesterService.GetLunasByNISN(nisn)
	log.Printf("data: %v", result)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error",
            "error" : err.Error(),
		})
		return
	}

	// file :=excelize.NewFile()

}

