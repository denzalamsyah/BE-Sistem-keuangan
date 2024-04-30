package services

// import (
// 	"github.com/denzalamsyah/simak/app/models"
// 	"github.com/denzalamsyah/simak/app/repository"
// )

// type SPPServices interface {
// 	Store(PembayaranSPP *models.PembayaranSPP) error
// 	Update(id int, PembayaranSPP models.PembayaranSPP) error
// 	Delete(id int) error
// 	GetByID(id int) (*models.PembayaranSPPResponse, error)
// 	GetList() ([]models.PembayaranSPPResponse, error)
// }

// type sppServices struct {
// 	sppRepo repository.SPPRepository
// }

// func NewSPPService(sppRepo repository.SPPRepository) SPPServices{
// 	return &sppServices{sppRepo}
// }

// func (c *sppServices) Store(PembayaranSPP *models.PembayaranSPP) error {
// 	err := c.sppRepo.Store(PembayaranSPP)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *sppServices) Update(id int, PembayaranSPP models.PembayaranSPP) error {
// 	err := c.sppRepo.Update(id, PembayaranSPP)
// 	if err != nil {

// 		return err
// 	}
// 	return nil
// }

// func (c *sppServices) Delete(id int) error {
// 	err := c.sppRepo.Delete(id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *sppServices) GetByID(id int) (*models.PembayaranSPPResponse, error) {
// 	PembayaranSPP, err := c.sppRepo.GetByID(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return PembayaranSPP, nil
// }

// func (c *sppServices) GetList() ([]models.PembayaranSPPResponse, error) {
// 	PembayaranSPPs, err := c.sppRepo.GetList()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return PembayaranSPPs, nil
// }
