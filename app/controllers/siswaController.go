package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/denzalamsyah/simak/app/middleware"
	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/services"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

type SiswaAPI interface {
	AddSiswa(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
	GetList(c *gin.Context)
	History(c *gin.Context)
    GetTotalGenderCount(c *gin.Context)
    Search(c *gin.Context)
    SearchByKode(c *gin.Context)
    ExportSiswa(c *gin.Context)
    DownloadSiswa(c *gin.Context)
    BiodataSiswa(c *gin.Context)
    ImportFromExcel(c *gin.Context)
}

type siswaAPI struct {
	siswaService services.SiswaServices
}

func NewSiswaAPI(siswaRepo services.SiswaServices) *siswaAPI {
	return &siswaAPI{siswaRepo}
}

func (s *siswaAPI) AddSiswa(c *gin.Context) {
    var newSiswa models.Siswa

    err := c.ShouldBind(&newSiswa)
    if err != nil {
       
		log.Printf("Error: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "invalid request body",
            "error" : err.Error(),
        })
        return
    }

    file, err := c.FormFile("file")
    if err != nil {
        // Jika tidak ada file yang diunggah, gunakan gambar default
        newSiswa.Gambar = "https://res.cloudinary.com/dgvkpzi4p/image/upload/v1706338009/149071_fxemnm.png"
    } else {
        // Jika ada file yang diunggah, upload ke Cloudinary dan dapatkan URL-nya
        imageURL, err := middleware.UploadToCloudinary(file)
        if err != nil {
			log.Printf("Error: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "failed to upload image to Cloudinary",
                "error":   err.Error(),
            })
            return
        }
        newSiswa.Gambar = imageURL
    }

   

    err = s.siswaService.Store(&newSiswa)
    if err != nil {
        // if newSiswa.Alamat == "" || newSiswa.AgamaID == 0 || newSiswa.Angkatan == "" || newSiswa.Email == "" || newSiswa.GenderID == 0 || newSiswa.JurusanID == 0 || newSiswa.KelasID == 0 || newSiswa.Nisn == 0 || newSiswa.Nama == "" || newSiswa.NamaAyah == "" || newSiswa.NamaIbu == "" || newSiswa.NomorTelepon == 0 || newSiswa.TanggalLahir== "" || newSiswa.TempatLahir == ""{ 
        //     c.JSON(http.StatusBadRequest, gin.H{"error":   "semua item harus di isi kecuali gambar",})
		// 	return
		// }
        existingSiswa, err := s.siswaService.GetUserNisn(newSiswa.Nisn)
        if err != nil {
            log.Printf("Error checking NISN: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   err.Error(),
            })
            return
        }
        
        if existingSiswa.Nisn == newSiswa.Nisn{
            c.JSON(http.StatusBadRequest, gin.H{
                "message":   "Siswa dengan NISN tersebut sudah ada",
                "error":   "Gagal menambah data",
            })
            return
        }
        
        return
    }

   
    

    c.JSON(http.StatusOK, gin.H{
        "message": "Berhasil menambahkan data",
        "data":    newSiswa,
    })
}

func (s *siswaAPI) Update(c *gin.Context) {
    siswaNISN := c.Param("nisn")
    if siswaNISN == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "invalid request body" ,
        })
        return
    }

    // nisn, err := strconv.Atoi(siswaNISN)
    // if err != nil {
	// 	log.Printf("Pesan error: %v", err)

    //     c.JSON(http.StatusBadRequest, gin.H{
    //         "message": "invalid request body",
    //         "error" : err.Error(),
    //     })
    //     return
    // }

    var newSiswa models.Siswa
    if err := c.ShouldBind(&newSiswa); err != nil {
		log.Printf("Encode error: %v", err)

        c.JSON(http.StatusBadRequest, gin.H{
            "message": "invalid request body",
            "error" : err.Error(),
        })
        return
    }

    // Jika ada file yang diunggah, perbarui gambar siswa
    file, err := c.FormFile("file")
    if err != nil && err != http.ErrMissingFile {
        
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "failed to get image from form-data",
            "error":   err.Error(),
        })
        return
    }

    if file != nil {
        imageURL, err := middleware.UploadToCloudinary(file)
        if err != nil {
		log.Printf("middleware upload gambar error: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "failed to upload image to Cloudinary" + err.Error(),
                "error":   err.Error(),
            })
            return
        }
        newSiswa.Gambar = imageURL
    }


    err = s.siswaService.Update(siswaNISN, newSiswa)
    if err != nil {
		log.Printf("Update error: %v", err)

        existingSiswa, err := s.siswaService.GetUserNisn(newSiswa.Nisn)
        if err != nil {
            log.Printf("Error checking NISN: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   err.Error(),
            })
            return
        }
        
        if existingSiswa.Nisn == newSiswa.Nisn{
            c.JSON(http.StatusBadRequest, gin.H{
                "message":   "NISN yang anda masukan sudah ada",
                "error":   "Gagal mengubah data",
            })
            return
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Berhasil mengubah siswa",
        "data":    newSiswa,
    })
}

func (s *siswaAPI) Delete(c *gin.Context) {

	siswaNISN := c.Param("nisn")
	
	if siswaNISN == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	// nisn, err := strconv.Atoi(siswaNISN)
	// if err != nil {
	// 	log.Printf("Pesan error: %v", err)

	// 	c.JSON(400, gin.H{
	// 		"message" : "invalid request body" + err.Error(),
            
	// 	})
	// 	return
	// }

	err := s.siswaService.Delete(siswaNISN)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"error" : err.Error(),
            "message" : "Gagal menghapus data",
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil menghapus data",
	})
	
}
func (s *siswaAPI) GetByID(c *gin.Context) {

	siswaNISN := c.Param("nisn")
	
	// if err != nil {
	// 	log.Printf("Pesan error: %v", err)

	// 	c.JSON(400, gin.H{
	// 		"message" : "invalid request body",
    //         "error" : err.Error(),
	// 	})
	// 	return
	// }

	result, err := s.siswaService.GetByID(siswaNISN)	
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error",
            "error" : err.Error(),

		})
		return
	}

	c.JSON(200, result)
	
}
func (s *siswaAPI) GetList(c *gin.Context) {
    page, err := strconv.Atoi(c.Query("page"))
    if err != nil || page <= 0 {
        page = 1
    }

    pageSize, err := strconv.Atoi(c.Query("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 5000
    }

    result, totalPage, err := s.siswaService.GetList(page, pageSize)
    if err != nil {
		log.Printf("Pesan error: %v", err)

        c.JSON(500, gin.H{
            "message": "internal server error",
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


func (s *siswaAPI) History(c *gin.Context) {
    page, err := strconv.Atoi(c.Query("page"))
    if err != nil || page <= 0 {
        page = 1
    }

    pageSize, err := strconv.Atoi(c.Query("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 100
    }

	siswaNISN :=c.Param("nisn")
    kategori := c.Query("kategori")
    nama := c.Query("nama")
    tanggal := c.Query("tanggal")
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	result, totalPage, err := s.siswaService.HistoryPembayaranSiswa(siswaNISN,nama,tanggal,kategori, page, pageSize)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error",
            "error" : err.Error(),
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

func (s *siswaAPI) GetTotalGenderCount(c *gin.Context) {
    countLakiLaki, countPerempuan, err := s.siswaService.GetTotalGenderCount()
    if err != nil {
        log.Printf("Error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "count_laki_laki": countLakiLaki,
        "count_perempuan": countPerempuan,
    })
}

func (s *siswaAPI) Search(c *gin.Context) {
	name := c.Query("nama")
	kelas := c.Query("kelas")
	nisn := c.Query("nisn")
    jurusan := c.Query("jurusan")
    angkatan := c.Query("angkatan")


	siswaList, err := s.siswaService.Search(name,  nisn, kelas, jurusan, angkatan)
	if err != nil {
        log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
   

	c.JSON(http.StatusOK, gin.H{
        "data": siswaList,
        "message" : "Berhasil mendapatkan data"})
    }

func (s *siswaAPI)SearchByKode(c *gin.Context){
    name := c.Query("nama")
	kode := c.Query("kode")
	nisn := c.Query("nisn")


	siswaList, err := s.siswaService.SearchByKodeKelas(name,  nisn, kode)
	if err != nil {
        log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
   
    fmt.Println(siswaList)

	c.JSON(http.StatusOK, gin.H{
        "data": siswaList,
        "message" : "Berhasil mendapatkan data"})
    }


func (s *siswaAPI) ExportSiswa(c *gin.Context) {
    page, err := strconv.Atoi(c.Query("page"))
    if err != nil || page <= 0 {
        page = 1
    }

    pageSize, err := strconv.Atoi(c.Query("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 5000
    }

    result, _, err := s.siswaService.GetList(page, pageSize)
    if err != nil {
        log.Printf("Pesan error: %v", err)

        c.JSON(500, gin.H{
            "message": "internal server error",
            "error":   err.Error(),
        })
        return
    }

    file := excelize.NewFile()
    index, err := file.NewSheet("Siswa")
    if err != nil {
        log.Printf("Pesan error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    file.SetActiveSheet(index)

    // Add header row
    header := []string{"NO",  "NISN", "Nama", "Kelas", "Jurusan", "Agama", "Tempat Lahir", "Tanggal Lahir", "Gender", "Nama Ayah", "Nama Ibu", "Nomor Telepon", "Angkatan", "Email", "Alamat"}
    for col, val := range header {
        colName, err := excelize.ColumnNumberToName(col + 1)
        if err != nil {
        log.Printf("Pesan error: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        cell := colName + "1"
        file.SetCellValue("Siswa", cell, val)
    }

    // Add data rows
    for i, data := range result {
        row := i + 2
        file.SetCellValue("Siswa", "A"+strconv.Itoa(row), i+1)
        file.SetCellValue("Siswa", "B"+strconv.Itoa(row), data.NISN)
        file.SetCellValue("Siswa", "C"+strconv.Itoa(row), data.Nama)
        file.SetCellValue("Siswa", "D"+strconv.Itoa(row), data.Kelas)
        file.SetCellValue("Siswa", "E"+strconv.Itoa(row), data.Jurusan)
        file.SetCellValue("Siswa", "F"+strconv.Itoa(row), data.Agama)
        file.SetCellValue("Siswa", "G"+strconv.Itoa(row), data.TempatLahir)
        file.SetCellValue("Siswa", "H"+strconv.Itoa(row), data.TanggalLahir)
        file.SetCellValue("Siswa", "I"+strconv.Itoa(row), data.Gender)
        file.SetCellValue("Siswa", "J"+strconv.Itoa(row), data.NamaAyah)
        file.SetCellValue("Siswa", "K"+strconv.Itoa(row), data.NamaIbu)
        file.SetCellValue("Siswa", "L"+strconv.Itoa(row), data.NomorTelepon)
        file.SetCellValue("Siswa", "M"+strconv.Itoa(row), data.Angkatan)
        file.SetCellValue("Siswa", "N"+strconv.Itoa(row), data.Email)
        file.SetCellValue("Siswa", "O"+strconv.Itoa(row), data.Alamat)
    }

    fileName := "siswa.xlsx"
    err = file.SaveAs("./app/files/"+fileName)
    if err != nil {
        log.Printf("Pesan error: %v", err)

        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    filePath := "./app/files/" + fileName
    defer func() {
        if err := os.Remove(filePath); err != nil {
            log.Printf("Gagal menghapus file: %v", err)
        }
    }()
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Disposition", "attachment; filename="+fileName)
    c.Header("Content-Type", "application/octet-stream")
    c.Header("Content-Transfer-Encoding", "binary")
    c.File(filePath)
}

func (s *siswaAPI) DownloadSiswa(c *gin.Context) {
    name := c.Query("nama")
	kelas := c.Query("kelas")
	nisn := c.Query("nisn")
    jurusan := c.Query("jurusan")
    angkatan := c.Query("angkatan")


	result, err := s.siswaService.Search(name,  nisn, kelas, jurusan, angkatan)
	if err != nil {
        log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
    if len(result) == 0 {
        log.Printf("Error: %v", "data tidak ada")
        c.JSON(http.StatusNotFound, gin.H{"message": "Data tidak ditemukan"})
        return
    }
    

    file := excelize.NewFile()
    index, err := file.NewSheet("Siswa")
    if err != nil {
        log.Printf("Pesan error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    file.SetActiveSheet(index)

    // Add header row
    header := []string{"NO", "Nama", "NISN", "Kelas", "Jurusan", "Agama", "Tempat Lahir", "Tanggal Lahir", "Gender", "Nama Ayah", "Nama Ibu", "Nomor Telepon", "Angkatan", "Email", "Alamat"}
    for col, val := range header {
        colName, err := excelize.ColumnNumberToName(col + 1)
        if err != nil {
        log.Printf("Pesan error: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        cell := colName + "1"
        file.SetCellValue("Siswa", cell, val)
    }

    // Add data rows
    for i, data := range result {
        row := i + 2
        file.SetCellValue("Siswa", "A"+strconv.Itoa(row), i+1)
        file.SetCellValue("Siswa", "B"+strconv.Itoa(row), data.Nama)
        file.SetCellValue("Siswa", "C"+strconv.Itoa(row), data.NISN)
        file.SetCellValue("Siswa", "D"+strconv.Itoa(row), data.Kelas)
        file.SetCellValue("Siswa", "E"+strconv.Itoa(row), data.Jurusan)
        file.SetCellValue("Siswa", "F"+strconv.Itoa(row), data.Agama)
        file.SetCellValue("Siswa", "G"+strconv.Itoa(row), data.TempatLahir)
        file.SetCellValue("Siswa", "H"+strconv.Itoa(row), data.TanggalLahir)
        file.SetCellValue("Siswa", "I"+strconv.Itoa(row), data.Gender)
        file.SetCellValue("Siswa", "J"+strconv.Itoa(row), data.NamaAyah)
        file.SetCellValue("Siswa", "K"+strconv.Itoa(row), data.NamaIbu)
        file.SetCellValue("Siswa", "L"+strconv.Itoa(row), data.NomorTelepon)
        file.SetCellValue("Siswa", "M"+strconv.Itoa(row), data.Angkatan)
        file.SetCellValue("Siswa", "N"+strconv.Itoa(row), data.Email)
        file.SetCellValue("Siswa", "O"+strconv.Itoa(row), data.Alamat)
    }

    fileName := "siswa.xlsx"
    err = file.SaveAs("./app/files/"+fileName)
    if err != nil {
        log.Printf("Pesan error: %v", err)

        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    filePath := "./app/files/" + fileName
    defer func() {
        if err := os.Remove(filePath); err != nil {
            log.Printf("Gagal menghapus file: %v", err)
        }
    }()
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Disposition", "attachment; filename="+fileName)
    c.Header("Content-Type", "application/octet-stream")
    c.Header("Content-Transfer-Encoding", "binary")
    c.File(filePath)
    
}

func (s *siswaAPI) BiodataSiswa(c *gin.Context) {
	name := c.Query("nama")
	kelas := c.Query("kelas")
	nisn := c.Query("nisn")
	jurusan := c.Query("jurusan")
	angkatan := c.Query("angkatan")

	result, err := s.siswaService.Search(name, nisn, kelas, jurusan, angkatan)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(result) == 0 {
		log.Printf("Error: %v", "data tidak ada")
		c.JSON(http.StatusNotFound, gin.H{"message": "Data tidak ditemukan"})
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.SetHeaderFunc(func() {

		y := pdf.GetY()

		imageWidth := 20 
		imageHeight := 20 

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
	pdf.SetFont("Times", "BU", 14)
	pdf.CellFormat(0, 10, "Biodata Siswa", "0", 1, "C", false, 0, "")
	pdf.Ln(3)
	pdf.SetFont("Times", "", 12)
	pdf.SetX(10)
	pdf.CellFormat(0, 10, "Nama : "+result[0].Nama, "0", 1, "", false, 0, "")
	pdf.SetX(10)
	pdf.CellFormat(0, 10, "NISN : "+result[0].NISN, "0", 1, "", false, 0, "")
	pdf.SetX(10)
	pdf.CellFormat(0, 10, "Kelas : "+result[0].Kelas, "0", 1, "", false, 0, "")
	pdf.SetX(10)
	pdf.CellFormat(0, 10, "Jurusan : "+result[0].Jurusan, "0", 1, "", false, 0, "")
	pdf.SetX(10)
	pdf.CellFormat(0, 10, "Angkatan : "+result[0].Angkatan, "0", 1, "", false, 0, "")
	pdf.SetX(10)
	pdf.CellFormat(0, 10, "Alamat : "+result[0].Alamat, "0", 1, "", false, 0, "")
	pdf.SetX(10)
	pdf.CellFormat(0, 10, "Nomor Telepon : "+result[0].NomorTelepon, "0", 1, "", false, 0, "")
	pdf.SetX(10)
	pdf.CellFormat(0, 10, "Email : "+result[0].Email, "0", 1, "", false, 0, "")
	pdf.SetX(10)
	pdf.CellFormat(0, 10, "Tempat, Tanggal Lahir : "+result[0].TempatLahir+", "+result[0].TanggalLahir, "0", 1, "", false, 0, "")
	pdf.SetX(10)
	pdf.CellFormat(0, 10, "Nama Ayah : "+result[0].NamaAyah, "0", 1, "", false, 0, "")
	pdf.SetX(10)
	pdf.CellFormat(0, 10, "Nama Ibu : "+result[0].NamaIbu, "0", 1, "", false, 0, "")
	pdf.Ln(10)

    pdf.SetFont("Times", "BI", 12)
	// pdf.SetX(float64(10))
	pdf.CellFormat(0, 2, "*Staf Tata Usaha SMA Plus Nurul Iman Leles,*", "0", 1, "C", false, 0, "")
	// pdf.Ln(25)
	// pdf.SetFont("Times", "B", 11)
	// pdf.SetX(float64(150))
    // pdf.CellFormat(0, 3, "*Bendahara Sekolah*", "0", 1, "", false, 0, "")

	fileName := "biodata_siswa.pdf"
	err = pdf.OutputFileAndClose("./app/files/" + fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer func() {
		if err := os.Remove("./app/files/" + fileName); err != nil {
			log.Printf("Gagal menghapus file: %v", err)
		}
	}()

	c.File("./app/files/" + fileName)
}

func (s *siswaAPI) ImportFromExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Pesan error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tidak ada file yang diterima"})
		return
	}

	filePath := "./app/files/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("Pesan error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
		return
	}

	if err := s.siswaService.ImportFromExcel(filePath); err != nil {
		log.Printf("Pesan error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
    defer func() {
		if err := os.Remove("./app/files/" + file.Filename); err != nil {
			log.Printf("Gagal menghapus file: %v", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil menambahkan data"})
}