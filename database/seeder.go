package database

import (
	"context"
	"log"
	"mbkm-api/utils"
)

type Seeder struct {
	db *Database
}

func NewSeeder(db *Database) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) SeedUsers() error {
	ctx := context.Background()

	users := []struct {
		Username string
		Email    string
		Password string
		FullName string
		Phone    string
		Role     string
	}{
		{
			Username: "admin",
			Email:    "admin@mbkm.ac.id",
			Password: "admin123",
			FullName: "Administrator",
			Phone:    "081234567890",
			Role:     "admin",
		},
		{
			Username: "lecturer1",
			Email:    "lecturer1@mbkm.ac.id",
			Password: "lecturer123",
			FullName: "Dr. Budi Santoso",
			Phone:    "081234567891",
			Role:     "lecturer",
		},
		{
			Username: "lecturer2",
			Email:    "lecturer2@mbkm.ac.id",
			Password: "lecturer123",
			FullName: "Dr. Siti Nurhaliza",
			Phone:    "081234567892",
			Role:     "lecturer",
		},
		{
			Username: "student1",
			Email:    "student1@mbkm.ac.id",
			Password: "student123",
			FullName: "Ahmad Fauzi",
			Phone:    "081234567893",
			Role:     "student",
		},
		{
			Username: "student2",
			Email:    "student2@mbkm.ac.id",
			Password: "student123",
			FullName: "Siti Rahmawati",
			Phone:    "081234567894",
			Role:     "student",
		},
		{
			Username: "student3",
			Email:    "student3@mbkm.ac.id",
			Password: "student123",
			FullName: "Andi Wijaya",
			Phone:    "081234567895",
			Role:     "student",
		},
	}

	log.Println("🌱 Seeding users...")

	for _, user := range users {
		// Check if user already exists
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM "user" WHERE email = $1)`
		err := s.db.Pool.QueryRow(ctx, checkQuery, user.Email).Scan(&exists)
		if err != nil {
			log.Printf("❌ Error checking user %s: %v", user.Email, err)
			continue
		}

		if exists {
			log.Printf("⏭️  User %s already exists, skipping...", user.Email)
			continue
		}

		// Hash password
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			log.Printf("❌ Error hashing password for %s: %v", user.Email, err)
			continue
		}

		// Insert user
		query := `
			INSERT INTO "user" (username, email, password_hash, full_name, phone, role, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`
		_, err = s.db.Pool.Exec(ctx, query, user.Username, user.Email, hashedPassword, user.FullName, user.Phone, user.Role)
		if err != nil {
			log.Printf("❌ Error inserting user %s: %v", user.Email, err)
			continue
		}

		log.Printf("✅ User created: %s (%s)", user.Email, user.Role)
	}

	log.Println("✅ User seeding completed!")
	return nil
}

func (s *Seeder) SeedPrograms() error {
	ctx := context.Background()

	// Get lecturer IDs from lecturer table
	var lecturer1ID, lecturer2ID int
	err := s.db.Pool.QueryRow(ctx, `SELECT id FROM "lecturer" WHERE nidn = $1`, "0123456789").Scan(&lecturer1ID)
	if err != nil {
		log.Println("❌ Lecturer1 not found, please seed lecturers first")
		return err
	}

	err = s.db.Pool.QueryRow(ctx, `SELECT id FROM "lecturer" WHERE nidn = $1`, "0987654321").Scan(&lecturer2ID)
	if err != nil {
		log.Println("❌ Lecturer2 not found, please seed lecturers first")
		return err
	}

	programs := []struct {
		Code        string
		Name        string
		Description string
		Credits     int
		Semester    int
		LecturerID  int
	}{
		{
			Code:        "MBKM001",
			Name:        "Studi Independen - Web Development",
			Description: "Program studi independen fokus pada pengembangan web modern menggunakan React, Node.js, dan PostgreSQL",
			Credits:     20,
			Semester:    5,
			LecturerID:  lecturer1ID,
		},
		{
			Code:        "MBKM002",
			Name:        "Magang Industri - Software Engineering",
			Description: "Program magang di perusahaan teknologi untuk pengalaman langsung software engineering",
			Credits:     20,
			Semester:    6,
			LecturerID:  lecturer1ID,
		},
		{
			Code:        "MBKM003",
			Name:        "Kampus Mengajar - Pendidikan Digital",
			Description: "Program mengajar di sekolah dengan fokus pada literasi digital dan teknologi",
			Credits:     20,
			Semester:    5,
			LecturerID:  lecturer2ID,
		},
		{
			Code:        "MBKM004",
			Name:        "Studi Independen - Data Science",
			Description: "Program pembelajaran data science, machine learning, dan analisis data",
			Credits:     20,
			Semester:    6,
			LecturerID:  lecturer2ID,
		},
		{
			Code:        "MBKM005",
			Name:        "Proyek Kemanusiaan - Tech for Good",
			Description: "Mengembangkan solusi teknologi untuk mengatasi masalah sosial",
			Credits:     20,
			Semester:    7,
			LecturerID:  lecturer1ID,
		},
	}

	log.Println("🌱 Seeding programs...")

	for _, program := range programs {
		// Check if program already exists
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM "program" WHERE code = $1)`
		err := s.db.Pool.QueryRow(ctx, checkQuery, program.Code).Scan(&exists)
		if err != nil {
			log.Printf("❌ Error checking program %s: %v", program.Code, err)
			continue
		}

		if exists {
			log.Printf("⏭️  Program %s already exists, skipping...", program.Code)
			continue
		}

		// Insert program
		query := `
			INSERT INTO "program" (code, name, description, credits, semester, lecturer_id, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`
		_, err = s.db.Pool.Exec(ctx, query, program.Code, program.Name, program.Description, program.Credits, program.Semester, program.LecturerID)
		if err != nil {
			log.Printf("❌ Error inserting program %s: %v", program.Code, err)
			continue
		}

		log.Printf("✅ Program created: %s - %s", program.Code, program.Name)
	}

	log.Println("✅ Program seeding completed!")
	return nil
}

func (s *Seeder) SeedLecturers() error {
	ctx := context.Background()

	// Get lecturer user IDs
	var lecturer1UserID, lecturer2UserID int
	err := s.db.Pool.QueryRow(ctx, `SELECT id FROM "user" WHERE email = $1`, "lecturer1@mbkm.ac.id").Scan(&lecturer1UserID)
	if err != nil {
		log.Println("❌ Lecturer1 user not found, please seed users first")
		return err
	}

	err = s.db.Pool.QueryRow(ctx, `SELECT id FROM "user" WHERE email = $1`, "lecturer2@mbkm.ac.id").Scan(&lecturer2UserID)
	if err != nil {
		log.Println("❌ Lecturer2 user not found, please seed users first")
		return err
	}

	lecturers := []struct {
		UserID     int
		NIDN       string
		FullName   string
		Phone      string
		Department string
	}{
		{
			UserID:     lecturer1UserID,
			NIDN:       "0123456789",
			FullName:   "Dr. Budi Santoso, M.Kom",
			Phone:      "081234567891",
			Department: "Informatika",
		},
		{
			UserID:     lecturer2UserID,
			NIDN:       "0987654321",
			FullName:   "Dr. Siti Nurhaliza, M.T",
			Phone:      "081234567892",
			Department: "Sistem Informasi",
		},
	}

	log.Println("🌱 Seeding lecturers...")

	for _, lecturer := range lecturers {
		// Check if lecturer already exists
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM "lecturer" WHERE nidn = $1 OR user_id = $2)`
		err := s.db.Pool.QueryRow(ctx, checkQuery, lecturer.NIDN, lecturer.UserID).Scan(&exists)
		if err != nil {
			log.Printf("❌ Error checking lecturer %s: %v", lecturer.NIDN, err)
			continue
		}

		if exists {
			log.Printf("⏭️  Lecturer %s already exists, skipping...", lecturer.NIDN)
			continue
		}

		// Insert lecturer
		query := `
			INSERT INTO "lecturer" (user_id, nidn, full_name, phone, department, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`
		_, err = s.db.Pool.Exec(ctx, query, lecturer.UserID, lecturer.NIDN, lecturer.FullName, lecturer.Phone, lecturer.Department)
		if err != nil {
			log.Printf("❌ Error inserting lecturer %s: %v", lecturer.NIDN, err)
			continue
		}

		log.Printf("✅ Lecturer created: %s - %s", lecturer.NIDN, lecturer.FullName)
	}

	log.Println("✅ Lecturer seeding completed!")
	return nil
}

func (s *Seeder) SeedAll() error {
	log.Println("🌱 Starting database seeding...")

	if err := s.SeedUsers(); err != nil {
		return err
	}

	if err := s.SeedLecturers(); err != nil {
		return err
	}

	if err := s.SeedPrograms(); err != nil {
		return err
	}

	log.Println("✅ All seeding completed successfully!")
	return nil
}
