package handlers

import (
	"context"
	"mbkm-api/database"
	"mbkm-api/models"
	"mbkm-api/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type EnrollmentHandler struct {
	db *database.Database
}

func NewEnrollmentHandler(db *database.Database) *EnrollmentHandler {
	return &EnrollmentHandler{db: db}
}

// GetAll godoc
// @Summary Get all enrollments
// @Description Retrieve list of all enrollments (admin/lecturer only)
// @Tags Enrollments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Enrollment "Enrollments retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /enrollments [get]
func (h *EnrollmentHandler) GetAll(c *fiber.Ctx) error {
	ctx := context.Background()
	query := `SELECT id, student_id, program_id, status, enrolled_at, created_at, updated_at FROM "enrollment" ORDER BY created_at DESC`

	rows, err := h.db.Pool.Query(ctx, query)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch enrollments")
	}
	defer rows.Close()

	var enrollments []models.Enrollment
	for rows.Next() {
		var e models.Enrollment
		var id, studentID, programID int64
		err := rows.Scan(&id, &studentID, &programID, &e.Status, &e.EnrolledAt, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return utils.InternalServerErrorResponse(c, "Failed to scan enrollment data: "+err.Error())
		}
		e.ID = int(id)
		e.StudentID = int(studentID)
		e.ProgramID = int(programID)
		enrollments = append(enrollments, e)
	}

	if enrollments == nil {
		enrollments = []models.Enrollment{}
	}

	return utils.SuccessResponse(c, "Enrollments retrieved successfully", enrollments)
}

// GetByStudent godoc
// @Summary Get enrollments by student
// @Description Retrieve all enrollments for a specific student
// @Tags Enrollments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param studentId path int true "Student ID"
// @Success 200 {array} models.Enrollment "Student enrollments retrieved"
// @Failure 400 {object} map[string]interface{} "Invalid student ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /enrollments/student/{studentId} [get]
func (h *EnrollmentHandler) GetByStudent(c *fiber.Ctx) error {
	studentID, err := strconv.Atoi(c.Params("studentId"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid student ID")
	}

	ctx := context.Background()
	query := `SELECT id, student_id, program_id, status, enrolled_at, created_at, updated_at FROM "enrollment" WHERE student_id = $1 ORDER BY created_at DESC`

	rows, err := h.db.Pool.Query(ctx, query, studentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch enrollments")
	}
	defer rows.Close()

	var enrollments []models.Enrollment
	for rows.Next() {
		var e models.Enrollment
		var id, studentID, programID int64
		err := rows.Scan(&id, &studentID, &programID, &e.Status, &e.EnrolledAt, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return utils.InternalServerErrorResponse(c, "Failed to scan enrollment data")
		}
		e.ID = int(id)
		e.StudentID = int(studentID)
		e.ProgramID = int(programID)
		enrollments = append(enrollments, e)
	}

	if enrollments == nil {
		enrollments = []models.Enrollment{}
	}

	return utils.SuccessResponse(c, "Student enrollments retrieved successfully", enrollments)
}

// Create godoc
// @Summary Create new enrollment
// @Description Enroll a student in a program
// @Tags Enrollments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateEnrollmentRequest true "Enrollment details"
// @Success 201 {object} map[string]interface{} "Enrollment created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or already enrolled"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /enrollments [post]
func (h *EnrollmentHandler) Create(c *fiber.Ctx) error {
	var req models.CreateEnrollmentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	ctx := context.Background()
	var enrollmentID int
	query := `INSERT INTO "enrollment" (student_id, program_id) VALUES ($1, $2) RETURNING id`

	err := h.db.Pool.QueryRow(ctx, query, req.StudentID, req.ProgramID).Scan(&enrollmentID)
	if err != nil {
		return utils.ConflictResponse(c, "Failed to create enrollment or already enrolled")
	}

	return utils.CreatedResponse(c, "Enrollment created successfully", fiber.Map{"id": enrollmentID})
}

func (h *EnrollmentHandler) UpdateStatus(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid enrollment ID")
	}

	var req models.UpdateEnrollmentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	ctx := context.Background()
	query := `UPDATE "enrollment" SET status = $1 WHERE id = $2`

	result, err := h.db.Pool.Exec(ctx, query, req.Status, id)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to update enrollment status")
	}

	if result.RowsAffected() == 0 {
		return utils.NotFoundResponse(c, "Enrollment not found")
	}

	return utils.SuccessResponse(c, "Enrollment status updated successfully", nil)
}

func (h *EnrollmentHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid enrollment ID")
	}

	ctx := context.Background()
	query := `DELETE FROM "enrollment" WHERE id = $1`

	result, err := h.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to delete enrollment")
	}

	if result.RowsAffected() == 0 {
		return utils.NotFoundResponse(c, "Enrollment not found")
	}

	return utils.SuccessResponse(c, "Enrollment deleted successfully", nil)
}
