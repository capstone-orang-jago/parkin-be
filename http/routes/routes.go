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
	v1.POST("/midtrans/callback", application.MidtransHandler.MidtransCallback)

	auth := v1.Group("/auth")
	{
		auth.POST("/register", application.AuthHandler.Register)
		auth.POST("/verify-otp", application.AuthHandler.VerifyRegister)
		auth.POST("/send-otp", application.OtpCodeHandler.ResendOtpCode)

		auth.POST("/login", application.AuthHandler.Login)
		auth.POST("/refresh", application.AuthHandler.RefreshToken)
		auth.POST("/logout", middlewares.AuthMiddleware(), application.AuthHandler.Logout)
	}

	user := v1.Group("/user")
	user.Use(middlewares.AuthMiddleware()).Use(middlewares.RoleMiddleware(domain.RolePengguna))
	{
		user.GET("/profile", application.UserHandler.GetUserProfile)
		user.GET("/transactions", application.TransaksiParkirHandler.GetAllTransaksiParkirByUserID)
		user.GET("/transactions/:id", application.TransaksiParkirHandler.GetTransaksiParkirByID)
		user.POST("/parking/pay", application.TransaksiParkirHandler.CreateTransaksiParkir)
		user.POST("/topup/create", application.MidtransHandler.CreateTopUp)
		user.GET("/topup/:order_id/status", application.MidtransHandler.GetTopUpStatus)

	}

	officer := v1.Group("/officer")
	officer.Use(middlewares.AuthMiddleware()).Use(middlewares.RoleMiddleware(domain.RolePetugas))
	{
		officer.GET("/profile", application.UserHandler.GetUserProfile)
		officer.GET("/dashboard", application.DashboardHandler.GetOfficerDashboard)
		officer.POST("/qr-code", application.QRSessionHandler.CreateQRSession)
		officer.GET("/transactions", application.TransaksiParkirHandler.GetAllTransaksiParkirByUserID)
	}

	admin := v1.Group("/admin")
	admin.Use(middlewares.AuthMiddleware()).Use(middlewares.RoleMiddleware("admin"))
	{
		admin.GET("/", application.UserHandler.GetAllUsers)
		admin.GET("/:id", application.UserHandler.GetUserByID)
		admin.POST("/", application.UserHandler.CreateUser)
		admin.PUT("/:id", application.UserHandler.UpdateUser)
		admin.DELETE("/:id", application.UserHandler.DeleteUser)
		admin.GET("/me/saldo", application.UserHandler.GetUserSaldo)
		admin.PUT("/me/saldo", application.UserHandler.DeductUserSaldo)
	}

	transaksiParkir := v1.Group("/transaksi-parkir")
	transaksiParkir.Use(middlewares.AuthMiddleware())
	{
		transaksiParkir.GET("/:id", application.TransaksiParkirHandler.GetTransaksiParkirByID)
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

	return router
}
