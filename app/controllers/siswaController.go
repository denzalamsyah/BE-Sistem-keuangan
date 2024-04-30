package controllers

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/denzalamsyah/simak/app/middleware"
	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/services"
	"github.com/gin-gonic/gin"
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
    ExportSiswa(c *gin.Context)
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
            "message": "invalid request body" + err.Error(),
            "error":   "Gagal mengubah data",
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
        if newSiswa.Alamat == "" || newSiswa.AgamaID == 0 || newSiswa.Angkatan == "" || newSiswa.Email == "" || newSiswa.GenderID == 0 || newSiswa.JurusanID == 0 || newSiswa.KelasID == 0 || newSiswa.NISN == 0 || newSiswa.Nama == "" || newSiswa.NamaAyah == "" || newSiswa.NamaIbu == "" || newSiswa.NomorTelepon == 0 || newSiswa.TanggalLahir== "" || newSiswa.TempatLahir == ""{ 
            c.JSON(http.StatusBadRequest, gin.H{"error":   "semua item harus di isi kecuali gambar",})
			return
		}
		log.Printf("Error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "message": err.Error(),
            "error":   "Gagal menambah data",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Berhasil menambahkan data",
        "data":    newSiswa,
    })
}

func (s *siswaAPI) Update(c *gin.Context) {
    siswaID := c.Param("id")

    if siswaID == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "invalid request body" ,
        })
        return
    }

    id, err := strconv.Atoi(siswaID)
    if err != nil {
		log.Printf("Pesan error: %v", err)

        c.JSON(http.StatusBadRequest, gin.H{
            "message": "invalid request body" + err.Error(),
        })
        return
    }

    var existingSiswa models.Siswa
    if err := c.ShouldBind(&existingSiswa); err != nil {
		log.Printf("Pesan error: %v", err)

        c.JSON(http.StatusBadRequest, gin.H{
            "message": "invalid request body" + err.Error(),
            "error":  "Gagal mengubah data",
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
		log.Printf("Pesan error: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "failed to upload image to Cloudinary" + err.Error(),
                "error":   "Gagal mengubah data",
            })
            return
        }
        existingSiswa.Gambar = imageURL
    }

    // Lakukan pembaruan data siswa di database
    err = s.siswaService.Update(id, existingSiswa)
    if err != nil {
		log.Printf("Pesan error: %v", err)

        c.JSON(http.StatusInternalServerError, gin.H{
            "message": err.Error(),
            "error":   "Gagal mengubah data",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Berhasil mengubah siswa",
        "data":    existingSiswa,
    })
}

func (s *siswaAPI) Delete(c *gin.Context) {

	siswaID := c.Param("id")
	
	if siswaID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(siswaID)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body" + err.Error(),
            
		})
		return
	}

	err = s.siswaService.Delete(id)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : err.Error(),
            "error" : "Gagal menghapus data",
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil menghapus data",
	})
	
}
func (s *siswaAPI) GetByID(c *gin.Context) {

	siswaID, err := strconv.Atoi(c.Param("id"))
	
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body" + err.Error(),
		})
		return
	}

	result, err := s.siswaService.GetByID(siswaID)	
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "internal server error" + err.Error(),
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

	siswaID, err := strconv.Atoi(c.Param("id"))
	
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	result, totalPage, err := s.siswaService.HistoryPembayaranSiswa(siswaID, page, pageSize)
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

	siswaList, err := s.siswaService.Search(name,  nisn, kelas, jurusan)
	if err != nil {
        log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": siswaList})
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
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    file.SetActiveSheet(index)

    // Add header row
    header := []string{"NO", "Nama", "NISN", "Kelas", "Jurusan", "Agama", "Tempat Lahir", "Tanggal Lahir", "Gender", "Nama Ayah", "Nama Ibu", "Nomor Telepon", "Angkatan", "Email", "Alamat"}
    for col, val := range header {
        colName, err := excelize.ColumnNumberToName(col + 1)
        if err != nil {
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
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    defer os.Remove(fileName)

    // Return the file as attachment
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Disposition", "attachment; filename="+fileName)
    c.Header("Content-Type", "application/octet-stream")
    c.Header("Content-Transfer-Encoding", "binary")
    c.File(fileName)
}


