package models

import (
	"time"
)

type Credential struct {
	Host         string
	Username     string
	Password     string
	DatabaseName string
	Port         int
	Schema       string
}
type Siswa struct {
	Nisn 		 string `gorm:"type:varchar(10);primaryKey" json:"nisn" form:"nisn"`
	Nama         string `gorm:"type:varchar(35)" form:"nama"`
	KelasID      string `gorm:"type:varchar(10);foreignKey:FkId_kelas;references:KodeKelas" form:"kelas"`
	Kelas        Kelas  `json:"kelas"`
	JurusanID  string   `gorm:"type:varchar(10);foreignKey:FkId_jurusan;references:KodeJurusan" form:"jurusan"`
	Jurusan      Jurusan`json:"jurusan"`
	AgamaID      int    `gorm:"foreignKey:FkId_agama;references:IDAgama" json:"id_agama" form:"id_agama"`
	Agama        Agama  `json:"agama"`
	TempatLahir  string `gorm:"type:varchar(15)" json:"tempat_lahir" form:"tempat_lahir"`
	TanggalLahir string `gorm:"type:varchar(15)" json:"tanggal_lahir"  form:"tanggal_lahir"`
	GenderID     int    `gorm:"foreignKey:FkId_gender;references:IDGendeSr" json:"id_gender" form:"id_gender"`	
	Gender       Gender `json:"gender"`
	NamaAyah     string `gorm:"type:varchar(35)" json:"nama_ayah" form:"nama_ayah"`
	NamaIbu      string `gorm:"type:varchar(35)" json:"nama_ibu" form:"nama_ibu"`
	NomorTelepon string `gorm:"type:varchar(15)" json:"nomor_telepon" form:"nomor_telepon"`
	Angkatan     string `gorm:"type:varchar(10)" json:"angkatan" form:"angkatan"`
	Email        string `gorm:"type:varchar(35)" json:"email" form:"email"`
	Alamat       string `gorm:"type:varchar(100)" json:"alamat" form:"alamat"`
	Gambar string `gorm:"type:varchar(255)" json:"gambar" form:"gambar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PembayaranSemester struct {
	ID             int         `gorm:"primaryKey" json:"id"`
	SiswaID        string      `gorm:"type:varchar(10);foreignKey:FkId_siswa;references:Nisn" json:"id_siswa"`
	Siswa          Siswa       `json:"siswa"`
	TransaksiID    int         `gorm:"foreignKey:FkId_transaksi;references:ID" json:"id_transaksi"`
	Transaksi  Transaksi   `json:"transaksi"`
	Bulan       string      `gorm:"type:varchar(15)" json:"bulan"`
	Semester       string      `gorm:"type:varchar(15)" json:"semester"`
	Tanggal        string      `gorm:"type:varchar(15)" json:"tanggal"`
	Jumlah         int         `json:"total_jumlah"`
	Status         string      `gorm:"type:varchar(15)" json:"status"`
	Biaya  int  `json:"biaya"`
	Tunggakan int  `json:"tunggakan"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type HistoryPembayaran struct {
	ID             int    `json:"id"`
	Siswa          string `gorm:"type:varchar(35)" json:"siswa"`
	NISN           string `gorm:"type:varchar(10)" json:"nisn"`
	Nama_transaksi string `gorm:"type:varchar(35)" json:"nama_transaksi"`
	Biaya          int    `json:"biaya"`
	Tanggal        string `gorm:"type:varchar(15)" json:"tanggal"`
	Status         string `gorm:"type:varchar(15)" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Transaksi struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	Nama      string `gorm:"type:varchar(35)" json:"nama"`
	Jumlah 	int `json:"jumlah"`
	Kategori 	string `gorm:"type:varchar(5)" json:"kategori"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Kelas struct {
	KodeKelas string    `gorm:"primaryKey" json:"kode_kelas"`
	Kelas         string `gorm:"type:varchar(35)" json:"kelas"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Jurusan struct{
	KodeJurusan string    `gorm:"primaryKey" json:"kode_jurusan" `
	Jurusan    string `gorm:"type:varchar(35)" json:"jurusan" `
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Guru struct {
	Nip        string    `gorm:"type:varchar(10);primaryKey" json:"nip" form:"nip"`
	Nama         string `gorm:"type:varchar(35)" json:"nama" form:"nama"`
	AgamaID      int    `gorm:"foreignKey:FkId_agama;references:IDAgama" json:"id_agama" form:"id_agama"`
	Agama        Agama  `json:"agama"`
	JabatanID    int    `gorm:"foreignKey:FkId_jabatan;references:IDJabatan" json:"id_jabatan " form:"id_jabatan"`
	Jabatan      Jabatan  `json:"jabatan"`
	TempatLahir  string `gorm:"type:varchar(15)" json:"tempat_lahir" form:"tempat_lahir"`
	TanggalLahir string `gorm:"type:varchar(15)" json:"tanggal_lahir" form:"tanggal_lahir"`
	GenderID     int    `gorm:"foreignKey:FkId_gender;references:IDGender" json:"id_gender" form:"id_gender"`
	Gender       Gender `json:"gender"`
	NomorTelepon string    `gorm:"type:varchar(15)" json:"nomor_telepon" form:"nomor_telepon"`
	Email        string `gorm:"type:varchar(35)" json:"email" form:"email"`
	Alamat       string `gorm:"type:varchar(100)" json:"alamat" form:"alamat"`
	Gambar string `gorm:"type:varchar(255)" json:"gambar" form:"gambar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


type Pemasukan struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	IDBayar   int    `gorm:"foreignKey:ID" json:"id_bayar"`
	Nama      string `gorm:"type:varchar(35)" json:"nama_pemasukan"`
	Tanggal   string `gorm:"type:varchar(15)" json:"tanggal"`
	Jumlah    int    `json:"total_jumlah"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Pemasukanlainnya struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	Nama      string `gorm:"type:varchar(35)" json:"nama_pemasukan"`
	Tanggal   string `gorm:"type:varchar(15)" json:"tanggal"`
	Jumlah    int    `json:"total_jumlah"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Pengeluaran struct {
	ID      int    `gorm:"primaryKey" json:"id"`
	Nama    string `gorm:"type:varchar(35)" json:"nama_pengeluaran"`
	Tanggal string `gorm:"type:varchar(15)" json:"tanggal"`
	Jumlah  int    `json:"total_jumlah"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}



type Login struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"type:varchar(35)" json:"email"`
	Password string `gorm:"type:varchar(15)" json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


type Gender struct {
	IDGender int    `gorm:"primaryKey" json:"id"`
	Nama     string `gorm:"type:varchar(15)" json:"nama"`
	
}

type Agama struct {
	IDAgama int    `gorm:"primaryKey" json:"id"`
	Nama    string `gorm:"type:varchar(30)" json:"nama"`
}

type Session struct {
	ID     int       `gorm:"primaryKey" json:"id"`
	Token  string    `json:"token"`
	Email  string    `json:"email"`
	Expiry time.Time `json:"expiry"`
}

type UserLogin struct {
	Email    string `gorm:"type:varchar(35)" json:"email" binding:"required"`
	Password string `gorm:"type:varchar(15)" json:"password" binding:"required"`
}

type UserRegister struct {
	Email    string `gorm:"type:varchar(35)" json:"email" binding:"required"`
	Password string `gorm:"type:varchar(15)" json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password"`

}
type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Email     string    `json:"email" gorm:"type:varchar(55);not null"`
	Password  string    `json:"-" gorm:"type:varchar(105);not null"`
	ConfirmPassword  string `json:"-" gorm:"type:varchar(105);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}



type Jabatan struct{
	IDJabatan int    `gorm:"primaryKey" json:"id"`
	Nama    string `gorm:"type:varchar(30)" json:"nama"`
}

type ResetToken struct{
	Email   string `json:"email" binding:"required"`
	TokenHash   string `json:"token" binding:"required"`
	CreatedAt time.Time
	ExpirationTime time.Time
}

type KasGuru struct{
	ID int  `gorm:"primaryKey" json:"id"`
	GuruID string `gorm:"type:varchar(10);foreignKey:FkId_stakeholder;references:Nip" json:"id_guru"`
	Guru Guru `json:"nama_guru"`
	Jumlah int `json:"jumlah_bayar"`
	Saldo int `json:"saldo"`
	TanggalBayar string `gorm:"type:varchar(15)" json:"tanggal_bayar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PengambilanKas struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	GuruID      string      `gorm:"type:varchar(10)" json:"nip"`
	Nama      string    `gorm:"type:varchar(35)" json:"nama"`
	JumlahAmbil int       `json:"jumlah_ambil"`
	TanggalAmbil string `gorm:"type:varchar(15)" json:"tanggal_ambil"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TotalSaldoKas struct {
	Nama     string `json:"nama"`
	NIP      string    `gorm:"type:varchar(10)" json:"nip"`
	TotalKas int    `json:"total_kas"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

