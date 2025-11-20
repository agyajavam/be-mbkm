package handlers

import (
	"context"
	"mbkm-api/database"
	"mbkm-api/models"
	"mbkm-api/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AssessmentHandler struct {
	db *database.Database
}

func NewAssessmentHandler(db *database.Database) *AssessmentHandler {
	return &AssessmentHandler{db: db}
}

// GetByEnrollment godoc
// @Summary Get assessments by enrollment
// @Description Retrieve all assessments for a specific enrollment
// @Tags Assessments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param enrollmentId path int true "Enrollment ID"
// @Success 200 {array} models.Assessment "Assessments retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid enrollment ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /assessments/enrollment/{enrollmentId} [get]
func (h *AssessmentHandler) GetByEnrollment(c *fiber.Ctx) error {
	enrollmentID, err := strconv.Atoi(c.Params("enrollmentId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid enrollment ID")
	}

	ctx := context.Background()
	query := `SELECT id, enrollment_id, student_id, program_id, category, score, max_score, weight, notes, created_at, updated_at FROM "assessment" WHERE enrollment_id = $1 ORDER BY created_at DESC`

	rows, err := h.db.Pool.Query(ctx, query, enrollmentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch assessments")
	}
	defer rows.Close()

	var assessments []models.Assessment
	for rows.Next() {
		var a models.Assessment
		var id, enrollmentID, studentID, programID int64
		err := rows.Scan(&id, &enrollmentID, &studentID, &programID, &a.Category, &a.Score, &a.MaxScore, &a.Weight, &a.Notes, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return utils.InternalServerErrorResponse(c, "Failed to scan assessment data")
		}
		a.ID = int(id)
		a.EnrollmentID = int(enrollmentID)
		a.StudentID = int(studentID)
		a.ProgramID = int(programID)
		assessments = append(assessments, a)
	}

	if assessments == nil {
		assessments = []models.Assessment{}
	}

	return utils.SuccessResponse(c, "Assessments retrieved successfully", assessments)
}

func (h *AssessmentHandler) Create(c *fiber.Ctx) error {
	var req models.CreateAssessmentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	ctx := context.Background()
	var sid, pid int64
	enrollmentQuery := `SELECT student_id, program_id FROM "enrollment" WHERE id = $1`
	err := h.db.Pool.QueryRow(ctx, enrollmentQuery, req.EnrollmentID).Scan(&sid, &pid)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid enrollment ID")
	}
	studentID := int(sid)
	programID := int(pid)

	var assessmentID int
	insertQuery := `INSERT INTO "assessment" (enrollment_id, student_id, program_id, category, score, max_score, weight, notes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err = h.db.Pool.QueryRow(ctx, insertQuery, req.EnrollmentID, studentID, programID, req.Category, req.Score, req.MaxScore, req.Weight, req.Notes).Scan(&assessmentID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create assessment")
	}

	return utils.CreatedResponse(c, "Assessment created successfully", fiber.Map{"id": assessmentID})
}

func (h *AssessmentHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid assessment ID")
	}

	var req models.UpdateAssessmentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	ctx := context.Background()
	query := `UPDATE "assessment" SET score = $1, max_score = $2, weight = $3, notes = $4 WHERE id = $5`

	result, err := h.db.Pool.Exec(ctx, query, req.Score, req.MaxScore, req.Weight, req.Notes, id)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to update assessment")
	}

	if result.RowsAffected() == 0 {
		return utils.NotFoundResponse(c, "Assessment not found")
	}

	return utils.SuccessResponse(c, "Assessment updated successfully", nil)
}

func (h *AssessmentHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid assessment ID")
	}

	ctx := context.Background()
	query := `DELETE FROM "assessment" WHERE id = $1`

	result, err := h.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to delete assessment")
	}

	if result.RowsAffected() == 0 {
		return utils.NotFoundResponse(c, "Assessment not found")
	}

	return utils.SuccessResponse(c, "Assessment deleted successfully", nil)
}
