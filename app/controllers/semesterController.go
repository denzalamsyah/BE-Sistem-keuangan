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
		tanggalSekarang := time.Now().Format("02 January 2006") 
		pdf.CellFormat(0, 10, "Garut, "+tanggalSekarang, "0", 1, "", false, 0, "")
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
	pdf.CellFormat(0, 10, "NISN : "+result.NISN, "0", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, "Nama Pembayaran : "+result.Transaksi, "0", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, "Semester : "+result.Bulan, "0", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, "Kelas : "+result.Semester, "0", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, "Jumlah Bayar : Rp. "+strconv.Itoa(result.Jumlah), "0", 1, "", false, 0, "")
    pdf.CellFormat(0, 10, "Tanggal Pembayaran : "+result.Tanggal, "0", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, "Status : "+result.Status, "0", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, "Waktu Update : "+result.UpdatedAt, "0", 1, "", false, 0, "")
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

// download report pembayaran berdasarkan nisn siswa
func (s *semesterAPI) DownloadReportSiswa(c *gin.Context) {
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
    pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "RINCIAN BIAYA", "0", 1, "C", false, 0, "")
    pdf.Ln(5)

	pdf.SetFont("Arial", "B", 12)

	// Tampilkan data siswa
	pdf.CellFormat(0, 10, "Nama : "+result[0].Siswa, "0", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, "NISN : "+result[0].NISN, "0", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, "TIPE : "+kategori, "0", 1, "", false, 0, "")
	// pdf.Ln(10)

		xStart := 10 
		xEnd := 200 
		tableWidth := xEnd - xStart

		// Jumlah kolom dalam tabel
		numColumns := 4

		// Menghitung lebar setiap sel
		columnWidth := float64(tableWidth) / float64(numColumns)

		
		printedTitle := false
		for _, data := range result {
			if data.Transaksi == "SPP" {
				if !printedTitle {
					// Buat tabel khusus untuk pembayaran SPP
					// pdf.CellFormat(columnWidth*4, 7, "Data Pembayaran Siswa", "1", 0, "C", false, 0, "")
					pdf.Ln(-1)
					pdf.CellFormat(columnWidth, 7, "Nama", "1", 0, "C", false, 0, "")
					pdf.CellFormat(columnWidth, 7, "Bulan/Semester", "1", 0, "C", false, 0, "")
					// pdf.CellFormat(columnWidth, 7, "Biaya", "1", 0, "C", false, 0, "")
					pdf.CellFormat(columnWidth, 7, "Bayar", "1", 0, "C", false, 0, "")
					// pdf.CellFormat(columnWidth, 7, "Tunggakan", "1", 0, "C", false, 0, "")
					pdf.CellFormat(columnWidth, 7, "Status", "1", 0, "C", false, 0, "")
					pdf.Ln(-1)
					printedTitle = true
				}
			} else {
				if !printedTitle {
					// Buat tabel untuk pembayaran selain SPP
					// pdf.CellFormat(columnWidth*6, 7, "Data Pembayaran", "1", 0, "C", false, 0, "")
					pdf.Ln(-1)
					pdf.CellFormat(columnWidth, 7, "Nama Pembayaran", "1", 0, "C", false, 0, "")
					pdf.CellFormat(columnWidth, 7, "Semester", "1", 0, "C", false, 0, "")
					// pdf.CellFormat(columnWidth, 7, "Biaya", "1", 0, "C", false, 0, "")
					pdf.CellFormat(columnWidth, 7, "Bayar", "1", 0, "C", false, 0, "")
					// pdf.CellFormat(columnWidth, 7, "Tunggakan", "1", 0, "C", false, 0, "")
					pdf.CellFormat(columnWidth, 7, "Status", "1", 0, "C", false, 0, "")
					pdf.Ln(-1)
					printedTitle = true
				}
			}

			// Set font untuk data
			pdf.SetFont("Arial", "", 12)

			// Menampilkan data
			pdf.CellFormat(columnWidth, 7, data.Transaksi, "1", 0, "C", false, 0, "")
			if data.Transaksi == "SPP" {
				// Hilangkan kolom semester jika pembayaran SPP
				pdf.CellFormat(columnWidth, 7, data.Bulan, "1", 0, "C", false, 0, "")
			} else {
				pdf.CellFormat(columnWidth, 7, data.Semester, "1", 0, "C", false, 0, "")
			}
			// pdf.CellFormat(columnWidth, 7, strconv.Itoa(int(data.Biaya)), "1", 0, "C", false, 0, "")
			pdf.CellFormat(columnWidth, 7, "Rp. " +strconv.Itoa(int(data.Jumlah)), "1", 0, "C", false, 0, "")
			// pdf.CellFormat(columnWidth, 7, strconv.Itoa(int(data.Tunggakan)), "1", 0, "C", false, 0, "")
			pdf.CellFormat(columnWidth, 7, data.Status, "1", 0, "C", false, 0, "")
			pdf.Ln(-1)
		}

		

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

