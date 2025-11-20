package main

import (
	"log"
	"mbkm-api/config"
	"mbkm-api/database"
	_ "mbkm-api/docs"
	"mbkm-api/models"
	"mbkm-api/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// @title MBKM API
// @version 1.0
// @description API for MBKM (Merdeka Belajar Kampus Merdeka) Program Management System
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@mbkm-api.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("❌ Failed to load config:", err)
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer db.Close()

	// Auto-migrate menggunakan GORM
	if err := db.AutoMigrate(
		&models.User{},
		&models.Lecturer{},
		&models.Program{},
		&models.Enrollment{},
		&models.Assessment{},
	); err != nil {
		log.Fatal("Auto-migration failed:", err)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			log.Println("Migration complete, exiting...")
			return
		case "seed":
			log.Println("🌱 Running seeder...")
			seeder := database.NewSeeder(db)
			if err := seeder.SeedAll(); err != nil {
				log.Fatal("Seeding failed:", err)
			}
			log.Println("Seeding complete, exiting...")
			return
		case "seed:users":
			log.Println("🌱 Running user seeder...")
			seeder := database.NewSeeder(db)
			if err := seeder.SeedUsers(); err != nil {
				log.Fatal("User seeding failed:", err)
			}
			log.Println("User seeding complete, exiting...")
			return
		case "seed:programs":
			log.Println("🌱 Running program seeder...")
			seeder := database.NewSeeder(db)
			if err := seeder.SeedPrograms(); err != nil {
				log.Fatal("Program seeding failed:", err)
			}
			log.Println("Program seeding complete, exiting...")
			return
		case "seed:lecturers":
			log.Println("🌱 Running lecturer seeder...")
			seeder := database.NewSeeder(db)
			if err := seeder.SeedLecturers(); err != nil {
				log.Fatal("Lecturer seeding failed:", err)
			}
			log.Println("Lecturer seeding complete, exiting...")
			return
		}
	}

	app := fiber.New(fiber.Config{
		AppName: "MBKM API v1.0",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	routes.SetupRoutes(app, db, cfg)

	log.Printf("Server running on port %s\n", cfg.ServerPort)
	log.Printf("Health check: http://localhost:%s/health\n", cfg.ServerPort)

	if err := app.Listen(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
