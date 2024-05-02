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
	Nisn int `gorm:"primaryKey;size:20" json:"nisn" form:"nisn"`
	Nama         string `gorm:"type:varchar(35)" json:"nama" form:"nama"`
	KelasID      int    `gorm:"foreignKey:FkId_kelas;references:IDKelas" json:"id_kelas" form:"id_kelas"`
	Kelas        Kelas  `json:"kelas"`
	JurusanID    int    `gorm:"foreignKey:FkId_jurusan;references:IDJurusan" json:"id_jurusan" form:"id_jurusan"`
	Jurusan      Jurusan  `json:"jurusan"`
	AgamaID      int    `gorm:"foreignKey:FkId_agama;references:IDAgama" json:"id_agama" form:"id_agama"`
	Agama        Agama  `json:"agama"`
	TempatLahir  string `gorm:"type:varchar(15)" json:"tempat_lahir" form:"tempat_lahir"`
	TanggalLahir string `gorm:"type:varchar(15)" json:"tanggal_lahir"  form:"tanggal_lahir"`
	GenderID     int    `gorm:"foreignKey:FkId_gender;references:IDGender" json:"id_gender" form:"id_gender"`	
	Gender       Gender `json:"gender"`
	NamaAyah     string `gorm:"type:varchar(35)" json:"nama_ayah" form:"nama_ayah"`
	NamaIbu      string `gorm:"type:varchar(35)" json:"nama_ibu" form:"nama_ibu"`
	NomorTelepon int    `json:"nomor_telepon" form:"nomor_telepon"`
	Angkatan     string ` gorm:"type:varchar(10)" json:"angkatan" form:"angkatan"`
	Email        string `gorm:"type:varchar(35)" json:"email" form:"email"`
	Alamat       string `gorm:"type:varchar(100)" json:"alamat" form:"alamat"`
	Gambar string `gorm:"type:varchar(255)" json:"gambar" form:"gambar"`
}
type Guru struct {
	Nip        int    `gorm:"primaryKey" json:"nip" form:"nip"`
	Nama         string `gorm:"type:varchar(35)" json:"nama" form:"nama"`
	AgamaID      int    `gorm:"foreignKey:FkId_agama;references:IDAgama" json:"id_agama" form:"id_agama"`
	Agama        Agama  `json:"agama"`
	JabatanID    int    `gorm:"foreignKey:FkId_jabatan;references:IDJabatan" json:"id_jabatan " form:"id_jabatan"`
	Jabatan      Jabatan  `json:"jabatan"`
	TempatLahir  string `gorm:"type:varchar(15)" json:"tempat_lahir" form:"tempat_lahir"`
	TanggalLahir string `gorm:"type:varchar(15)" json:"tanggal_lahir" form:"tanggal_lahir"`
	GenderID     int    `gorm:"foreignKey:FkId_gender;references:IDGender" json:"id_gender" form:"id_gender"`
	Gender       Gender `json:"gender"`
	NomorTelepon int    `json:"nomor_telepon" form:"nomor_telepon"`
	Email        string `gorm:"type:varchar(35)" json:"email" form:"email"`
	Alamat       string `gorm:"type:varchar(100)" json:"alamat" form:"alamat"`
	Gambar string `gorm:"type:varchar(255)" json:"gambar" form:"gambar"`
}


type Pemasukan struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	IDBayar   int    `gorm:"foreignKey:ID" json:"id_bayar"`
	Nama      string `gorm:"type:varchar(35)" json:"nama_pemasukan"`
	Tanggal   string `gorm:"type:varchar(15)" json:"tanggal"`
	Jumlah    int    `json:"total_jumlah"`
}
type Pemasukanlainnya struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	Nama      string `gorm:"type:varchar(35)" json:"nama_pemasukan"`
	Tanggal   string `gorm:"type:varchar(15)" json:"tanggal"`
	Jumlah    int    `json:"total_jumlah"`
}

type Pengeluaran struct {
	ID      int    `gorm:"primaryKey" json:"id"`
	Nama    string `gorm:"type:varchar(35)" json:"nama_pengeluaran"`
	Tanggal string `gorm:"type:varchar(15)" json:"tanggal"`
	Jumlah  int    `json:"total_jumlah"`
}



type PembayaranSemester struct {
	ID             int         `gorm:"primaryKey" json:"id"`
	SiswaID        int         `gorm:"foreignKey:FkId_siswa;references:Nisn" json:"id_siswa"`
	Siswa          Siswa       `json:"siswa"`
	TransaksiID    int         `gorm:"foreignKey:FkId_transaksi;references:ID" json:"id_transaksi"`
	Transaksi  Transaksi   `json:"transaksi"`
	Semester       string      `gorm:"type:varchar(15)" json:"semester"`
	TahunAjar      string      `gorm:"type:varchar(15)" json:"tahun_ajar"`
	Tanggal        string      `gorm:"type:varchar(15)" json:"tanggal"`
	Jumlah         int         `json:"total_jumlah"`
	Status         string      `gorm:"type:varchar(15)" json:"status"`
}


type Transaksi struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	Nama      string `gorm:"type:varchar(35)" json:"nama"`
	JumlahBayar int    `json:"jumlah_bayar"`
}

type Login struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"type:varchar(35)" json:"email"`
	Password string `gorm:"type:varchar(15)" json:"password"`
}

type Kelas struct {
	IDKelas int    `gorm:"primaryKey" json:"id"`
	Kelas   string `gorm:"type:varchar(35)" json:"kelas"`
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

type Jurusan struct{
	IDJurusan int    `gorm:"primaryKey" json:"id"`
	Nama    string `gorm:"type:varchar(35)" json:"nama"`
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

type PesertaArisan struct{
	ID int  `gorm:"primaryKey" json:"id"`
	IDPeserta int `gorm:"foreignKey:FkId_stakeholders;references:NIP" json:"id_peserta"`
	Peserta Guru  `json:"peserta"`
	NomorArisan int `json:"no_arisan"`
	Jumlah int `json:"jumlah_bayar"`
	TanggalBayar string `gorm:"type:varchar(15)" json:"tanggal_bayar"`
}

type Arisan struct{
	ID int  `gorm:"primaryKey" json:"id"`
	NamaArisan string  `gorm:"type:varchar(35)" json:"nama_arisan"`
	JumlahPeserta int `json:"jumlah_peserta"`
	TanggalMulai string `gorm:"type:varchar(15)" json:"tanggal_mulai"`
	TanggalBerakhir string `gorm:"type:varchar(15)" json:"tanggal_berakhir"`
}

type KasGuru struct{
	ID int  `gorm:"primaryKey" json:"id"`
	GuruID int `gorm:"foreignKey:FkId_stakeholder;references:Nip" json:"id_guru"`
	Guru Guru `json:"nama_guru"`
	Jumlah int `json:"jumlah_bayar"`
	TanggalBayar string `gorm:"type:varchar(15)" json:"tanggal_bayar"`
}

type PengambilanKas struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	GuruID      uint      `json:"nip"`
	Nama      string    `gorm:"type:varchar(35)" json:"nama"`
	JumlahAmbil int       `json:"jumlah_ambil"`
	TanggalAmbil string `gorm:"type:varchar(15)" json:"tanggal_ambil"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type HistoryPengambilanKas struct {
	NIP           uint   `json:"nip"`
	Nama          string `json:"nama"`
	JumlahAmbil   int    `json:"jumlah_ambil"`
	TanggalAmbil  string `json:"tanggal_ambil"`
}
