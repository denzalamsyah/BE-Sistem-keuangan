package services

import (
	"log"

	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/repository"
)

type SiswaServices interface {
	Store(Siswa *models.Siswa) error
	Update(nisn string, Siswa models.Siswa) error
	Delete(nisn string) error
	GetByID(nisn string) (*models.SiswaResponse, error)
	GetList(page, pageSize int) ([]models.SiswaResponse, int, error)
	HistoryPembayaranSiswa(siswaID,nama, tanggal, kategori string, page, pageSize int) ([]models.HistoryPembayaran, int, error)
	GetTotalGenderCount() (int, int, error)
	Search(name, kelas, nisn, jurusan, angkatan string) ([]models.SiswaResponse, error)
	SearchByKodeKelas(name, nisn, kodeKelas string) ([]models.SiswaResponse, error)
	GetUserNisn(nisn string) (models.Siswa, error)
	ImportFromExcel(filePath string) error

}

type siswaServices struct {
	siswaRepo repository.SiswaRepository
}
func NewSiswaService(siswaRepo repository.SiswaRepository) SiswaServices {
	return &siswaServices{siswaRepo}
}


func (c *siswaServices) Store(Siswa *models.Siswa) error {
    err := c.siswaRepo.Store(Siswa)
    if err != nil {
        return err
    }
    return nil
}

func (c *siswaServices) Update(nisn string, siswa models.Siswa) error {
	err := c.siswaRepo.Update(nisn, siswa)
	if err != nil {
		
		return err
	}
	return nil
}

func (c *siswaServices) Delete(nisn string) error {
	err := c.siswaRepo.Delete(nisn)
	if err != nil {
		return err
	}
	return nil
}

func (c *siswaServices) GetByID(nisn string) (*models.SiswaResponse, error) {
	siswa, err := c.siswaRepo.GetByID(nisn)
	if err != nil {
		return nil, err
	}
	return siswa, nil
}

func (c *siswaServices) GetList(page, pageSize int) ([]models.SiswaResponse, int, error) {
    siswas, totalPage, err := c.siswaRepo.GetList(page, pageSize)
    if err != nil {
        return nil, 0, err
    }
    return siswas, totalPage, nil
}


func (c *siswaServices) HistoryPembayaranSiswa(siswaID,nama, tanggal, kategori string, page, pageSize int) ([]models.HistoryPembayaran, int, error) {
	history, totalPage, err := c.siswaRepo.HistoryPembayaranSiswa(siswaID,nama, tanggal, kategori, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return history, totalPage, nil
}

func (c *siswaServices) GetTotalGenderCount() (int, int, error) {
	countLakiLaki, countPerempuan, err := c.siswaRepo.GetTotalGenderCount()
	if err != nil {
		return 0, 0, err
	}
	return int(countLakiLaki), int(countPerempuan), nil
}

func(c *siswaServices) Search(name, nisn, kelas, jurusan, angkatan string) ([]models.SiswaResponse, error){
	siswa, err := c.siswaRepo.Search(name, nisn, kelas, jurusan, angkatan)
	if err != nil {
        return nil, err
    }
	return siswa, nil
}

func(c *siswaServices) SearchByKodeKelas(name, nisn, kodeKelas string) ([]models.SiswaResponse, error){
	siswa, err := c.siswaRepo.SearchByKodeKelas(name, nisn, kodeKelas)
	if err != nil {
        return nil, err
    }
	return siswa, nil
}

func (c *siswaServices) GetUserNisn(nisn string) (models.Siswa, error) {
	return c.siswaRepo.GetUserNisn(nisn)
}

func (c *siswaServices) ImportFromExcel(filePath string) error {
    err := c.siswaRepo.ImportFromExcel(filePath)
    if err != nil {
        log.Printf("Error importing from Excel: %v", err)
        return err
    }
    return nil
}