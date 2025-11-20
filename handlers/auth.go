package handlers

import (
	"context"
	"mbkm-api/config"
	"mbkm-api/database"
	"mbkm-api/models"
	"mbkm-api/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	db  *database.Database
	cfg *config.Config
}

func NewAuthHandler(db *database.Database, cfg *config.Config) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Register Request"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or user already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return utils.BadRequestResponse(c, "Username, email, and password are required")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to hash password")
	}

	ctx := context.Background()
	var userID int
	query := `
		INSERT INTO "user" (username, email, password_hash, full_name, phone, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err = h.db.Pool.QueryRow(ctx, query, req.Username, req.Email, hashedPassword, req.FullName, req.Phone, req.Role).Scan(&userID)
	if err != nil {
		return utils.ConflictResponse(c, "Username or email already exists")
	}

	token, err := utils.GenerateToken(userID, req.Email, req.Role, h.cfg.JWTSecret, h.cfg.JWTExpiration)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to generate token")
	}

	return utils.CreatedResponse(c, "User registered successfully", fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":       userID,
			"username": req.Username,
			"email":    req.Email,
			"role":     req.Role,
		},
	})
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login Request"
// @Success 200 {object} models.LoginResponse "Login successful"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Failure 403 {object} map[string]interface{} "Account inactive"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	ctx := context.Background()
	var user models.User
	query := `SELECT id, username, email, password_hash, full_name, role, is_active FROM "user" WHERE email = $1`
	err := h.db.Pool.QueryRow(ctx, query, req.Email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.IsActive,
	)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid email or password")
	}

	if !user.IsActive {
		return utils.ForbiddenResponse(c, "Account is inactive")
	}

	if err := utils.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return utils.UnauthorizedResponse(c, "Invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, h.cfg.JWTSecret, h.cfg.JWTExpiration)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to generate token")
	}

	user.PasswordHash = ""
	return utils.SuccessResponse(c, "Login successful", models.LoginResponse{
		Token: token,
		User:  user,
	})
}

// GetMe godoc
// @Summary Get current user profile
// @Description Get authenticated user's profile information
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User "User profile retrieved"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	ctx := context.Background()
	var user models.User
	query := `SELECT id, username, email, full_name, phone, role, is_active, created_at, updated_at FROM "user" WHERE id = $1`
	err := h.db.Pool.QueryRow(ctx, query, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.FullName, &user.Phone, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return utils.NotFoundResponse(c, "User not found")
	}

	return utils.SuccessResponse(c, "User profile retrieved", user)
}
