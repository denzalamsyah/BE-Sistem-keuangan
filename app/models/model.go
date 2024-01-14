package models

import "time"

type Credential struct {
	Host         string
	Username     string
	Password     string
	DatabaseName string
	Port         int
	Schema       string
}
type Siswa struct {
	ID           int    `gorm:"primaryKey" json:"id"`
	Nama         string `gorm:"type:varchar(35)" json:"nama"`
	NISN         int    `json:"nisn"`
	KelasID      int    `gorm:"foreignKey:FkId_kelas;references:IDKelas" json:"id_kelas"`
	Kelas        Kelas  `json:"kelas"`
	AgamaID      int    `gorm:"foreignKey:FkId_agama;references:IDAgama" json:"id_agama"`
	Agama        Agama  `json:"agama"`
	TempatLahir  string `gorm:"type:varchar(15)" json:"tempat_lahir"`
	TanggalLahir string `json:"tanggal_lahir"`
	GenderID     int    `gorm:"foreignKey:FkId_gender;references:IDGender" json:"id_gender"`
	Gender       Gender `json:"gender"`
	NamaAyah     string `gorm:"type:varchar(35)" json:"nama_ayah"`
	NamaIbu      string `gorm:"type:varchar(35)" json:"nama_ibu"`
	NomorTelepon int    `json:"nomor_telepon"`
	Email        string `gorm:"type:varchar(35)" json:"email"`
	Alamat       string `gorm:"type:varchar(100)" json:"alamat"`
}

type HistoryPembayaran struct {
	ID             int         `gorm:"primaryKey" json:"id"`
	SiswaID        int         `gorm:"foreignKey:FkId_siswa;references:ID" json:"id_siswa"`
	Siswa          Siswa       `json:"siswa"`
	Nama_transaksi string      `gorm:"type:varchar(35)" json:"nama_transaksi"`
	Biaya          int         `json:"biaya"`
	Tanggal        string      `json:"tanggal"`
	PenerimaID     int         `gorm:"foreignKey:FkId_stakeholder;references:ID" json:"id_penerima"`
	Penerima       Stakeholder `json:"penerima"`
}

type ReportPembayaran struct {
	ID              int    `gorm:"primaryKey" json:"id"`
	SiswaID         int    `gorm:"foreignKey:FkId_siswa;references:ID" json:"id_siswa"`
	Siswa           Siswa  `json:"siswa"`
	Tahun_pelajaran string `json:"tahun_pelajaran"`
}

type Pemasukan struct {
	ID      int    `gorm:"primaryKey" json:"id"`
	Nama    string `json:"nama_pemasukan"`
	Tanggal string `json:"tanggal"`
	Jumlah  int    `json:"total_jumlah"`
}

type Pengeluaran struct {
	ID      int    `gorm:"primaryKey" json:"id"`
	Nama    string `gorm:"type:varchar(35)" json:"nama_pemasukan"`
	Tanggal string `json:"tanggal"`
	Jumlah  int    `json:"total_jumlah"`
}

type PembayaranSPP struct {
	ID         int         `gorm:"primaryKey" json:"id"`
	SiswaID    int         `gorm:"foreignKey:FkId_siswa;references:ID" json:"id_siswa"`
	Siswa      Siswa       `json:"siswa"`
	Tanggal    string      `json:"tanggal"`
	Jumlah     int         `json:"total_jumlah"`
	PenerimaID int         `gorm:"foreignKey:FkId_stakeholder;references:ID" json:"id_penerima"`
	Penerima   Stakeholder `json:"penerima"`
}

type PembayaranSemester struct {
	ID             int         `gorm:"primaryKey" json:"id"`
	SiswaID        int         `gorm:"foreignKey:FkId_siswa;references:ID" json:"id_siswa"`
	Siswa          Siswa       `json:"siswa"`
	Nama_transaksi string      `gorm:"type:varchar(35)" json:"nama_transaksi"`
	Tanggal        string      `json:"tanggal"`
	Jumlah         int         `json:"total_jumlah"`
	PenerimaID     int         `gorm:"foreignKey:FkId_stakeholder;references:ID" json:"id_penerima"`
	Penerima       Stakeholder `json:"penerima"`
}
type Stakeholder struct {
	ID           int    `gorm:"primaryKey" json:"id"`
	Nama         string `gorm:"type:varchar(35)" json:"nama"`
	NIP          int    `json:"nip"`
	AgamaID      int    `gorm:"foreignKey:FkId_agama;references:IDAgama" json:"id_agama"`
	Agama        Agama  `json:"agama"`
	Jabatan      string `gorm:"type:varchar(25)" json:"jabatan"`
	TempatLahir  string `gorm:"type:varchar(15)" json:"tempat_lahir"`
	TanggalLahir string `json:"tanggal_lahir"`
	GenderID     int    `gorm:"foreignKey:FkId_gender;references:IDGender" json:"id_gender"`
	Gender       Gender `json:"gender"`
	NomorTelepon int    `json:"nomor_telepon"`
	Email        string `gorm:"type:varchar(35)" json:"email"`
	Alamat       string `gorm:"type:varchar(100)" json:"alamat"`
}

type Login struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"type:varchar(35)" json:"email"`
	Password string `gorm:"type:varchar(15)" json:"password"`
}

type Kelas struct {
	IDKelas int    `gorm:"primaryKey" json:"id"`
	Kelas   string `gorm:"type:varchar(15)" json:"kelas"`
}

type Gender struct {
	IDGender int    `gorm:"primaryKey" json:"id"`
	Nama     string `gorm:"type:varchar(15)" json:"nama"`
}

type Agama struct {
	IDAgama int    `gorm:"primaryKey" json:"id"`
	Nama    string `gorm:"type:varchar(15)" json:"nama"`
}

type Session struct {
	ID     int       `gorm:"primaryKey" json:"id"`
	Token  string    `json:"token"`
	Email  string    `json:"email"`
	Expiry time.Time `json:"expiry"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Email     string    `json:"email" gorm:"type:varchar(255);not null"`
	Password  string    `json:"-" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}