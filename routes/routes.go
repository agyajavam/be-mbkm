package routes

import (
	"mbkm-api/config"
	"mbkm-api/database"
	"mbkm-api/handlers"
	"mbkm-api/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App, db *database.Database, cfg *config.Config) {
	authHandler := handlers.NewAuthHandler(db, cfg)
	programHandler := handlers.NewProgramHandler(db)
	enrollmentHandler := handlers.NewEnrollmentHandler(db)
	assessmentHandler := handlers.NewAssessmentHandler(db)
	lecturerHandler := handlers.NewLecturerHandler(db)

	api := app.Group("/api/v1")

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "MBKM API is running",
		})
	})

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	protected := api.Use(middleware.AuthMiddleware(cfg))

	protected.Get("/auth/me", authHandler.GetMe)

	programs := protected.Group("/programs")
	programs.Get("/", programHandler.GetAll)
	programs.Get("/:id", programHandler.GetByID)
	programs.Post("/", middleware.RoleMiddleware("admin", "lecturer"), programHandler.Create)
	programs.Put("/:id", middleware.RoleMiddleware("admin", "lecturer"), programHandler.Update)
	programs.Delete("/:id", middleware.RoleMiddleware("admin"), programHandler.Delete)

	lecturers := protected.Group("/lecturers")
	lecturers.Get("/", lecturerHandler.GetAll)
	lecturers.Get("/:id", lecturerHandler.GetByID)
	lecturers.Post("/", middleware.RoleMiddleware("admin"), lecturerHandler.Create)
	lecturers.Put("/:id", middleware.RoleMiddleware("admin"), lecturerHandler.Update)
	lecturers.Delete("/:id", middleware.RoleMiddleware("admin"), lecturerHandler.Delete)

	enrollments := protected.Group("/enrollments")
	enrollments.Get("/", middleware.RoleMiddleware("admin", "lecturer"), enrollmentHandler.GetAll)
	enrollments.Get("/student/:studentId", enrollmentHandler.GetByStudent)
	enrollments.Post("/", middleware.RoleMiddleware("admin", "lecturer", "student"), enrollmentHandler.Create)
	enrollments.Put("/:id/status", middleware.RoleMiddleware("admin", "lecturer"), enrollmentHandler.UpdateStatus)
	enrollments.Delete("/:id", middleware.RoleMiddleware("admin"), enrollmentHandler.Delete)

	assessments := protected.Group("/assessments")
	assessments.Get("/enrollment/:enrollmentId", assessmentHandler.GetByEnrollment)
	assessments.Post("/", middleware.RoleMiddleware("admin", "lecturer"), assessmentHandler.Create)
	assessments.Put("/:id", middleware.RoleMiddleware("admin", "lecturer"), assessmentHandler.Update)
	assessments.Delete("/:id", middleware.RoleMiddleware("admin", "lecturer"), assessmentHandler.Delete)
}
