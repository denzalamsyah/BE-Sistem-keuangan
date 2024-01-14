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
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type APIHandler struct {
	UserAPIHandler controllers.UserAPI
	SiswaAPIHandler controllers.SiswaAPI
	StakeAPIHandler controllers.StakeAPI
	SppAPIHandler controllers.SppAPI
	SemesterAPIHandler controllers.SemesterAPI
	
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

		conn.AutoMigrate(&models.User{}, &models.Session{}, &models.Siswa{}, &models.Stakeholder{}, &models.HistoryPembayaran{}, &models.ReportPembayaran{}, &models.Login{}, &models.Pemasukan{}, &models.PembayaranSPP{}, &models.PembayaranSemester{}, &models.Pengeluaran{}, &models.Gender{}, &models.Agama{}, models.Kelas{} )

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
	sppRepo := repository.NewSPPRepo(db)
	semesterRepo := repository.NewSemesterRepo(db)
	

	userService := services.NewUserService(userRepo, sessionRepo)
	siswaService := services.NewSiswaService(siswaRepo)
	stakeService := services.NewStakeService(stakeRepo)
	sppService := services.NewSPPService(sppRepo)
	semesterService := services.NewSemesterService(semesterRepo)

	userAPIHandler := controllers.NewUserAPI(userService)
	siswaAPI := controllers.NewSiswaAPI(siswaService)
	stakeAPI := controllers.NewStakeAPI(stakeService)
	sppAPI := controllers.NewSPPAPI(sppService)
	semesterAPI := controllers.NewSemesterAPI(semesterService)
	

	apiHandler := APIHandler{
	UserAPIHandler: userAPIHandler,
	SiswaAPIHandler: siswaAPI,
	StakeAPIHandler: stakeAPI,
	SppAPIHandler: sppAPI,
	SemesterAPIHandler: semesterAPI,
	
}
	
version := gin.Group("/api")
{
	user := version.Group("/user")
		{
			user.POST("/login", apiHandler.UserAPIHandler.Login)
			user.Use(middleware.Auth())
		}
	siswa := version.Group("/siswa")
	{
		siswa.Use(middleware.Auth())
		siswa.POST("/", apiHandler.SiswaAPIHandler.AddSiswa)
		siswa.PUT("/:id", apiHandler.SiswaAPIHandler.Update)
		siswa.DELETE("/:id", apiHandler.SiswaAPIHandler.Delete)
		siswa.GET("/:id", apiHandler.SiswaAPIHandler.GetByID)
		siswa.GET("/", apiHandler.SiswaAPIHandler.GetList)
	}
	stake := version.Group("/stake")
	{
		stake.Use(middleware.Auth())
		stake.POST("/", apiHandler.StakeAPIHandler.AddStake)
		stake.PUT("/:id", apiHandler.StakeAPIHandler.Update)
		stake.DELETE("/:id", apiHandler.StakeAPIHandler.Delete)
		stake.GET("/:id", apiHandler.StakeAPIHandler.GetByID)
		stake.GET("/", apiHandler.StakeAPIHandler.GetList)
	}

	Spp := version.Group("/spp")
	{
		Spp.Use(middleware.Auth())
		Spp.POST("/", apiHandler.SppAPIHandler.AddSPP)
		Spp.PUT("/:id", apiHandler.SppAPIHandler.Update)
		Spp.DELETE("/:id", apiHandler.SppAPIHandler.Delete)
		Spp.GET("/:id", apiHandler.SppAPIHandler.GetByID)
		Spp.GET("/", apiHandler.SppAPIHandler.GetList)
	}
	
	Semester := version.Group("/semester")
	{
		Semester.Use(middleware.Auth())
		Semester.POST("/", apiHandler.SemesterAPIHandler.AddSemester)
		Semester.PUT("/:id", apiHandler.SemesterAPIHandler.Update)
		Semester.DELETE("/:id", apiHandler.SemesterAPIHandler.Delete)
		Semester.GET("/:id", apiHandler.SemesterAPIHandler.GetByID)
		Semester.GET("/", apiHandler.SemesterAPIHandler.GetList)
	}
	

}
return gin

}