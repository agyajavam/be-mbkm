.PHONY: run build clean migrate help

run:
	@echo "🚀 Starting server..."
	@go run cmd/main.go

build:
	@echo "🔨 Building binary..."
	@go build -o bin/mbkm-api cmd/main.go
	@echo "✅ Build complete: bin/mbkm-api"

clean:
	@echo "🧹 Cleaning..."
	@rm -rf bin/
	@echo "✅ Clean complete"

install:
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy

migrate:
	@echo "🔄 Running migrations..."
	@go run cmd/main.go migrate

seed:
	@echo "🌱 Running seeders..."
	@go run cmd/main.go seed

seed-users:
	@echo "🌱 Running user seeder..."
	@go run cmd/main.go seed:users

seed-programs:
	@echo "🌱 Running program seeder..."
	@go run cmd/main.go seed:programs

seed-lecturers:
	@echo "🌱 Running lecturer seeder..."
	@go run cmd/main.go seed:lecturers

swagger:
	@echo "📚 Generating Swagger documentation..."
	@swag init -g cmd/main.go -o docs
	@echo "✅ Swagger docs generated at docs/"

help:
	@echo "Available commands:"
	@echo "  make run            - Run the application"
	@echo "  make build          - Build binary"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make install        - Install dependencies"
	@echo "  make migrate        - Run database migrations"
	@echo "  make seed           - Run all seeders"
	@echo "  make seed-users     - Run user seeder only"
	@echo "  make seed-lecturers - Run lecturer seeder only"
	@echo "  make seed-programs  - Run program seeder only"
	@echo "  make swagger        - Generate Swagger documentation"
