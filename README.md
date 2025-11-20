# MBKM API - Golang + Fiber + PostgreSQL (Native SQL)

Backend API untuk sistem pembelajaran MBKM (Merdeka Belajar Kampus Merdeka) menggunakan Golang dengan Fiber framework dan PostgreSQL (tanpa ORM, native SQL dengan pgx driver).

## 🚀 Features

- ✅ **Fiber Web Framework** - Fast, Express-inspired framework
- ✅ **PostgreSQL Native SQL** - Menggunakan pgx driver tanpa ORM
- ✅ **Auto-Migration** - Generate tables dari struct models
- ✅ **JWT Authentication** - Stateless authentication
- ✅ **Role-Based Access Control** - Admin, Lecturer, Student roles
- ✅ **CRUD Operations** - Users, Programs, Enrollments, Assessments
- ✅ **Middleware** - Auth, CORS, Logger, Recovery
- ✅ **Clean Architecture** - Handlers → Database (simple 2-layer)

## 📁 Project Structure

```
mbkm-api/
├── cmd/
│   └── main.go              # Application entry point
├── config/
│   └── config.go            # Configuration management
├── database/
│   ├── database.go          # PostgreSQL connection & auto-migrate
│   └── schema.sql           # Manual schema (optional)
├── handlers/
│   ├── auth.go              # Authentication handlers
│   ├── program.go           # Program CRUD handlers
│   ├── enrollment.go        # Enrollment handlers
│   └── assessment.go        # Assessment handlers
├── middleware/
│   └── auth.go              # JWT & Role middleware
├── models/
│   └── models.go            # Data models & DTOs
├── routes/
│   └── routes.go            # Route definitions
├── utils/
│   ├── jwt.go               # JWT token generation/validation
│   ├── password.go          # Password hashing (bcrypt)
│   └── response.go          # Standard API responses
├── .env                     # Environment variables
├── .env.example             # Environment template
├── go.mod                   # Go dependencies
├── Makefile                 # Build automation
└── README.md                # This file
```

## 🛠️ Installation

### Prerequisites
- Go 1.21+
- PostgreSQL 15+

### Setup

1. **Clone & Navigate**
```bash
cd mbkm-api
```

2. **Copy Environment File**
```bash
cp .env.example .env
# Edit .env with your database credentials
```

3. **Install Dependencies**
```bash
make install
# atau
go mod download && go mod tidy
```

4. **Create Database**
```bash
createdb -U postgres mbkm_db
```

5. **Run Auto-Migration**
```bash
make migrate
# atau
go run cmd/main.go migrate
```

6. **Start Server**
```bash
make run
# atau
go run cmd/main.go
```

Server akan berjalan di `http://localhost:8080`

## 📊 Database Schema

### Tables
- **users** - User accounts (admin, lecturer, student)
- **programs** - Study programs/courses
- **enrollments** - Student enrollments in programs
- **assessments** - Student grades/assessments

### Auto-Migration
Project ini menggunakan **auto-migration dari struct**. Table akan otomatis dibuat berdasarkan definisi struct di `models/models.go`:

```go
db.AutoMigrate(
    models.User{},
    models.Program{},
    models.Enrollment{},
    models.Assessment{},
)
```

## 🔌 API Endpoints

### Authentication (Public)
```
POST   /api/v1/auth/register   - Register new user
POST   /api/v1/auth/login      - Login user
```

### Authentication (Protected)
```
GET    /api/v1/auth/me         - Get current user profile
```

### Programs (Protected)
```
GET    /api/v1/programs        - Get all programs
GET    /api/v1/programs/:id    - Get program by ID
POST   /api/v1/programs        - Create program (admin/lecturer)
PUT    /api/v1/programs/:id    - Update program (admin/lecturer)
DELETE /api/v1/programs/:id    - Delete program (admin)
```

### Enrollments (Protected)
```
GET    /api/v1/enrollments                  - Get all enrollments (admin/lecturer)
GET    /api/v1/enrollments/student/:id      - Get student enrollments
POST   /api/v1/enrollments                  - Create enrollment
PUT    /api/v1/enrollments/:id/status       - Update enrollment status (admin/lecturer)
DELETE /api/v1/enrollments/:id              - Delete enrollment (admin)
```

### Assessments (Protected)
```
GET    /api/v1/assessments/enrollment/:id   - Get assessments by enrollment
POST   /api/v1/assessments                  - Create assessment (admin/lecturer)
PUT    /api/v1/assessments/:id              - Update assessment (admin/lecturer)
DELETE /api/v1/assessments/:id              - Delete assessment (admin/lecturer)
```

### Health Check
```
GET    /health                 - Server health status
```

## 🔐 Authentication

### Register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "student1",
    "email": "student1@mbkm.ac.id",
    "password": "password123",
    "full_name": "Ahmad Rizki",
    "role": "student"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student1@mbkm.ac.id",
    "password": "password123"
  }'
```

### Use Token
```bash
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🎯 Role-Based Access

- **admin**: Full access to all resources
- **lecturer**: Manage programs, enrollments, assessments
- **student**: View programs, manage own enrollments

## ⚙️ Configuration

Edit `.env` file:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=mbkm_db
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key-min-32-chars
JWT_EXPIRATION=24

# Server
SERVER_PORT=8080
```

## 🧪 Testing

```bash
# Test health endpoint
curl http://localhost:8080/health

# Register user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@mbkm.ac.id","password":"pass123","full_name":"Test User","role":"student"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@mbkm.ac.id","password":"pass123"}'
```

## 📝 Sample Data

Default users (password: `password123`):
- **admin@mbkm.ac.id** (Admin)
- **dosen1@mbkm.ac.id** (Lecturer)
- **mhs1@mbkm.ac.id** (Student)

## 🚧 Development

```bash
# Run with live reload (install air first)
air

# Build binary
make build

# Clean build artifacts
make clean

# Run migration only
make migrate
```

## 📦 Dependencies

- `github.com/gofiber/fiber/v2` - Web framework
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/joho/godotenv` - Environment loader
- `github.com/golang-jwt/jwt/v5` - JWT implementation
- `golang.org/x/crypto` - Password hashing

## 🔧 Troubleshooting

### Database connection error
```bash
# Check PostgreSQL is running
pg_isready

# Create database if not exists
createdb -U postgres mbkm_db

# Test connection
psql -h localhost -U postgres -d mbkm_db
```

### Migration issues
```bash
# Drop and recreate database
dropdb -U postgres mbkm_db
createdb -U postgres mbkm_db

# Run migration again
make migrate
```

### Port already in use
```bash
# Change SERVER_PORT in .env
SERVER_PORT=3000

# Or kill process on port 8080
lsof -ti:8080 | xargs kill -9
```

## 📄 License

MIT License

## 👨‍💻 Author

Workshop MBKM - Golang PostgreSQL API
