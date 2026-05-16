package routes

import (
	"net/http"
	"time"

	"github.com/capstone-orang-jago/parkin-be/app"
	"github.com/capstone-orang-jago/parkin-be/config"
	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/http/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// middleware "github.com/rumahkueku-azzahra-dev/terumpit-coklat/http/middlewares"

	_ "github.com/capstone-orang-jago/parkin-be/docs"
)

// HealthCheck godoc
// @Summary Check server status
// @Description get server health
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func health(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func RegisterRoutes(application app.App) http.Handler {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     config.AppConfig.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "application/json"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Serve static files from uploads directory
	// router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")

	// PUBLIC ROUTES
	v1.GET("/health", health)

	auth := v1.Group("/auth")
	{
		auth.POST("/register", application.AuthHandler.Register)
		auth.POST("/verify-otp", application.AuthHandler.VerifyRegister)
		auth.POST("/send-otp", application.OtpCodeHandler.ResendOtpCode)

		auth.POST("/login", application.AuthHandler.Login)
		auth.POST("/refresh", application.AuthHandler.RefreshToken)
		auth.POST("/logout", middlewares.AuthMiddleware(), application.AuthHandler.Logout)
	}

	user := v1.Group("/users")
	user.Use(middlewares.AuthMiddleware()).Use(middlewares.RoleMiddleware(domain.RolePengguna))
	{
		user.GET("/", application.UserHandler.GetAllUsers)
		user.GET("/:id", application.UserHandler.GetUserByID)
		user.POST("/", application.UserHandler.CreateUser)
		user.PUT("/:id", application.UserHandler.UpdateUser)
		user.DELETE("/:id", application.UserHandler.DeleteUser)
		user.GET("/me/saldo", application.UserHandler.GetUserSaldo)
		user.PUT("/me/saldo", application.UserHandler.DeductUserSaldo)
		user.GET("/profile", application.UserHandler.GetUserProfile)
		user.GET("/transactions", application.TransaksiParkirHandler.GetAllTransaksiParkirByPenggunaID)
	}

	lokasiParkir := v1.Group("/lokasi-parkir")
	lokasiParkir.Use(middlewares.AuthMiddleware())
	{
		lokasiParkir.GET("/", application.LokasiParkirHandler.GetAllLokasiParkir)
		lokasiParkir.GET("/:id", application.LokasiParkirHandler.GetLokasiParkirByID)
		lokasiParkir.GET("/petugas/:petugas_id", application.LokasiParkirHandler.GetLokasiParkirByPetugasID)
		lokasiParkir.POST("/", application.LokasiParkirHandler.CreateLokasiParkir)
		lokasiParkir.PUT("/:id", application.LokasiParkirHandler.UpdateLokasiParkir)
		lokasiParkir.DELETE("/:id", application.LokasiParkirHandler.DeleteLokasiParkir)
	}

	officer := v1.Group("/officer")
	officer.Use(middlewares.AuthMiddleware())
	{
		officer.POST("/qr-code", application.QRSessionHandler.CreateQRSession)
	}

	transaksiParkir := v1.Group("/transaksi-parkir")
	transaksiParkir.Use(middlewares.AuthMiddleware())
	{
		transaksiParkir.POST("/", application.TransaksiParkirHandler.CreateTransaksiParkir)
		transaksiParkir.GET("/:id", application.TransaksiParkirHandler.GetTransaksiParkirByID)
	}
	return router
}
