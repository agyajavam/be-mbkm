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

type ProgramHandler struct {
	db *database.Database
}

func NewProgramHandler(db *database.Database) *ProgramHandler {
	return &ProgramHandler{db: db}
}

// GetAll godoc
// @Summary Get all programs
// @Description Retrieve list of all MBKM programs
// @Tags Programs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Program "Programs retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /programs [get]
func (h *ProgramHandler) GetAll(c *fiber.Ctx) error {
	ctx := context.Background()
	query := `SELECT id, code, name, description, credits, semester, lecturer_id, is_active, created_at, updated_at FROM "program" ORDER BY created_at DESC`

	rows, err := h.db.Pool.Query(ctx, query)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch programs")
	}
	defer rows.Close()

	var programs []models.Program
	for rows.Next() {
		var p models.Program
		var id, credits, semester, lecturerID int64
		err := rows.Scan(&id, &p.Code, &p.Name, &p.Description, &credits, &semester, &lecturerID, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return utils.InternalServerErrorResponse(c, "Failed to scan program data: "+err.Error())
		}
		p.ID = int(id)
		p.Credits = int(credits)
		p.Semester = int(semester)
		p.LecturerID = int(lecturerID)
		programs = append(programs, p)
	}

	if programs == nil {
		programs = []models.Program{}
	}

	return utils.SuccessResponse(c, "Programs retrieved successfully", programs)
}

// GetByID godoc
// @Summary Get program by ID
// @Description Retrieve a specific program by its ID
// @Tags Programs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Program ID"
// @Success 200 {object} models.Program "Program retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid program ID"
// @Failure 404 {object} map[string]interface{} "Program not found"
// @Router /programs/{id} [get]
func (h *ProgramHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid program ID")
	}

	ctx := context.Background()
	var program models.Program
	query := `SELECT id, code, name, description, credits, semester, lecturer_id, is_active, created_at, updated_at FROM "program" WHERE id = $1`

	var pid, credits, semester, lecturerID int64
	err = h.db.Pool.QueryRow(ctx, query, id).Scan(&pid, &program.Code, &program.Name, &program.Description, &credits, &semester, &lecturerID, &program.IsActive, &program.CreatedAt, &program.UpdatedAt)
	if err != nil {
		return utils.NotFoundResponse(c, "Program not found")
	}
	program.ID = int(pid)
	program.Credits = int(credits)
	program.Semester = int(semester)
	program.LecturerID = int(lecturerID)

	return utils.SuccessResponse(c, "Program retrieved successfully", program)
}

// Create godoc
// @Summary Create new program
// @Description Create a new MBKM program (admin/lecturer only)
// @Tags Programs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateProgramRequest true "Program details"
// @Success 201 {object} map[string]interface{} "Program created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /programs [post]
func (h *ProgramHandler) Create(c *fiber.Ctx) error {
	var req models.CreateProgramRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate required fields
	if req.Code == "" || req.Name == "" {
		return utils.BadRequestResponse(c, "Code and name are required")
	}

	ctx := context.Background()

	// Check if lecturer exists
	var lecturerExists bool
	err := h.db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM "lecturer" WHERE id = $1 AND is_active = true)`, req.LecturerID).Scan(&lecturerExists)
	if err != nil || !lecturerExists {
		return utils.BadRequestResponse(c, "Invalid lecturer ID")
	}

	var programID int64
	query := `INSERT INTO "program" (code, name, description, credits, semester, lecturer_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id`

	err = h.db.Pool.QueryRow(ctx, query, req.Code, req.Name, req.Description, req.Credits, req.Semester, req.LecturerID).Scan(&programID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return utils.ConflictResponse(c, "Program code already exists")
		}
		return utils.InternalServerErrorResponse(c, "Failed to create program")
	}

	return utils.CreatedResponse(c, "Program created successfully", fiber.Map{"id": int(programID)})
}

// Update godoc
// @Summary Update program
// @Description Update an existing program (admin/lecturer only)
// @Tags Programs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Program ID"
// @Param request body models.UpdateProgramRequest true "Updated program details"
// @Success 200 {object} map[string]interface{} "Program updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /programs/{id} [put]
func (h *ProgramHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid program ID")
	}

	var req models.UpdateProgramRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	ctx := context.Background()

	// Check if program exists
	var exists bool
	err = h.db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM "program" WHERE id = $1)`, id).Scan(&exists)
	if err != nil || !exists {
		return utils.NotFoundResponse(c, "Program not found")
	}

	query := `UPDATE "program" SET code = $1, name = $2, description = $3, credits = $4, semester = $5, updated_at = CURRENT_TIMESTAMP WHERE id = $6`

	result, err := h.db.Pool.Exec(ctx, query, req.Code, req.Name, req.Description, req.Credits, req.Semester, id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return utils.ConflictResponse(c, "Program code already exists")
		}
		return utils.InternalServerErrorResponse(c, "Failed to update program")
	}

	if result.RowsAffected() == 0 {
		return utils.NotFoundResponse(c, "Program not found")
	}

	return utils.SuccessResponse(c, "Program updated successfully", nil)
}

// Delete godoc
// @Summary Delete program
// @Description Delete a program (admin only)
// @Tags Programs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Program ID"
// @Success 200 {object} map[string]interface{} "Program deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid program ID"
// @Failure 404 {object} map[string]interface{} "Program not found"
// @Failure 409 {object} map[string]interface{} "Cannot delete, program has related data"
// @Router /programs/{id} [delete]
func (h *ProgramHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid program ID")
	}

	ctx := context.Background()

	// Check if program exists
	var exists bool
	err = h.db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM "program" WHERE id = $1)`, id).Scan(&exists)
	if err != nil || !exists {
		return utils.NotFoundResponse(c, "Program not found")
	}

	query := `DELETE FROM "program" WHERE id = $1`

	result, err := h.db.Pool.Exec(ctx, query, id)
	if err != nil {
		if strings.Contains(err.Error(), "foreign key") {
			return utils.ConflictResponse(c, "Cannot delete program, it has related enrollments")
		}
		return utils.InternalServerErrorResponse(c, "Failed to delete program")
	}

	if result.RowsAffected() == 0 {
		return utils.NotFoundResponse(c, "Program not found")
	}

	return utils.SuccessResponse(c, "Program deleted successfully", nil)
}
