package repository

import (
	"math"
	"strings"

	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type SemesterRepository interface {
	Store(PembayaranSemester *models.PembayaranSemester) error
	Update(id int, PembayaranSemester models.PembayaranSemester) error
	Delete(id int) error
	GetByID(id int) (*models.PembayaranSemesterResponse, error)
	GetList(page, pageSize int) ([]models.PembayaranSemesterResponse, int, error)
    Search(siswa, tahunAjar, transaksi, semester, tanggal, penerima string) ([]models.PembayaranSemesterResponse, error)
}

type semesterRepository struct{
	db *gorm.DB
}

func NewSemesterRepo(db *gorm.DB) *semesterRepository{
	return &semesterRepository{db}
}

func (c *semesterRepository) Store(PembayaranSemester *models.PembayaranSemester) error {
	tx := c.db.Begin()

    if err := tx.Create(PembayaranSemester).Error; err != nil {
        tx.Rollback()
        return err
    }
	 
	var transaksi models.Transaksi
    if err := tx.First(&transaksi, PembayaranSemester.TransaksiID).Error; err != nil {
        tx.Rollback()
        return err
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

    // Memperbarui data PembayaranSemester dengan nilai baru
    if err := tx.Model(&pembayaranSemester).Updates(&PembayaranSemester).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Mengambil data Transaksi untuk mendapatkan nilai nama transaksi
    var transaksi models.Transaksi
    if err := tx.First(&transaksi, pembayaranSemester.TransaksiID).Error; err != nil {
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
        }).Error; err != nil {
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
	err := c.db.Preload("Siswa").Preload("Penerima").Preload("Transaksi").Where("id = ?", id).First(&PembayaranSemester).Error
	if err != nil {
		return nil, err
	}

	PembayaranSemesterResponse := models.PembayaranSemesterResponse{
		ID:             PembayaranSemester.ID,
		Siswa:          PembayaranSemester.Siswa.Nama,
		Transaksi: PembayaranSemester.Transaksi.Nama,
		Semester:       PembayaranSemester.Semester,
		TahunAjar:      PembayaranSemester.TahunAjar,
		Status:         PembayaranSemester.Status,
		Tanggal:        PembayaranSemester.Tanggal,
		Jumlah:         PembayaranSemester.Jumlah,
		Penerima:       PembayaranSemester.Penerima.Nama,
	}

	return &PembayaranSemesterResponse, nil
}

func (c *semesterRepository) GetList(page, pageSize int) ([]models.PembayaranSemesterResponse, int, error) {
	var PembayaranSemester []models.PembayaranSemester
	err := c.db.Preload("Siswa").Preload("Penerima").Preload("Transaksi").
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
			Transaksi: s.Transaksi.Nama,
			Semester:       s.Semester,
			TahunAjar:      s.TahunAjar,
			Status:         s.Status,
			Tanggal:        s.Tanggal,
			Jumlah:         s.Jumlah,
			Penerima:       s.Penerima.Nama,
		})
	}
	totalPage := int(math.Ceil(float64(totalData) / float64(pageSize)))

	return PembayaranSemesterResponse, totalPage, nil
}

func (c *semesterRepository) Search(siswa, tahunAjar, transaksi, semester, tanggal, penerima string) ([]models.PembayaranSemesterResponse, error){
    siswa = strings.ToLower(siswa)
    tahunAjar = strings.ToLower(tahunAjar)
    transaksi = strings.ToLower(transaksi)
    semester = strings.ToLower(semester)
    tanggal = strings.ToLower(tanggal)
    penerima = strings.ToLower(penerima)

    var pembayaran []models.PembayaranSemesterResponse

    query := c.db.Table("pembayaran_semesters").
    Select("pembayaran_semesters.id, siswas.nama as siswa, transaksis.nama as transaksi, pembayaran_semesters.semester, pembayaran_semesters.tahun_ajar, stakeholders.nama as penerima, pembayaran_semesters.jumlah, pembayaran_semesters.tanggal, pembayaran_semesters.status").
    Joins("JOIN siswas ON pembayaran_semesters.siswa_id = siswas.id").
    Joins("JOIN transaksis ON pembayaran_semesters.transaksi_id = transaksis.id").
    Joins("JOIN stakeholders ON pembayaran_semesters.penerima_id = stakeholders.id").
    Where("LOWER(siswas.nama) LIKE ? AND LOWER(pembayaran_semesters.tahun_ajar) LIKE ? AND LOWER(transaksis.nama) LIKE ? AND LOWER(pembayaran_semesters.semester) LIKE ? AND LOWER(pembayaran_semesters.tanggal) LIKE ? AND LOWER(stakeholders.nama) LIKE ?", "%"+siswa+"%", "%"+tahunAjar+"%", "%"+transaksi+"%", "%"+semester+"%", "%"+tanggal+"%", "%"+penerima+"%")

if err := query.Find(&pembayaran).Error; err != nil {
    return nil, err
}

return pembayaran, nil
}