package services

import (
	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/repository"
)

type GuruServices interface {
	Store(Guru *models.Guru) error
	Update(nip int, Guru models.Guru) error
	Delete(nip int) error
	GetByID(nip int) (*models.GuruResponse, error)
	GetList(page, pageSize int) ([]models.GuruResponse, int, error)
	GetTotalGenderCount() (int, int, error)
	Search(nama, nip, jabatan string) ([]models.GuruResponse, error)
	GetUserNIP(nip int) (models.Guru, error)
	HistoryPembayaranGuru(nip, page, pageSize int) ([]models.HistoryPembayaranKas, int, error)
	SaldoKasByNIP(nip int) (int, error)
	AmbilKasGuru(nip, jumlah int,nama, tanggal string) error
	HistoryPengambilanKas(nip, page, pageSize int) ([]models.HistoryPengambilanKas, int, error)
}

type guruServices struct{
	guruRepo repository.GuruRepository
}
func NewGuruService(guruRepo repository.GuruRepository) GuruServices{
	return &guruServices{guruRepo}
}
func (c *guruServices) Store(Guru *models.Guru) error{
	err := c.guruRepo.Store(Guru)
	if err != nil {
		return err
	}
	return nil
}

func (c *guruServices) Update(nip int, Guru models.Guru) error {
	err := c.guruRepo.Update(nip, Guru)
	if err != nil {
		
		return err
	}
	return nil
}
func (c *guruServices) Delete(nip int) error{
	err := c.guruRepo.Delete(nip)
	if err != nil {
		return err
	}
	return nil
}

func (c *guruServices) GetByID(nip int) (*models.GuruResponse, error){
	guru, err := c.guruRepo.GetByID(nip)
	if err != nil {
		return nil, err
	}
	return guru, nil
}

func (c *guruServices) GetList(page, pageSize int)([]models.GuruResponse, int, error){
	guru, totalPage, err := c.guruRepo.GetList(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return guru, totalPage, nil
}

func (c *guruServices) GetTotalGenderCount() (int, int, error) {
	countLakiLaki, countPerempuan, err := c.guruRepo.GetTotalGenderCount()
	if err != nil {
		return 0, 0, err
	}
	return int(countLakiLaki), int(countPerempuan), nil
}

func(c *guruServices) Search(nama, nip, jabatan string) ([]models.GuruResponse, error){
	guru, err := c.guruRepo.Search(nama, nip, jabatan)

	if err != nil {
        return nil, err
    }
	return guru, nil
}
func (c *guruServices) GetUserNIP(guru int) (models.Guru, error) {
	return c.guruRepo.GetUserNIP(guru)
}

func (c *guruServices) HistoryPembayaranGuru(nip, page, pageSize int) ([]models.HistoryPembayaranKas, int, error){
	history, totalPage, err := c.guruRepo.HistoryPembayaranGuru(nip, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return history, totalPage, nil
}

// HitungTotalKasGuru menghitung total kas guru berdasarkan NIP
func (c *guruServices) SaldoKasByNIP(nip int) (int, error) {
	totalKas, err := c.guruRepo.SaldoKasByNIP(nip)
	if err != nil {
		return 0, err
	}
	return totalKas, nil
}

// AmbilKasGuru mengurangi saldo uang kas guru berdasarkan NIP
func (c *guruServices) AmbilKasGuru(nip, jumlah int, nama, tanggal string) error {
	err := c.guruRepo.AmbilKasGuru(nip, jumlah, nama, tanggal)
	if err != nil {
		return err
	}
	return nil
}


func (c *guruServices) HistoryPengambilanKas(nip, page, pageSize int) ([]models.HistoryPengambilanKas, int, error){
	history, totalPage, err := c.guruRepo.HistoryPengambilanKas(nip, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return history, totalPage, nil
}