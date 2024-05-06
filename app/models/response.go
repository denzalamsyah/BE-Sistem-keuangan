package models

type SiswaResponse struct {
	NISN         int    `json:"nisn"`
	Nama         string `json:"nama"`
	Kelas        string `json:"kelas"`
	Jurusan      string `json:"jurusan"`
	Agama        string `json:"agama"`
	TempatLahir  string `json:"tempat_lahir"`
	TanggalLahir string `json:"tanggal_lahir"`
	Gender       string `json:"gender"`
	NamaAyah     string `json:"nama_ayah"`
	NamaIbu      string `json:"nama_ibu"`
	NomorTelepon int    `json:"nomor_telepon"`
	Angkatan     string `json:"angkatan" form:"angkatan"`
	Email        string `json:"email"`
	Alamat       string `json:"alamat"`
	Gambar       string `json:"gambar"`
}

type GuruResponse struct {
	Nip          int    `json:"nip"`
	Nama         string `json:"nama"`
	Agama        string `json:"agama"`
	Jabatan      string `json:"jabatan"`
	TempatLahir  string `json:"tempat_lahir"`
	TanggalLahir string `json:"tanggal_lahir"`
	Gender       string `json:"gender"`
	NomorTelepon int    `json:"nomor_telepon"`
	Email        string `json:"email"`
	Alamat       string `json:"alamat"`
	Gambar       string `json:"gambar"`
}
type HistoryPembayaran struct {
	ID             int    `json:"id"`
	Siswa          string `gorm:"type:varchar(35)" json:"siswa"`
	NISN           int    `json:"nisn"`
	Nama_transaksi string `gorm:"type:varchar(35)" json:"nama_transaksi"`
	Biaya          int    `json:"biaya"`
	Tanggal        string `gorm:"type:varchar(15)" json:"tanggal"`
	Status         string `gorm:"type:varchar(15)" json:"status"`
}
type PembayaranSPPResponse struct {
	ID        int    `json:"id"`
	Siswa     string `json:"siswa"`
	Bulan     string `json:"bulan"`
	Semester  string `json:"semester"`
	Transaksi string `json:"nama_transaksi"`
	TahunAjar string `json:"tahun_ajar"`
	Tanggal   string `json:"tanggal"`
	Jumlah    int    `json:"total_jumlah"`
	Penerima  string `json:"penerima"`
	Status    string `json:"status"`
}

type PembayaranSemesterResponse struct {
	ID        int    `json:"id"`
	Siswa     string `json:"siswa"`
	NISN      int    `json:"nisn"`
	TahunAjar string `json:"tahun_ajar"`
	Semester  string `json:"semester"`
	Transaksi string `json:"nama_transaksi"`
	Tanggal   string `json:"tanggal"`
	Jumlah    int    `json:"total_jumlah"`
	Status    string `json:"status"`
}

type PemasukanResponse struct {
	ID      int    `json:"id"`
	Nama    string `json:"nama_pemasukan"`
	Tanggal string `json:"tanggal"`
	Jumlah  int    `json:"total_jumlah"`
}

type Total struct {
	Pemasukan   int `json:"total_pemasukan"`
	Pengeluaran int `json:"total_pengeluaran"`
	Saldo       int `json:"total_saldo"`
}

// type KelasResponse struct {
// 	ID    int    `json:"kode_kelas"`
// 	Kelas string `json:"kelas"`
// }
// type JurusanResponse struct {
// 	ID   int    `json:"kode_jurusan"`
// 	Nama string `json:"nama"`
// }

type TransaksiResponse struct {
	ID          int    `json:"id"`
	Nama        string `json:"nama"`
	JumlahBayar int    `json:"jumlah_bayar"`
}

type KasGuruResponse struct {
	ID           int    `json:"id"`
	NIP          int    `json:"nip"`
	NamaGuru     string `json:"nama_guru"`
	Jumlah       int    `json:"jumlah_bayar"`
	Saldo        int    `json:"saldo"`
	TanggalBayar string `json:"tanggal_bayar"`
}

type HistoryPembayaranKas struct {
	ID           uint   `json:"id"`
	Nama         string `json:"nama"`
	NIP          int    `json:"nip"`
	Saldo        int    `json:"saldo"`
	Jumlah_Bayar int    `json:"jumlah_bayar"`
	Tanggal      string `json:"tanggal_bayar"`
}

type HistoryPengambilanKas struct {
	ID           uint   `json:"id"`
	NIP          uint   `json:"nip"`
	Nama         string `json:"nama"`
	JumlahAmbil  int    `json:"jumlah_ambil"`
	TanggalAmbil string `json:"tanggal_ambil"`
}