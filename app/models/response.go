package models

type SiswaResponse struct {
	ID           int    `gorm:"primaryKey" json:"id"`
	Nama         string `gorm:"type:varchar(35)" json:"nama"`
	NISN         int    `json:"nisn"`
	Kelas        string `json:"kelas"`
	Jurusan      string `json:"jurusan"`
	Agama        string `json:"agama"`
	TempatLahir  string `gorm:"type:varchar(15)" json:"tempat_lahir"`
	TanggalLahir string `json:"tanggal_lahir"`
	Gender       string `json:"gender"`
	NamaAyah     string `gorm:"type:varchar(35)" json:"nama_ayah"`
	NamaIbu      string `gorm:"type:varchar(35)" json:"nama_ibu"`
	NomorTelepon int    `json:"nomor_telepon"`
	Email        string `gorm:"type:varchar(35)" json:"email"`
	Alamat       string `gorm:"type:varchar(100)" json:"alamat"`
}

type StakeholderResponse struct {
	ID           int    `gorm:"primaryKey" json:"id"`
	Nama         string `gorm:"type:varchar(35)" json:"nama"`
	NIP          int    `json:"nip"`
	Agama        string `json:"agama"`
	Jabatan      string `json:"jabatan"`
	TempatLahir  string `gorm:"type:varchar(15)" json:"tempat_lahir"`
	TanggalLahir string `json:"tanggal_lahir"`
	Gender       string `json:"gender"`
	NomorTelepon int    `json:"nomor_telepon"`
	Email        string `gorm:"type:varchar(35)" json:"email"`
	Alamat       string `gorm:"type:varchar(100)" json:"alamat"`
}
type HistoryPembayaran struct {
	Siswa          string `json:"siswa"`
	Nama_transaksi string `gorm:"type:varchar(35)" json:"nama_transaksi"`
	Biaya          int    `json:"biaya"`
	Tanggal        string `json:"tanggal"`
	Penerima       string `json:"penerima"`
}
type PembayaranSPPResponse struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	Siswa     string `json:"siswa"`
	Bulan     string `json:"bulan"`
	Semester  string `json:"semester"`
	Transaksi string `gorm:"type:varchar(35)" json:"nama_transaksi"`
	TahunAjar string `json:"tahun_ajar"`
	Tanggal   string `json:"tanggal"`
	Jumlah    int    `json:"total_jumlah"`
	Penerima  string `json:"penerima"`
	Status    string `json:"status"`
}

type PembayaranSemesterResponse struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	Siswa     string `json:"siswa"`
	TahunAjar string `json:"tahun_ajar"`
	Semester  string `json:"semester"`
	Transaksi string `gorm:"type:varchar(35)" json:"nama_transaksi"`
	Tanggal   string `json:"tanggal"`
	Jumlah    int    `json:"total_jumlah"`
	Penerima  string `json:"penerima"`
	Status    string `json:"status"`
}

type PemasukanResponse struct {
	ID      int    `gorm:"primaryKey" json:"id"`
	Nama    string `json:"nama_pemasukan"`
	Tanggal string `json:"tanggal"`
	Jumlah  int    `json:"total_jumlah"`
}

type Total struct {
	Pemasukan   int `json:"total_pemasukan"`
	Pengeluaran int `json:"total_pengeluaran"`
	Saldo       int `json:"total_saldo"`
}
