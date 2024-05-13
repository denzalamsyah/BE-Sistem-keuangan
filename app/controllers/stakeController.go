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

type GuruAPI interface {
	AddStake(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
	GetList(c *gin.Context)
	GetTotalGenderCount(c *gin.Context)
	Search(c *gin.Context)
    GetHistoriKas(c *gin.Context)
    GetTotalKasByNIP(c *gin.Context)
	AmbilKasGuru(c *gin.Context)
    GetHistoriPengambilanKas(c *gin.Context)
    ExportGuru(c *gin.Context)

}

type guruAPI struct{
	guruService services.GuruServices
}

func NewGuruAPI(guruRepo services.GuruServices) *guruAPI{
	return &guruAPI{guruRepo}
}

func (s *guruAPI) AddStake(c *gin.Context){
	var stake models.Guru

	if err := c.ShouldBind(&stake); err != nil{
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	file, err := c.FormFile("file")
    if err != nil {
        // Jika tidak ada file yang diunggah, gunakan gambar default
        stake.Gambar = "https://res.cloudinary.com/dgvkpzi4p/image/upload/v1706338009/149071_fxemnm.png"
    } else {
        // Jika ada file yang diunggah, upload ke Cloudinary dan dapatkan URL-nya
        imageURL, err := middleware.UploadToCloudinary(file)
        if err != nil {
		log.Printf("Pesan error: %v", err)

            c.JSON(http.StatusInternalServerError, gin.H{
                "message": "failed to upload image to Cloudinary",
                "error":   err.Error(),
            })
            return
        }
        stake.Gambar = imageURL
    }

	err = s.guruService.Store(&stake)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		existingStake, err := s.guruService.GetUserNIP(stake.Nip)
        if err != nil {
            log.Printf("Error checking NIP: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   err.Error(),
            })
            return
        }
        
        if existingStake.Nip == stake.Nip{
            c.JSON(http.StatusBadRequest, gin.H{
                "message":   "Guru dengan NIP tersebut sudah ada",
                "error":   "Gagal menambah data",
            })
            return
        }
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil menambah data",
		"data" : stake,
	})	

}
func (s *guruAPI) Update(c *gin.Context) {
    guruNIP := c.Param("nip")
    if guruNIP == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "invalid request body" ,
        })
        return
    }

   

    var newGuru models.Guru
    if err := c.ShouldBind(&newGuru); err != nil {
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
                "message": "failed to upload image to Cloudinary",
                "error":   err.Error(),
            })
            return
        }
        newGuru.Gambar = imageURL
    }


    err = s.guruService.Update(guruNIP, newGuru)
    if err != nil {
		log.Printf("Update error: %v", err)
        existingSiswa, err := s.guruService.GetUserNIP(newGuru.Nip)
        if err != nil {
            log.Printf("Error checking NISN: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   err.Error(),
            })
            return
        }
        
        if existingSiswa.Nip== newGuru.Nip{
            c.JSON(http.StatusBadRequest, gin.H{
                "message":   "NIP yang anda masukan sudah ada",
                "error":   "Gagal mengubah data",
            })
            return
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Berhasil mengubah siswa",
        "data":    newGuru,
    })
}

func (s *guruAPI) Delete(c *gin.Context){
	stakeNIP := c.Param("nip")
	
	if stakeNIP == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	

	err := s.guruService.Delete(stakeNIP)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message" : "Gagal menghapus data",
			"error" : err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "Berhasil menghapus",
	})
}

func (s *guruAPI) GetByID(c *gin.Context){
	stakeNIP := c.Param("nip")
	
	result, err := s.guruService.GetByID(stakeNIP)	
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

func (s *guruAPI) GetList(c *gin.Context){
	page, err := strconv.Atoi(c.Query("page"))
    if err != nil || page <= 0 {
        page = 1
    }

    pageSize, err := strconv.Atoi(c.Query("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 100
    }

	result, totalPage, err := s.guruService.GetList(page, pageSize)
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

func (s *guruAPI) GetTotalGenderCount(c *gin.Context) {
    countLakiLaki, countPerempuan, err := s.guruService.GetTotalGenderCount()
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

func (s *guruAPI) Search(c *gin.Context){
	nama := c.Query("nama")
	nip := c.Query("nip")
	jabatan := c.Query("jabatan")

	stakeList, err := s.guruService.Search(nama,  nip, jabatan)
	if err != nil {
        log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stakeList})
}

func ( s *guruAPI) GetHistoriKas(c *gin.Context){
    page, err := strconv.Atoi(c.Query("page"))
    if err != nil || page <= 0 {
        page = 1
    }

    pageSize, err := strconv.Atoi(c.Query("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 100
    }

	guruNIP := c.Param("nip")
	
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	result, totalPage, err := s.guruService.HistoryPembayaranGuru(guruNIP, page, pageSize)
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

// GetTotalKasGuru mengembalikan total kas guru dalam bentuk JSON
func (s *guruAPI) GetTotalKasByNIP(c *gin.Context) {
	guruNIP :=c.Param("nip")

	totalKas, err := s.guruService.SaldoKasByNIP(guruNIP)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message": "NIP tidak ditemukan",
			"error":   err.Error(),
		})
		return
	}

	guru, err := s.guruService.GetUserNIP(guruNIP)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"message": "internal server error",
			"error":   err.Error(),
		})
		return
	}

	

	response := models.TotalSaldoKas{
		NIP: guru.Nip,
		Nama: guru.Nama,
		TotalKas: totalKas,
	}

	c.JSON(200, gin.H{
		"data" : response,
	})
}

// AmbilKasGuru mengurangi saldo uang kas guru dalam bentuk JSON
func (s *guruAPI) AmbilKasGuru(c *gin.Context) {

	var requestBody models.PengambilanKas

	if err := c.BindJSON(&requestBody); err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message": "invalid request body",
		})
		return
	}

	guruNIP := c.Param("nip")
	 existingSiswa, err := s.guruService.GetUserNIP(guruNIP)
        if err != nil {
            log.Printf("Error checking NIP: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   err.Error(),
            })
            return
        }

	err = s.guruService.AmbilKasGuru(requestBody.JumlahAmbil, guruNIP, existingSiswa.Nama, requestBody.TanggalAmbil)
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(500, gin.H{
			"error": "internal server error",
			"message":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Berhasil mengambil kas guru",
	})
}

// histori pengambilan
func ( s *guruAPI) GetHistoriPengambilanKas(c *gin.Context){
    page, err := strconv.Atoi(c.Query("page"))
    if err != nil || page <= 0 {
        page = 1
    }

    pageSize, err := strconv.Atoi(c.Query("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 100
    }

	guruNIP :=c.Param("nip")
	
	if err != nil {
		log.Printf("Pesan error: %v", err)

		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	result, totalPage, err := s.guruService.HistoryPengambilanKas(guruNIP, page, pageSize)
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

// download data guru
func (s *guruAPI) ExportGuru(c *gin.Context) {
    page, err := strconv.Atoi(c.Query("page"))
    if err != nil || page <= 0 {
        page = 1
    }

    pageSize, err := strconv.Atoi(c.Query("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 5000
    }

    result, _, err := s.guruService.GetList(page, pageSize)
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
    header := []string{"NO","NIP", "Nama", "Agama", "Jabatan", "Tempat Lahir", "Tanggal Lahir", "Gender", "Nomor Telepon", "Email", "Alamat"}
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
        file.SetCellValue("Siswa", "B"+strconv.Itoa(row), data.Nip)
        file.SetCellValue("Siswa", "C"+strconv.Itoa(row), data.Nama)
        file.SetCellValue("Siswa", "D"+strconv.Itoa(row), data.Agama)
        file.SetCellValue("Siswa", "E"+strconv.Itoa(row), data.Jabatan)
        file.SetCellValue("Siswa", "F"+strconv.Itoa(row), data.TempatLahir)
        file.SetCellValue("Siswa", "G"+strconv.Itoa(row), data.TanggalLahir)
        file.SetCellValue("Siswa", "H"+strconv.Itoa(row), data.Gender)
        file.SetCellValue("Siswa", "I"+strconv.Itoa(row), data.NomorTelepon)
        file.SetCellValue("Siswa", "J"+strconv.Itoa(row), data.Email)
        file.SetCellValue("Siswa", "K"+strconv.Itoa(row), data.Alamat)
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