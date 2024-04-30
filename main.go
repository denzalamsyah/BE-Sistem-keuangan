package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/denzalamsyah/simak/app/controllers"
	"github.com/denzalamsyah/simak/app/initializers"
	"github.com/denzalamsyah/simak/app/middleware"
	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/repository"
	"github.com/denzalamsyah/simak/app/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type APIHandler struct {
	UserAPIHandler controllers.UserAPI
	SiswaAPIHandler controllers.SiswaAPI
	StakeAPIHandler controllers.StakeAPI
	// SppAPIHandler controllers.SppAPI
	SemesterAPIHandler controllers.SemesterAPI
	PemasukanAPIHandler controllers.PemasukanAPI
	PengeluaranAPIHandler controllers.PengeluaranAPI
	KelasAPIHandler controllers.KelasAPI
	JurusanAPIHandler controllers.JurusanAPI
	TransaksiAPIHandler controllers.TransaksiAPI
	ArisanAPIHandler controllers.ArisanAPI
	KasAPIHandler controllers.KasAPI
}
func main() {
	gin.SetMode(gin.ReleaseMode) //release

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		router := gin.New()
		db := initializers.NewDB()
		router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("[%s] \"%s %s %s\"\n",
				param.TimeStamp.Format(time.RFC822),
				param.Method,
				param.Path,
				param.ErrorMessage,
			)
		}))
		router.Use(gin.Recovery())

		dbCredential := models.Credential{
			Host:         "localhost",
			Username:     "postgres",
			Password:     "rizwan123",
			DatabaseName: "simaks",
			Port:         5432,
			Schema:       "public",
		}

		conn, err := db.ConnectToDB(&dbCredential)
		if err != nil {
			panic(err)
		}

		conn.AutoMigrate(&models.User{}, &models.Session{}, &models.Siswa{}, &models.Stakeholder{}, &models.Login{}, &models.Pemasukan{}, &models.Pemasukanlainnya{}, &models.HistoryPembayaran{}, &models.PembayaranSemester{}, &models.Pengeluaran{}, &models.Gender{}, &models.Agama{}, models.Kelas{}, &models.Transaksi{}, &models.ResetToken{}, &models.Arisan{}, &models.KasGuru{} )

		router = RunServer(conn, router)
	

		fmt.Println("Server is running on port 8080")
		err = router.Run(":8080")
		if err != nil {
			panic(err)
		}

	}()

	wg.Wait()
}


func RunServer(db *gorm.DB,  gin *gin.Engine) *gin.Engine{
	userRepo := repository.NewUserRepo(db)
	sessionRepo := repository.NewSessionsRepo(db)
	siswaRepo := repository.NewSiswaRepo(db)
	stakeRepo := repository.NewStakeholderRepo(db)
	// sppRepo := repository.NewSPPRepo(db)
	semesterRepo := repository.NewSemesterRepo(db)
	pemasukanRepo := repository.NewPemasukanRepo(db)
	pengeluaranRepo := repository.NewPengeluaranRepo(db)
	kelasRepo := repository.NewKelasRepo(db)
	jurusanRepo := repository.NewJurusanRepo(db)
	transaksiRepo := repository.NewTransaksiRepo(db)
	arisanRepo := repository.NewArisanRepo(db)
	kasRepo := repository.NewKasRepo(db)

	userService := services.NewUserService(userRepo, sessionRepo,)
	siswaService := services.NewSiswaService(siswaRepo)
	stakeService := services.NewStakeService(stakeRepo)
	// sppService := services.NewSPPService(sppRepo)
	semesterService := services.NewSemesterService(semesterRepo)
	pemasukanService := services.NewPemasukanService(pemasukanRepo)
	pengeluaranService := services.NewPengeluaranService(pengeluaranRepo)
	kelasService := services.NewKelasService(kelasRepo)
	jurusanService := services.NewJurusanService(jurusanRepo)
	transaksiService := services.NewTransaksiService(transaksiRepo)
	arisanService := services.NewArisanService(arisanRepo)
	kasService := services.NewKasService(kasRepo)

	userAPIHandler := controllers.NewUserAPI(userService)
	siswaAPI := controllers.NewSiswaAPI(siswaService)
	stakeAPI := controllers.NewStakeAPI(stakeService)
	// sppAPI := controllers.NewSPPAPI(sppService)
	semesterAPI := controllers.NewSemesterAPI(semesterService)
	pemasukanAPI := controllers.NewPemasukanAPI(pemasukanService)
	pengeluaranAPI := controllers.NewPengeluaranAPI(pengeluaranService)
	kelasAPI := controllers.NewKelasAPI(kelasService)
	jurusanAPI := controllers.NewJurusanAPI(jurusanService)
	transaksiAPI := controllers.NewTransaksiAPI(transaksiService)
	arisanAPI := controllers.NewArisanAPI(arisanService)
	kasAPI := controllers.NewKasAPI(kasService)

	apiHandler := APIHandler{
	UserAPIHandler: userAPIHandler,
	SiswaAPIHandler: siswaAPI,
	StakeAPIHandler: stakeAPI,
	// SppAPIHandler: sppAPI,
	SemesterAPIHandler: semesterAPI,
	PemasukanAPIHandler: pemasukanAPI,
	PengeluaranAPIHandler: pengeluaranAPI,
	KelasAPIHandler: kelasAPI,
	JurusanAPIHandler: jurusanAPI,
	TransaksiAPIHandler: transaksiAPI,
	ArisanAPIHandler: arisanAPI,
	KasAPIHandler: kasAPI,
}

gin.Use(cors.New(cors.Config{
	AllowOrigins:     []string{"*"},
	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
	AllowCredentials: true,
	MaxAge:           12 * time.Hour,
}))
	
version := gin.Group("/api")
{
	user := version.Group("/user")
		{
			user.POST("/login", apiHandler.UserAPIHandler.Login)
			user.POST("/register", apiHandler.UserAPIHandler.Register)
			user.POST("/reset-password/request", apiHandler.UserAPIHandler.RequestResetToken)
			user.POST("/reset-password/reset", apiHandler.UserAPIHandler.ResetPassword)
			user.Use(middleware.Auth())
			
		}
	siswa := version.Group("/siswa")
	{
		siswa.Use(middleware.Auth())
		siswa.POST("/", apiHandler.SiswaAPIHandler.AddSiswa)
		siswa.PUT("/:id", apiHandler.SiswaAPIHandler.Update)
		siswa.DELETE("/:id", apiHandler.SiswaAPIHandler.Delete)
		siswa.GET("/:id", apiHandler.SiswaAPIHandler.GetByID)
		siswa.GET("/histori/:id", apiHandler.SiswaAPIHandler.History)
		siswa.GET("/", apiHandler.SiswaAPIHandler.GetList)
		siswa.GET("/gender", apiHandler.SiswaAPIHandler.GetTotalGenderCount)
		siswa.GET("/search", apiHandler.SiswaAPIHandler.Search)
		siswa.GET("/export", apiHandler.SiswaAPIHandler.ExportSiswa)
	}
	stake := version.Group("/stake")
	{
		stake.Use(middleware.Auth())
		stake.POST("/", apiHandler.StakeAPIHandler.AddStake)
		stake.PUT("/:id", apiHandler.StakeAPIHandler.Update)
		stake.DELETE("/:id", apiHandler.StakeAPIHandler.Delete)
		stake.GET("/:id", apiHandler.StakeAPIHandler.GetByID)
		stake.GET("/", apiHandler.StakeAPIHandler.GetList)
		stake.GET("/gender", apiHandler.StakeAPIHandler.GetTotalGenderCount)
		stake.GET("/search", apiHandler.StakeAPIHandler.Search)
	}

	// Spp := version.Group("/spp")
	// {
	// 	Spp.Use(middleware.Auth())
	// 	Spp.POST("/", apiHandler.SppAPIHandler.AddSPP)
	// 	Spp.PUT("/:id", apiHandler.SppAPIHandler.Update)
	// 	Spp.DELETE("/:id", apiHandler.SppAPIHandler.Delete)
	// 	Spp.GET("/:id", apiHandler.SppAPIHandler.GetByID)
	// 	Spp.GET("/", apiHandler.SppAPIHandler.GetList)
	// }
	
	Semester := version.Group("/semester")
	{
		Semester.Use(middleware.Auth())
		Semester.POST("/", apiHandler.SemesterAPIHandler.AddSemester)
		Semester.PUT("/:id", apiHandler.SemesterAPIHandler.Update)
		Semester.DELETE("/:id", apiHandler.SemesterAPIHandler.Delete)
		Semester.GET("/:id", apiHandler.SemesterAPIHandler.GetByID)
		Semester.GET("/", apiHandler.SemesterAPIHandler.GetList)
		Semester.GET("/search", apiHandler.SemesterAPIHandler.Search)
	}
	pemasukan := version.Group("/pemasukan")
	{
		pemasukan.Use(middleware.Auth())
		pemasukan.GET("/", apiHandler.PemasukanAPIHandler.FindAll)
		pemasukan.GET("/total", apiHandler.PemasukanAPIHandler.TotalKeuangan)
		pemasukan.GET("/:id", apiHandler.PemasukanAPIHandler.GetByID)
		pemasukan.POST("/", apiHandler.PemasukanAPIHandler.Add)
		pemasukan.PUT("/:id", apiHandler.PemasukanAPIHandler.Update)
		pemasukan.DELETE("/:id", apiHandler.PemasukanAPIHandler.Delete)
		pemasukan.GET("/get", apiHandler.PemasukanAPIHandler.GetList)
		pemasukan.GET("/searchAll", apiHandler.PemasukanAPIHandler.SearchAll)
		pemasukan.GET("/search", apiHandler.PemasukanAPIHandler.Search)

	}
	pengeluaran := version.Group("/pengeluaran")
	{
		pengeluaran.Use(middleware.Auth())
		pengeluaran.GET("/:id", apiHandler.PengeluaranAPIHandler.GetByID)
		pengeluaran.POST("/", apiHandler.PengeluaranAPIHandler.Add)
		pengeluaran.PUT("/:id", apiHandler.PengeluaranAPIHandler.Update)
		pengeluaran.DELETE("/:id", apiHandler.PengeluaranAPIHandler.Delete)
		pengeluaran.GET("/", apiHandler.PengeluaranAPIHandler.GetList)
		pengeluaran.GET("/search", apiHandler.PengeluaranAPIHandler.Search)
	}

	kelas := version.Group("/kelas")
	{
		kelas.Use(middleware.Auth())
		kelas.POST("/", apiHandler.KelasAPIHandler.AddKelas)
		kelas.PUT("/:id", apiHandler.KelasAPIHandler.Update)
		kelas.DELETE("/:id", apiHandler.KelasAPIHandler.Delete)
		kelas.GET("/", apiHandler.KelasAPIHandler.GetList)
		kelas.GET("/total", apiHandler.KelasAPIHandler.GetTotalKelasCount)
		kelas.GET("/search", apiHandler.KelasAPIHandler.Search)
	}

	Jurusan := version.Group("/jurusan")
	{
		Jurusan.Use(middleware.Auth())
		Jurusan.POST("/", apiHandler.JurusanAPIHandler.AddJurusan)
		Jurusan.PUT("/:id", apiHandler.JurusanAPIHandler.Update)
		Jurusan.DELETE("/:id", apiHandler.JurusanAPIHandler.Delete)
		Jurusan.GET("/", apiHandler.JurusanAPIHandler.GetList)
		Jurusan.GET("/total", apiHandler.JurusanAPIHandler.GetTotalJurusanCount)
		Jurusan.GET("/search", apiHandler.JurusanAPIHandler.Search)
	}

	Transaksi := version.Group("/transaksi")
	{
		Transaksi.Use(middleware.Auth())
		Transaksi.POST("/", apiHandler.TransaksiAPIHandler.AddTransaksi)
		Transaksi.PUT("/:id", apiHandler.TransaksiAPIHandler.Update)
		Transaksi.DELETE("/:id", apiHandler.TransaksiAPIHandler.Delete)
		Transaksi.GET("/", apiHandler.TransaksiAPIHandler.GetList)
		Transaksi.GET("/search", apiHandler.TransaksiAPIHandler.Search)
	}

	Arisan := version.Group("/arisan")
	{
		Arisan.Use(middleware.Auth())
		Arisan.POST("/", apiHandler.ArisanAPIHandler.Store)
		Arisan.PUT("/:id", apiHandler.ArisanAPIHandler.Update)
		Arisan.DELETE("/:id", apiHandler.ArisanAPIHandler.Delete)
		Arisan.GET("/", apiHandler.ArisanAPIHandler.GetList)
		Arisan.GET("/search", apiHandler.ArisanAPIHandler.Search)
	}

	Kas := version.Group("/kas")
	{
		Kas.Use(middleware.Auth())
		Kas.POST("/", apiHandler.KasAPIHandler.Store)
		Kas.PUT("/:id", apiHandler.KasAPIHandler.Update)
		Kas.DELETE("/:id", apiHandler.KasAPIHandler.Delete)
		Kas.GET("/", apiHandler.KasAPIHandler.GetList)
		Kas.GET("/search", apiHandler.KasAPIHandler.Search)

	}
}
return gin

}