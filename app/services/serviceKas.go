package services

import (
	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/repository"
)

type KasServices interface {
	Store(KasGuru *models.KasGuru) error
	Update(id int, KasGuru models.KasGuru) error
	GetList(page, pageSize int) ([]models.KasGuruResponse,int, error)
	Delete(id int) error
	Search(nama, tanggal string) ([]models.KasGuruResponse, error)
	GetByID(id int) (*models.HistoryPembayaranKas, error)
	GetAmbilByID(id int) (*models.HistoryPengambilanKas, error)
}

type kasServices struct {
	kasRepo repository.KasRepository
}
func NewKasService(kasRepo repository.KasRepository) KasServices {
	return &kasServices{kasRepo}
}


func (c *kasServices) Store(KasGuru *models.KasGuru) error {
    err := c.kasRepo.Store(KasGuru)
    if err != nil {
        return err
    }
    return nil
}

func (c *kasServices) Update(id int, KasGuru models.KasGuru) error {
	err := c.kasRepo.Update(id, KasGuru)
	if err != nil {
		
		return err
	}
	return nil
}

func (c *kasServices) Delete(id int) error {
	err := c.kasRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *kasServices) GetList(page, pageSize int) ([]models.KasGuruResponse, int, error) {
    KasGurus, totalPage, err := c.kasRepo.GetList(page, pageSize)
    if err != nil {
        return nil, 0, err
    }
    return KasGurus, totalPage, nil
}

func(c *kasServices) Search(nama, tanggal string) ([]models.KasGuruResponse, error){
	KasGuru, err := c.kasRepo.Search(nama, tanggal)
	if err != nil {
        return nil, err
    }
	return KasGuru, nil
}

func (c *kasServices) GetByID(id int) (*models.HistoryPembayaranKas, error){
	kas, err := c.kasRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return kas, nil
}

func (c *kasServices) GetAmbilByID(id int) (*models.HistoryPengambilanKas, error){
	kas, err := c.kasRepo.GetAmbilByID(id)
	if err != nil {
		return nil, err
	}
	return kas, nil
}
