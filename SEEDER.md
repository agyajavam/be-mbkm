# Database Seeder

Seeder untuk mengisi database dengan data sample untuk development dan testing.

## Available Commands

### Seed Semua Data
```bash
make seed
```
Menjalankan semua seeder (users + programs).

### Seed Users Only
```bash
make seed-users
```
Hanya menjalankan seeder untuk users.

### Seed Programs Only
```bash
make seed-programs
```
Hanya menjalankan seeder untuk programs (memerlukan users sudah di-seed).

## Seeded Data

### Users (6 users)

#### 1. Admin
- **Email**: `admin@mbkm.ac.id`
- **Password**: `admin123`
- **Role**: `admin`
- **Full Name**: Administrator

#### 2. Lecturer 1
- **Email**: `lecturer1@mbkm.ac.id`
- **Password**: `lecturer123`
- **Role**: `lecturer`
- **Full Name**: Dr. Budi Santoso

#### 3. Lecturer 2
- **Email**: `lecturer2@mbkm.ac.id`
- **Password**: `lecturer123`
- **Role**: `lecturer`
- **Full Name**: Dr. Siti Nurhaliza

#### 4. Student 1
- **Email**: `student1@mbkm.ac.id`
- **Password**: `student123`
- **Role**: `student`
- **Full Name**: Ahmad Fauzi

#### 5. Student 2
- **Email**: `student2@mbkm.ac.id`
- **Password**: `student123`
- **Role**: `student`
- **Full Name**: Siti Rahmawati

#### 6. Student 3
- **Email**: `student3@mbkm.ac.id`
- **Password**: `student123`
- **Role**: `student`
- **Full Name**: Andi Wijaya

### Programs (5 programs)

#### 1. Studi Independen - Web Development
- **Code**: MBKM001
- **Credits**: 20
- **Semester**: 5
- **Lecturer**: Dr. Budi Santoso

#### 2. Magang Industri - Software Engineering
- **Code**: MBKM002
- **Credits**: 20
- **Semester**: 6
- **Lecturer**: Dr. Budi Santoso

#### 3. Kampus Mengajar - Pendidikan Digital
- **Code**: MBKM003
- **Credits**: 20
- **Semester**: 5
- **Lecturer**: Dr. Siti Nurhaliza

#### 4. Studi Independen - Data Science
- **Code**: MBKM004
- **Credits**: 20
- **Semester**: 6
- **Lecturer**: Dr. Siti Nurhaliza

#### 5. Proyek Kemanusiaan - Tech for Good
- **Code**: MBKM005
- **Credits**: 20
- **Semester**: 7
- **Lecturer**: Dr. Budi Santoso

## Testing Credentials

### Login as Admin
```json
{
  "email": "admin@mbkm.ac.id",
  "password": "admin123"
}
```

### Login as Lecturer
```json
{
  "email": "lecturer1@mbkm.ac.id",
  "password": "lecturer123"
}
```

### Login as Student
```json
{
  "email": "student1@mbkm.ac.id",
  "password": "student123"
}
```

## Notes

- Seeder akan **skip data yang sudah ada** (tidak akan duplikat)
- Password di-hash menggunakan bcrypt
- Semua users dibuat dengan status `is_active = true`
- Programs dibuat dengan status `is_active = true`

## Development Workflow

1. **Setup database baru**:
   ```bash
   make migrate
   make seed
   ```

2. **Reset database**:
   ```bash
   # Drop and recreate database
   psql -h localhost -U your_user -d postgres -c "DROP DATABASE IF EXISTS mbkm_db;"
   psql -h localhost -U your_user -d postgres -c "CREATE DATABASE mbkm_db;"
   
   # Run migration and seed
   make migrate
   make seed
   ```

3. **Update hanya users**:
   ```bash
   make seed-users
   ```

## Customization

Untuk menambah atau mengubah data seeder, edit file:
```
database/seeder.go
```

Struktur seeder:
- `SeedUsers()` - Seed user data
- `SeedPrograms()` - Seed program data
- `SeedAll()` - Run all seeders
