package repository

import (
	"errors"
	"math"
	"strings"
	"time"

	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type SemesterRepository interface {
	Store(PembayaranSemester *models.PembayaranSemester) error
	Update(id int, PembayaranSemester models.PembayaranSemester) error
	Delete(id int) error
	GetByID(id int) (*models.PembayaranSemesterResponse, error)
	GetList(page, pageSize int) ([]models.PembayaranSemesterResponse, int, error)
    Search(siswa, tahunAjar, transaksi, semester, tanggal string) ([]models.PembayaranSemesterResponse, error)
    GetLunasByNISN(nisn int) ([]models.PembayaranSemesterResponse, error)
}

type semesterRepository struct{
	db *gorm.DB
}

func NewSemesterRepo(db *gorm.DB) *semesterRepository{
	return &semesterRepository{db}
}

func (c *semesterRepository) Store(PembayaranSemester *models.PembayaranSemester) error {
	tx := c.db.Begin()

    var count int64
    if err := tx.Model(&models.PembayaranSemester{}).
        Where("transaksi_id = ?", PembayaranSemester.TransaksiID).
        Where("siswa_id = ?", PembayaranSemester.SiswaID).
        Where("kelas = ?", PembayaranSemester.Kelas).
        Where("semester = ?", PembayaranSemester.Semester).
        Where("tanggal = ?", PembayaranSemester.Tanggal).
        Count(&count).Error; err != nil {
        tx.Rollback()
        return err
    }

    if count > 0 {
        tx.Rollback()
        return errors.New("Pembayaran sudah ada")
    }

    PembayaranSemester.CreatedAt = time.Now()
    
    if err := tx.Create(PembayaranSemester).Error; err != nil {
        tx.Rollback()
        return err
    }

    var transaksi models.Transaksi
    if err := tx.First(&transaksi, PembayaranSemester.TransaksiID).Error; err != nil {
        tx.Rollback()
        return err
    }

    if PembayaranSemester.Jumlah > transaksi.Jumlah {
        tx.Rollback()
        return errors.New("Nilai pembayaran terlalu besar")
    }

    pemasukanSemester := models.Pemasukan{
        IDBayar: PembayaranSemester.ID,
        Nama:    transaksi.Nama,
        Tanggal: PembayaranSemester.Tanggal,
        Jumlah:  PembayaranSemester.Jumlah,
    }
    if err := tx.Create(&pemasukanSemester).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Cek apakah Jumlah yang dibayar sama dengan JumlahBayar pada transaksi
    status := "BELUM LUNAS"
    if PembayaranSemester.Jumlah == transaksi.Jumlah {
        status = "LUNAS"
    }

    // Update Status pada PembayaranSemester
    if err := tx.Model(&PembayaranSemester).Update("status", status).Error; err != nil {
        tx.Rollback()
        return err
    }
    return tx.Commit().Error
}

func (c *semesterRepository) Update(id int, PembayaranSemester models.PembayaranSemester) error {
    tx := c.db.Begin()
    var pembayaranSemester models.PembayaranSemester

    // Mengambil data PembayaranSemester yang akan diupdate
    if err := tx.Where("id = ?", id).First(&pembayaranSemester).Error; err != nil {
        tx.Rollback()
        return err
    }

     var transaksi models.Transaksi
     if err := tx.First(&transaksi, pembayaranSemester.TransaksiID).Error; err != nil {
         tx.Rollback()
         return err
     }
     if PembayaranSemester.Jumlah > transaksi.Jumlah {
        tx.Rollback()
        return errors.New("Jumlah terlalu besar")
    }

    // Memperbarui data PembayaranSemester dengan nilai baru
    if err := tx.Model(&pembayaranSemester).Updates(&PembayaranSemester).Error; err != nil {
        tx.Rollback()
        return err
    }
    

    // Memperbarui data Pemasukan yang terkait dengan PembayaranSemester
    if err := tx.Model(&models.Pemasukan{}).
        Where("id_bayar = ?", pembayaranSemester.ID).
        Updates(map[string]interface{}{
            "nama":    transaksi.Nama,
            "tanggal": pembayaranSemester.Tanggal,
            "jumlah":  pembayaranSemester.Jumlah,
            "updated_at": time.Now(),
        }).Error; err != nil {
        tx.Rollback()
        return err
    }
    status := "BELUM LUNAS"
    if pembayaranSemester.Jumlah == transaksi.Jumlah {
        status = "LUNAS"
    }
    if err := tx.Model(&pembayaranSemester).Update("status", status).Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}


func (c *semesterRepository) Delete(id int) error {
    tx := c.db.Begin()

    var pembayaranSemester models.PembayaranSemester

    if err := tx.Where("id = ?", id).First(&pembayaranSemester).Error; err != nil {
        tx.Rollback()
        return err
    }

    if err := tx.Delete(&pembayaranSemester).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Hapus data pembayaran yang terkait dengan pembayaran semester
    if err := tx.Where("id_bayar = ?", pembayaranSemester.ID).Delete(&models.Pemasukan{}).Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}


func (c *semesterRepository) GetByID(id int) (*models.PembayaranSemesterResponse, error) {
	var PembayaranSemester models.PembayaranSemester
	err := c.db.Preload("Siswa").Preload("Transaksi").Where("id = ?", id).First(&PembayaranSemester).Error
	if err != nil {
		return nil, err
	}

	PembayaranSemesterResponse := models.PembayaranSemesterResponse{
		ID:             PembayaranSemester.ID,
		Siswa:          PembayaranSemester.Siswa.Nama,
        NISN: PembayaranSemester.Siswa.Nisn,
		Transaksi: PembayaranSemester.Transaksi.Nama,
		Semester:       PembayaranSemester.Semester,
		Kelas:      PembayaranSemester.Kelas,
		Status:         PembayaranSemester.Status,
		Tanggal:        PembayaranSemester.Tanggal,
		Jumlah:         PembayaranSemester.Jumlah,
        CreatedAt: PembayaranSemester.CreatedAt,
        UpdatedAt: PembayaranSemester.UpdatedAt,
	}

	return &PembayaranSemesterResponse, nil
}

func (c *semesterRepository) GetList(page, pageSize int) ([]models.PembayaranSemesterResponse, int, error) {
	var PembayaranSemester []models.PembayaranSemester
	err := c.db.Preload("Siswa").Preload("Transaksi").
    Offset((page - 1) * pageSize).Limit(pageSize).Find(&PembayaranSemester).Error
	if err != nil {
		return nil, 0, err
	}

    var totalData int64
    if err := c.db.Model(&models.PembayaranSemester{}).Count(&totalData).Error; err != nil {
        return nil, 0, err
    }

	var PembayaranSemesterResponse []models.PembayaranSemesterResponse
	for _, s := range PembayaranSemester{
		PembayaranSemesterResponse = append(PembayaranSemesterResponse, models.PembayaranSemesterResponse{
			ID:             s.ID,
			Siswa:          s.Siswa.Nama,
            NISN: s.Siswa.Nisn,
			Transaksi: s.Transaksi.Nama,
			Semester:       s.Semester,
			Kelas:      s.Kelas,
			Status:         s.Status,
			Tanggal:        s.Tanggal,
			Jumlah:         s.Jumlah,
            CreatedAt: s.CreatedAt,
            UpdatedAt: s.UpdatedAt,
		})
	}
	totalPage := int(math.Ceil(float64(totalData) / float64(pageSize)))

	return PembayaranSemesterResponse, totalPage, nil
}

func (c *semesterRepository) Search(siswa, tahunAjar, transaksi, semester, tanggal string) ([]models.PembayaranSemesterResponse, error){
    siswa = strings.ToLower(siswa)
    tahunAjar = strings.ToLower(tahunAjar)
    transaksi = strings.ToLower(transaksi)
    semester = strings.ToLower(semester)
    tanggal = strings.ToLower(tanggal)
    // penerima = strings.ToLower(penerima)

    var pembayaran []models.PembayaranSemesterResponse

    query := c.db.Table("pembayaran_semesters").
    Select("pembayaran_semesters.id, siswas.nama as siswa, transaksis.nama as transaksi, pembayaran_semesters.semester, pembayaran_semesters.tahun_ajar, pembayaran_semesters.jumlah, pembayaran_semesters.tanggal, pembayaran_semesters.status, pembayaran_semesters.created_at, pembayaran_semesters.updated_at").
    Joins("JOIN siswas ON pembayaran_semesters.siswa_id = siswas.nisn").
    Joins("JOIN transaksis ON pembayaran_semesters.transaksi_id = transaksis.id").
    Where("LOWER(siswas.nama) LIKE ? AND LOWER(pembayaran_semesters.tahun_ajar) LIKE ? AND LOWER(transaksis.nama) LIKE ? AND LOWER(pembayaran_semesters.semester) LIKE ? AND LOWER(pembayaran_semesters.tanggal) LIKE ?", "%"+siswa+"%", "%"+tahunAjar+"%", "%"+transaksi+"%", "%"+semester+"%", "%"+tanggal+"%")

if err := query.Find(&pembayaran).Error; err != nil {
    return nil, err
}

return pembayaran, nil
}

func (c *semesterRepository) GetLunasByNISN(nisn int) ([]models.PembayaranSemesterResponse, error) {
	var pembayaran []models.PembayaranSemesterResponse

	if err := c.db.Model(&models.PembayaranSemester{}).
		Select("pembayaran_semesters.id, siswas.nisn, siswas.nama as siswa, transaksis.nama as transaksi, pembayaran_semesters.semester, pembayaran_semesters.tahun_ajar, pembayaran_semesters.tanggal, pembayaran_semesters.jumlah, pembayaran_semesters.status, pembayaran_semesters.created_at, pembayaran_semesters.updated_at").
		Joins("JOIN transaksis ON pembayaran_semesters.transaksi_id = transaksis.id").
		Joins("JOIN siswas ON pembayaran_semesters.siswa_id = siswas.nisn").
		Where("transaksis.jumlah_bayar = pembayaran_semesters.jumlah").
		Where("pembayaran_semesters.siswa_id = ?", nisn).
		Scan(&pembayaran).Error; err != nil {
		return nil, err
	}

	return pembayaran, nil
}