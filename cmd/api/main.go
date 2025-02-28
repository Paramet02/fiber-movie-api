package main

import (
	"fmt"
	// "log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	// "github.com/joho/godotenv"
	"github.com/paramet02/webapi/auth"
	"github.com/paramet02/webapi/config"
	"github.com/paramet02/webapi/handlers"
	"github.com/paramet02/webapi/repository"
	"github.com/paramet02/webapi/services"
)

type Application struct {
	auth         auth.Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
}


func main() {
	var application Application

	// Parse the command line arguments for JWT
	application.JWTSecret = os.Getenv("JWT_SECRET")
	application.JWTIssuer = os.Getenv("JWT_ISSUER")
	application.JWTAudience = os.Getenv("JWT_ADDIENCE")
	application.CookieDomain = os.Getenv("JWT_COOKIEDOMAIN")

	
	db := config.InitDatabase()
	defer config.CloseDatabase(db)
	
	application.auth = auth.Auth{
		JWTIssuer:        application.JWTIssuer,
		JWTAudience:      application.JWTAudience,
		JWTSecret:        application.JWTSecret,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "__Host-refresh_token",
		CookieDomain:  application.CookieDomain,
	}


	movieRepo := repository.NewMovieRepositoryDB(db)
	movieSer := services.NewMovieService(movieRepo)
	movieHandler := handlers.NewmovieHandler(movieSer)

	userRepo := repository.NewuserRepositoryDB(db)
	userSer := services.NewuserService(userRepo , &application.auth)
	userHandler := handlers.NewuserHandler(userSer , &application.auth)

	// สร้าง Fiber instance
	app := fiber.New()
	
	// ใช้ Middleware Recover
	app.Use(recover.New())

	// ใช้ Middleware CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173", // ระบุ Origin ที่อนุญาต
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS", // กำหนด Methods ที่อนุญาต
		AllowHeaders:     "Accept, Content-Type, X-CSRF-Token, Authorization", // กำหนด Headers ที่อนุญาต
		AllowCredentials: true, // อนุญาต Credentials
		MaxAge:           3600, // ระยะเวลาในวินาทีที่เบราว์เซอร์จะเก็บผลลัพธ์ของ CORS preflight cache
	}))

	
	// Routes that do not require authentication
	app.Post("/Login" , userHandler.Login)
	app.Post("/Register" , userHandler.Register)
	app.Get("/RefreshToken" , userHandler.RefreshToken)
	app.Get("/Logout" , userHandler.Logout)

	app.Get("/movie", movieHandler.GetsMovies)
	app.Get("/movies/:id" , movieHandler.GetMovie)
	app.Get("/genres" , movieHandler.GetsGenres)
	app.Get("/User/Email" , userHandler.GetUserByEmail)
	app.Get("/User/:id" , userHandler.GetUserByID)

	// Admin routes that require authentication
	adminGroup := app.Group("/admin")
	adminGroup.Use(application.authRequired())  // This ensures that the routes in this group require authentication
	adminGroup.Use(application.jwtMiddleware)
	adminGroup.Get("/movies/:id" , movieHandler.OneMovieForEdit)
	adminGroup.Post("/Insert" , movieHandler.InsertMovie)
	adminGroup.Put("/movie/:id" , movieHandler.UpdateMovie)
	adminGroup.Delete("/movie/:id" , movieHandler.DeleteMovie)

 
	// เปิด server ด้วย port	
	app.Listen(fmt.Sprintf(":%s" , os.Getenv("GO_PORT")))
	// เช็ค error ถ้ามี err ก็ down server
}

