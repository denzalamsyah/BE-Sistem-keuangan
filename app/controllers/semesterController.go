package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
	DownloadReportSiswa(c *gin.Context)
}

type semesterAPI struct{
	semesterService services.SemesterServices
}

func NewSemesterAPI(semesterRepo services.SemesterServices) *semesterAPI{
	return &semesterAPI{semesterRepo}
}

// menambah pembayaran siswa
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

// update pembayaran siswa
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

// menghapus pembayaran siswa berdasarkan id
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

// mendapatkan detail pembayaran siswa
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

// mendapatkan seluruh data pembayaran
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

// mendapatkan list data pembayaran berdasarkan kategori transaksi
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

// mendapatkan data dengan form pencarian
func (s *semesterAPI) Search(c *gin.Context){
	siswa := c.Query("siswa")
	transaksi := c.Query("transaksi")
	semester := c.Query("semester")
	tanggal := c.Query("tanggal")
	nisn := c.Query("nisn")
	kategori := c.Query("kategori")


	// penerima := c.Query("penerima")

	pembayaran, err := s.semesterService.Search(siswa, transaksi, semester, tanggal, nisn, kategori)
	if err != nil {
        log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": pembayaran})

}

// download kwetansi pembayaran berdasarkan id
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

	// pdf := gofpdf.New("P", "mm", "A4", "")
	// pdf.SetPage(7)
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
    pdf.CellFormat(0, 4, "Nama : "+result.Siswa, "0", 1, "", false, 0, "")
	pdf.SetX(float64(6))
	pdf.CellFormat(0, 4, "NISN : "+result.NISN, "0", 1, "", false, 0, "")
	pdf.SetX(float64(6))
	pdf.CellFormat(0, 4, "Nama Pembayaran : "+result.Transaksi, "0", 1, "", false, 0, "")
	pdf.SetX(float64(6))
	pdf.CellFormat(0, 4, "Semester : "+result.Semester, "0", 1, "", false, 0, "")
	pdf.SetX(float64(6))
	pdf.CellFormat(0, 4, "Bulan : "+result.Bulan, "0", 1, "", false, 0, "")
	pdf.SetX(float64(6))
	pdf.CellFormat(0, 4, "Jumlah Bayar : Rp. "+formatNumber(int(result.Jumlah)), "0", 1, "", false, 0, "")
	pdf.SetX(float64(6))
    pdf.CellFormat(0, 4, "Tanggal Pembayaran : "+result.Tanggal, "0", 1, "", false, 0, "")
	pdf.SetX(float64(6))
	pdf.CellFormat(0, 4, "Status : "+result.Status, "0", 1, "", false, 0, "")
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

// download report pembayaran berdasarkan nisn siswa
func (s *semesterAPI) DownloadReportSiswa(c *gin.Context) {
	nama := c.Query("nama")
	kategories := c.Query("kategories")

	trans, err := s.semesterService.SearchTransaksi(nama, kategories)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	siswa := c.Query("siswa")
	transaksi := c.Query("transaksi")
	bulan := c.Query("bulan")
	tanggal := c.Query("tanggal")
	nisn := c.Query("nisn")
	kategori := c.Query("kategori")

	result, err := s.semesterService.Search(siswa, transaksi, bulan, tanggal, nisn, kategori)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resultMap := make(map[int][]models.PembayaranSemesterResponse)
	for _, r := range result {
		resultMap[r.TransaksiID] = append(resultMap[r.TransaksiID], r)
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
		pdf.SetFont("Times", "B", 14)
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

	pdf.AddPage()
	pdf.SetFont("Times", "BU", 12)
	pdf.CellFormat(0, 10, "RINCIAN BIAYA", "0", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Times", "B", 11)

	// Tampilkan data siswa
	pdf.CellFormat(0, 10, "Nama : "+result[0].Siswa, "0", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, "NISN : "+result[0].NISN, "0", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, "TIPE : "+kategories, "0", 1, "", false, 0, "")
	pdf.Ln(5)


	pdf.SetFont("Times", "B", 11)

	header := []string{"NO.", "Pembayaran", "Bulan", "Semester", "Jumlah","Status"}
	widths := []float64{10, 40, 30, 30,40, 40}

	for i, str := range header {
		pdf.CellFormat(widths[i], 15, str, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	totalAmount := 0
	totalPaid := 0

	pdf.SetFont("Times", "", 11)
	rowNum := 1
	for _, t := range trans {
		if t.Nama == "SPP" {
			totalAmount += t.Jumlah * 12
		} else if t.Nama == "LKS" {
			totalAmount += t.Jumlah * 2
		} else {
			totalAmount += t.Jumlah
		}

		if details, exists := resultMap[t.ID]; exists {
			for _, detail := range details {
				pdf.CellFormat(widths[0], 10, fmt.Sprintf("%d", rowNum), "1", 0, "C", false, 0, "")
				pdf.CellFormat(widths[1], 10, t.Nama, "1", 0, "C", false, 0, "")
				pdf.CellFormat(widths[2], 10, detail.Bulan, "1", 0, "C", false, 0, "")
				pdf.CellFormat(widths[3], 10, detail.Semester, "1", 0, "C", false, 0, "")
				pdf.CellFormat(widths[4], 10, "Rp. "+formatNumber(int(detail.Jumlah)), "1", 0, "C", false, 0, "")
				pdf.CellFormat(widths[5], 10, detail.Status, "1", 0, "C", false, 0, "")
				totalPaid += detail.Jumlah
				pdf.Ln(-1)
				rowNum++
			}
		} else {
			// Tampilkan transaksi tanpa detail secara umum
			pdf.CellFormat(widths[0], 10, fmt.Sprintf("%d", rowNum), "1", 0, "C", false, 0, "")
			pdf.CellFormat(widths[1], 10, t.Nama, "1", 0, "C", false, 0, "")
			pdf.CellFormat(widths[2], 10, "", "1", 0, "C", false, 0, "")
			pdf.CellFormat(widths[3], 10, "", "1", 0, "C", false, 0, "")
			pdf.CellFormat(widths[4], 10, "", "1", 0, "C", false, 0, "")
			pdf.CellFormat(widths[5], 10, "", "1", 0, "C", false, 0, "")
			pdf.Ln(-1)
			rowNum++
		}
	}

	// Calculate the outstanding amount
	outstanding := totalAmount - totalPaid

	// Add totals row
	pdf.SetFont("Times", "B", 11)
	pdf.Ln(10)
	pdf.CellFormat(0, 10, "Keterangan :", "0", 1, "L", false, 0, "")
	pdf.CellFormat(0, 10, fmt.Sprintf("Wajib Bayar: Rp. %s", formatNumber(totalAmount)), "0", 1, "L", false, 0, "")
	pdf.CellFormat(0, 10, fmt.Sprintf("Dibayar: Rp. %s", formatNumber(totalPaid)), "0", 1, "L", false, 0, "")
	pdf.CellFormat(0, 10, fmt.Sprintf("Tunggakan: Rp. %s",formatNumber(outstanding)), "0", 1, "L", false, 0, "")
	pdf.Ln(8)
	pdf.SetFont("Times", "I", 11)
	pdf.SetX(float64(150))
	pdf.CellFormat(0, 2, "Mengetahui,", "0", 1, "", false, 0, "")
	pdf.Ln(25)
	pdf.SetFont("Times", "B", 11)
	pdf.SetX(float64(150))
    pdf.CellFormat(0, 3, "*Bendahara Sekolah*", "0", 1, "", false, 0, "")

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

func formatNumber(number int) string {
    str := strconv.Itoa(number)
    if len(str) <= 3 {
        return str
    }

    var result []string
    for i := len(str); i > 0; i -= 3 {
        start := i - 3
        if start < 0 {
            start = 0
        }
        result = append([]string{str[start:i]}, result...)
    }

    return strings.Join(result, ".")
}


