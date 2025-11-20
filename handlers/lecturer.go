package handlers

import (
	"context"
	"mbkm-api/database"
	"mbkm-api/models"
	"mbkm-api/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type LecturerHandler struct {
	db *database.Database
}

func NewLecturerHandler(db *database.Database) *LecturerHandler {
	return &LecturerHandler{db: db}
}

// GetAll godoc
// @Summary Get all lecturers
// @Description Retrieve list of all lecturers
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Lecturer "Lecturers retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /lecturers [get]
func (h *LecturerHandler) GetAll(c *fiber.Ctx) error {
	ctx := context.Background()
	query := `SELECT id, user_id, nidn, full_name, phone, department, is_active, created_at, updated_at FROM "lecturer" ORDER BY full_name ASC`

	rows, err := h.db.Pool.Query(ctx, query)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch lecturers")
	}
	defer rows.Close()

	var lecturers []models.Lecturer
	for rows.Next() {
		var l models.Lecturer
		var id, userID int64
		err := rows.Scan(&id, &userID, &l.NIDN, &l.FullName, &l.Phone, &l.Department, &l.IsActive, &l.CreatedAt, &l.UpdatedAt)
		if err != nil {
			return utils.InternalServerErrorResponse(c, "Failed to scan lecturer data: "+err.Error())
		}
		l.ID = int(id)
		l.UserID = int(userID)
		lecturers = append(lecturers, l)
	}

	if lecturers == nil {
		lecturers = []models.Lecturer{}
	}

	return utils.SuccessResponse(c, "Lecturers retrieved successfully", lecturers)
}

// GetByID godoc
// @Summary Get lecturer by ID
// @Description Retrieve a specific lecturer by ID
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Lecturer ID"
// @Success 200 {object} models.Lecturer "Lecturer retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid lecturer ID"
// @Failure 404 {object} map[string]interface{} "Lecturer not found"
// @Router /lecturers/{id} [get]
func (h *LecturerHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid lecturer ID")
	}

	ctx := context.Background()
	var lecturer models.Lecturer
	query := `SELECT id, user_id, nidn, full_name, phone, department, is_active, created_at, updated_at FROM "lecturer" WHERE id = $1`

	var lid, userID int64
	err = h.db.Pool.QueryRow(ctx, query, id).Scan(&lid, &userID, &lecturer.NIDN, &lecturer.FullName, &lecturer.Phone, &lecturer.Department, &lecturer.IsActive, &lecturer.CreatedAt, &lecturer.UpdatedAt)
	if err != nil {
		return utils.NotFoundResponse(c, "Lecturer not found")
	}
	lecturer.ID = int(lid)
	lecturer.UserID = int(userID)

	return utils.SuccessResponse(c, "Lecturer retrieved successfully", lecturer)
}

// Create godoc
// @Summary Create new lecturer
// @Description Create a new lecturer (admin only)
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateLecturerRequest true "Lecturer details"
// @Success 201 {object} map[string]interface{} "Lecturer created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 409 {object} map[string]interface{} "NIDN already exists"
// @Router /lecturers [post]
func (h *LecturerHandler) Create(c *fiber.Ctx) error {
	var req models.CreateLecturerRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate required fields
	if req.NIDN == "" || req.FullName == "" {
		return utils.BadRequestResponse(c, "NIDN and full name are required")
	}

	ctx := context.Background()

	// Check if user exists and is a lecturer
	var userExists bool
	err := h.db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM "user" WHERE id = $1 AND role = 'lecturer')`, req.UserID).Scan(&userExists)
	if err != nil || !userExists {
		return utils.BadRequestResponse(c, "Invalid user ID or user is not a lecturer")
	}

	var lecturerID int64
	query := `INSERT INTO "lecturer" (user_id, nidn, full_name, phone, department, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id`

	err = h.db.Pool.QueryRow(ctx, query, req.UserID, req.NIDN, req.FullName, req.Phone, req.Department).Scan(&lecturerID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return utils.ConflictResponse(c, "NIDN or User ID already exists")
		}
		return utils.InternalServerErrorResponse(c, "Failed to create lecturer")
	}

	return utils.CreatedResponse(c, "Lecturer created successfully", fiber.Map{"id": int(lecturerID)})
}

// Update godoc
// @Summary Update lecturer
// @Description Update an existing lecturer (admin only)
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Lecturer ID"
// @Param request body models.UpdateLecturerRequest true "Updated lecturer details"
// @Success 200 {object} map[string]interface{} "Lecturer updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Lecturer not found"
// @Failure 409 {object} map[string]interface{} "NIDN already exists"
// @Router /lecturers/{id} [put]
func (h *LecturerHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid lecturer ID")
	}

	var req models.UpdateLecturerRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	ctx := context.Background()

	// Check if lecturer exists
	var exists bool
	err = h.db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM "lecturer" WHERE id = $1)`, id).Scan(&exists)
	if err != nil || !exists {
		return utils.NotFoundResponse(c, "Lecturer not found")
	}

	query := `UPDATE "lecturer" SET nidn = $1, full_name = $2, phone = $3, department = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5`

	result, err := h.db.Pool.Exec(ctx, query, req.NIDN, req.FullName, req.Phone, req.Department, id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return utils.ConflictResponse(c, "NIDN already exists")
		}
		return utils.InternalServerErrorResponse(c, "Failed to update lecturer")
	}

	if result.RowsAffected() == 0 {
		return utils.NotFoundResponse(c, "Lecturer not found")
	}

	return utils.SuccessResponse(c, "Lecturer updated successfully", nil)
}

// Delete godoc
// @Summary Delete lecturer
// @Description Delete a lecturer (admin only)
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Lecturer ID"
// @Success 200 {object} map[string]interface{} "Lecturer deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid lecturer ID"
// @Failure 404 {object} map[string]interface{} "Lecturer not found"
// @Failure 409 {object} map[string]interface{} "Cannot delete, lecturer has related data"
// @Router /lecturers/{id} [delete]
func (h *LecturerHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid lecturer ID")
	}

	ctx := context.Background()

	// Check if lecturer exists
	var exists bool
	err = h.db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM "lecturer" WHERE id = $1)`, id).Scan(&exists)
	if err != nil || !exists {
		return utils.NotFoundResponse(c, "Lecturer not found")
	}

	query := `DELETE FROM "lecturer" WHERE id = $1`

	result, err := h.db.Pool.Exec(ctx, query, id)
	if err != nil {
		if strings.Contains(err.Error(), "foreign key") {
			return utils.ConflictResponse(c, "Cannot delete lecturer, has related programs or data")
		}
		return utils.InternalServerErrorResponse(c, "Failed to delete lecturer")
	}

	if result.RowsAffected() == 0 {
		return utils.NotFoundResponse(c, "Lecturer not found")
	}

	return utils.SuccessResponse(c, "Lecturer deleted successfully", nil)
}
