# Virtual Cuppa Backend

Backend API in Go using Gin framework with JWT authentication and PostgreSQL database.

## Requirements

- Go 1.21 or newer
- Docker and Docker Compose (for database)

## Installation

1. Install dependencies:

```bash
go mod download
```

2. Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

3. Start PostgreSQL database:

```bash
docker-compose up -d
```

## Running

```bash
go run main.go
```

Server will start on `http://localhost:8080`

## Project Structure

```
virtual-cuppa-be/
├── cmd/             # CLI tools
│   └── migrate/     # Migration CLI
├── config/          # Configuration (database, migrations)
├── handlers/        # HTTP handlers
├── middleware/      # Middleware (authorization)
├── migrations/      # SQL migration files
├── models/          # Data models
├── repositories/    # Data access layer
├── scripts/         # Helper scripts
├── services/        # Business logic layer
├── utils/           # Helper utilities (JWT, responses)
├── main.go          # Application entry point
├── docker-compose.yml
└── .env             # Environment variables
```

## Database Migrations

This project uses **golang-migrate** for database migrations.

### Running Migrations

Migrations run automatically when the application starts. To manually manage migrations:

**Run all pending migrations:**

```bash
go run cmd/migrate/main.go migrate -up
```

**Rollback last N migrations:**

```bash
go run cmd/migrate/main.go migrate -down 1
```

**Check current version:**

```bash
go run cmd/migrate/main.go migrate -version
```

### Creating New Migrations

**Windows (PowerShell):**

```powershell
.\scripts\create-migration.ps1 add_new_column
```

**Linux/Mac:**

```bash
chmod +x scripts/create-migration.sh
./scripts/create-migration.sh add_new_column
```

This will create two files:

- `migrations/000002_add_new_column.up.sql` - Forward migration
- `migrations/000002_add_new_column.down.sql` - Rollback migration

### Migration Best Practices

1. **Always test migrations locally first**
2. **Create rollback migrations** (down.sql)
3. **Never edit existing migrations** - create new ones
4. **Use transactions** where possible
5. **Backup production database** before running migrations
6. **Keep migrations small and focused**

Example migration:

```sql
-- 000002_add_avatar_column.up.sql
ALTER TABLE users ADD COLUMN avatar VARCHAR(255);

-- 000002_add_avatar_column.down.sql
ALTER TABLE users DROP COLUMN avatar;
```

## Project Structure

└── .env # Environment variables

````

## Available Endpoints

### Public

#### POST /api/auth/register

Register a new user.

**Body:**

```json
{
  "firstName": "John",
  "lastName": "Doe",
  "email": "user@example.com",
  "accountType": "User",
  "organisation": "My Company"
}
````

**Response:**

```json
{
  "token": "jwt.access.token",
  "refreshToken": "refresh.token.here",
  "user": {
    "id": 1,
    "firstName": "John",
    "lastName": "Doe",
    "email": "user@example.com",
    "accountType": "User",
    "organisation": "My Company",
    "isConfirmed": true,
    "createdAt": "2025-12-09T10:00:00Z",
    "updatedAt": "2025-12-09T10:00:00Z"
  }
}
```

#### POST /api/auth/request-code

Request a confirmation code to be sent via email (Step 1 of login).

**Body:**

```json
{
  "email": "user@example.com"
}
```

**Response:**

```json
{
  "message": "Confirmation code sent to email"
}
```

**Note:** The confirmCode will be sent to the user's email and is valid for 5 minutes.

#### POST /api/auth/login

Verify confirmation code and login (Step 2 of login).

**Body:**

```json
{
  "email": "user@example.com",
  "confirmCode": "123456"
}
```

**Response:**

```json
{
  "token": "jwt.access.token",
  "refreshToken": "refresh.token.here",
  "user": { ... }
}
```

#### POST /api/auth/refresh

Refresh access token using refresh token.

**Body:**

```json
{
  "refreshToken": "refresh.token.here"
}
```

**Response:**

```json
{
  "token": "new.jwt.access.token",
  "refreshToken": "new.refresh.token",
  "user": { ... }
}
```

### Protected (require JWT token)

Add header: `Authorization: Bearer <token>`

#### GET /api/profile

Get current user profile.

### Admin Only

Admin endpoints require both authentication and admin role.

#### POST /api/admin/import-csv

Import users from CSV file (admin with organisation only).

**Form Data:**

- `file`: CSV file with format: `firstName,lastName,email`

**Example CSV:**

```csv
firstName,lastName,email
John,Doe,john.doe@example.com
Jane,Smith,jane.smith@example.com
```

**cURL:**

```bash
curl -X POST http://localhost:8080/api/admin/import-csv \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -F "file=@users.csv"
```

**Response:**

```json
{
  "message": "Users imported successfully",
  "count": 2
}
```

#### POST /api/admin/confirm-user

Confirm imported user (admin only, same organisation).

**Body:**

```json
{
  "userId": 5
}
```

**Response:**

```json
{
  "message": "User confirmed successfully"
}
```

#### GET /api/admin/users?organisation=MyCompany

Get all users in an organisation.

**Response:**

```json
{
  "users": [...],
  "count": 10
}
```

#### GET /api/admin/dashboard

Admin dashboard endpoint.

### Health Check

#### GET /health

Check if server is running.

**Response:**

```json
{
  "status": "ok",
  "message": "Server is running"
}
```

## User Model

```go
{
  "firstName": string,
  "lastName": string,
  "email": string,
  "confirmCode": string (nullable),    // 6-digit code stored in cache (5min TTL)
  "accountType": "User" | "Admin",
  "organisation": string (optional),
  "isConfirmed": boolean,              // default: true for register, false for CSV import
  "refreshToken": string               // stored securely
}
```

## Authentication Flow

### Registration:

1. User submits registration data (no password required)
2. System generates 6-digit confirmCode
3. ConfirmCode stored in cache with 5-minute TTL
4. Email sent to user with confirmCode via SendGrid
5. User receives JWT tokens immediately

### Login (Two-Step Process):

**Step 1 - Request Code:**

1. User provides email via POST `/api/auth/request-code`
2. System generates 6-digit confirmCode
3. ConfirmCode stored in cache with 5-minute TTL
4. Email sent to user with confirmCode

**Step 2 - Verify Code:**

1. User submits email + confirmCode via POST `/api/auth/login`
2. System verifies code from cache
3. Code deleted after successful verification
4. User receives JWT tokens

## Environment Variables

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=virtual_cuppa
JWT_SECRET=your-secret-key
PORT=8080
GIN_MODE=debug

# SendGrid Configuration
SENDGRID_API_KEY=your-sendgrid-api-key-here
CONFIRM_CODE_TEMPLATE_ID=your-sendgrid-template-id-here
```

## Building

```bash
go build -o virtual-cuppa-be.exe
```

## Docker

Start database:

```bash
docker-compose up -d
```

Stop database:

```bash
docker-compose down
```

## API Testing Examples

### Register:

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"firstName":"John","lastName":"Doe","email":"test@example.com","accountType":"User"}'
```

### Request Login Code:

```bash
curl -X POST http://localhost:8080/api/auth/request-code \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'
```

### Login:

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","confirmCode":"123456"}'
```

**Note:** The confirmCode must be obtained from the email sent during registration or by requesting a new code.

### Refresh Token:

```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refreshToken":"YOUR_REFRESH_TOKEN"}'
```

### Profile (with token):

```bash
curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Import CSV (admin):

```bash
curl -X POST http://localhost:8080/api/admin/import-csv \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -F "file=@users.csv"
```

### Confirm User (admin):

```bash
curl -X POST http://localhost:8080/api/admin/confirm-user \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"userId":5}'
```

## Features

- ✅ JWT Authentication with Access & Refresh Tokens
- ✅ Role-based Authorization (User/Admin)
- ✅ User Registration & Login
- ✅ ConfirmCode Authentication (6-digit code)
- ✅ CSV Import for Bulk User Creation (Admin only)
- ✅ User Confirmation System
- ✅ Organisation-based User Management
- ✅ Layered Architecture (Handler → Service → Repository)
- ✅ PostgreSQL with GORM
- ✅ Database Migrations with golang-migrate
- ✅ Docker Compose for Development
- ✅ CORS Support

## Migration History

### 000001_create_users_table

- Creates `users` table with all fields
- Adds indexes for `email` and `deleted_at`
- Supports soft deletes

### 000002_remove_password_add_confirm_code

- Removes `password` column from `users` table
- Adds `confirm_code` VARCHAR(10) column
- Migrates from password-based to confirmCode authentication
