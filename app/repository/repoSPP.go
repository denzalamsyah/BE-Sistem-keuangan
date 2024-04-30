package repository

// type SPPRepository interface {
// 	Store(PembayaranSPP *models.PembayaranSPP) error
// 	Update(id int, PembayaranSPP models.PembayaranSPP) error
// 	Delete(id int) error
// 	GetByID(id int) (*models.PembayaranSPPResponse, error)
// 	GetList() ([]models.PembayaranSPPResponse, error)
// }
// type sppRepository struct{
// 	db *gorm.DB
// }

// func NewSPPRepo(db *gorm.DB) *sppRepository{
// 	return &sppRepository{db}
// }

// func generateRandomCode() string {
//     const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
//     seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

//     b := make([]byte, 5)
//     for i := range b {
//         b[i] = charset[seededRand.Intn(len(charset))]
//     }
//     return string(b)
// }
// func (c *sppRepository) Store(PembayaranSPP *models.PembayaranSPP) error {
//     tx := c.db.Begin()

//     // Create PembayaranSPP
//     if err := tx.Create(PembayaranSPP).Error; err != nil {
//         tx.Rollback()
//         return err
//     }

//     // Retrieve Transaksi data
//     var transaksi models.Transaksi
//     if err := tx.First(&transaksi, PembayaranSPP.TransaksiID).Error; err != nil {
//         tx.Rollback()
//         return err
//     }

//     // Create Pemasukan
//     pemasukanSPP := models.Pemasukan{
//         IDBayar: PembayaranSPP.ID,
//         Nama:    transaksi.Nama,
//         Tanggal: PembayaranSPP.Tanggal,
//         Jumlah:  PembayaranSPP.Jumlah,
//     }
//     if err := tx.Create(&pemasukanSPP).Error; err != nil {
//         tx.Rollback()
//         return err
//     }

//     return tx.Commit().Error
// }

// func (c *sppRepository) Update(id int, PembayaranSPP models.PembayaranSPP) error {
//     tx := c.db.Begin()

//     var pembayaranSPP models.PembayaranSPP

//     if err := tx.Where("id = ?", id).First(&pembayaranSPP).Error; err != nil {
//         tx.Rollback()
//         return err
//     }

//     if err := tx.Model(&pembayaranSPP).Updates(&PembayaranSPP).Error; err != nil {
//         tx.Rollback()
//         return err
//     }

//     if err := tx.Model(&models.Pemasukan{}).
//         Where("id_bayar = ?", pembayaranSPP.ID).
//         Updates(map[string]interface{}{
//             "nama": pembayaranSPP.Transaksi.Nama,
//             "tanggal":        pembayaranSPP.Tanggal,
//             "jumlah":   pembayaranSPP.Jumlah,
//         }).Error; err != nil {
//         tx.Rollback()
//         return err
//     }

//     return tx.Commit().Error
// }

// func (c *sppRepository) Delete(id int) error {
// 	tx := c.db.Begin()

// 	var pembayaranSPP models.PembayaranSPP

// 	if err := tx.Where("id = ?", id).First(&pembayaranSPP).Error; err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	if err := tx.Delete(&pembayaranSPP).Error; err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	if err := tx.Where("nama = ? AND tanggal = ? AND jumlah = ? AND id_bayar = ?", pembayaranSPP.Transaksi.Nama, pembayaranSPP.Tanggal, pembayaranSPP.Jumlah, pembayaranSPP.ID).Delete(&models.Pemasukan{}).Error; err != nil {
// 		tx.Rollback()
// 		return err
// 	}
// 	return tx.Commit().Error
// }

// func (c *sppRepository) GetByID(id int) (*models.PembayaranSPPResponse, error) {
// 	var PembayaranSPP models.PembayaranSPP
// 	err := c.db.Preload("Siswa").Preload("Penerima").Preload("Transaksi").Where("id = ?", id).First(&PembayaranSPP).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	PembayaranSPPResponse := models.PembayaranSPPResponse{
// 		ID:             PembayaranSPP.ID,
// 		Siswa:          PembayaranSPP.Siswa.Nama,
// 		Transaksi:  PembayaranSPP.Transaksi.Nama,
// 		Bulan:          PembayaranSPP.Bulan,
// 		Semester:       PembayaranSPP.Semester,
// 		TahunAjar:      PembayaranSPP.TahunAjar,
// 		Tanggal:        PembayaranSPP.Tanggal,
// 		Jumlah:         PembayaranSPP.Jumlah,
// 		Penerima:       PembayaranSPP.Penerima.Nama,
// 		Status:         PembayaranSPP.Status,
// 	}
// 	return &PembayaranSPPResponse, nil
// }

// func (c *sppRepository) GetList() ([]models.PembayaranSPPResponse, error) {
// 	var PembayaranSPP []models.PembayaranSPP
// 	err := c.db.Preload("Siswa").Preload("Penerima").Preload("Transaksi").Find(&PembayaranSPP).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	var PembayaranSPPResponse []models.PembayaranSPPResponse
// 	for _, s := range PembayaranSPP{
// 		PembayaranSPPResponse = append(PembayaranSPPResponse, models.PembayaranSPPResponse{
// 			ID:             s.ID,
// 			Siswa:          s.Siswa.Nama,
// 			Transaksi:  s.Transaksi.Nama,
// 			Bulan:          s.Bulan,
// 			Semester:       s.Semester,
// 			TahunAjar:      s.TahunAjar,
// 			Tanggal:        s.Tanggal,
// 			Jumlah:         s.Jumlah,
// 			Penerima:       s.Penerima.Nama,
// 			Status:         s.Status,
// 		})
// 	}
// 	return PembayaranSPPResponse, nil
// }
