# 📚 Project Structure Guide - For Golang Beginners

This guide explains every file and folder in the POS Backoffice project.

---

## 🏗️ Overall Architecture

```
Frontend (React) → Backend (Golang) → Database (Oracle)
     ↓                  ↓                  ↓
  TypeScript           Gin              Docker
  TailwindCSS       go-ora              XE 21c
```

---

## 📁 Complete Project Structure

```
pos-backoffice/
├── backend/              # Golang backend server
├── frontend/             # React frontend UI
├── database/             # Database scripts
├── docker-compose.yml    # Docker configuration
├── .gitignore           # Git ignore rules
├── .dockerignore        # Docker ignore rules
└── README.md            # Main documentation
```

---

## 🔧 Backend Structure (Golang)

### **Root Files**

#### `backend/go.mod`

**What it is**: Go module definition file (like package.json in Node.js)

**Purpose**:

- Defines module name: `pos-backoffice`
- Lists all dependencies (libraries)
- Specifies Go version (1.21)

**Key dependencies**:

```go
github.com/gin-gonic/gin          // Web framework (like Express.js)
github.com/sijms/go-ora/v2        // Oracle database driver
github.com/golang-jwt/jwt/v5      // JWT authentication
golang.org/x/crypto               // Password hashing (bcrypt)
```

**You don't edit this manually** - use `go mod tidy` to update it.

---

#### `backend/go.sum`

**What it is**: Checksums of all dependencies

**Purpose**: Ensures you download the exact same version of libraries

**You never edit this file** - Go manages it automatically.

---

#### `backend/.env`

**What it is**: Environment configuration file

**Purpose**: Stores sensitive configuration (passwords, secrets)

**Contents**:

```env
DB_HOST=localhost          # Database server address
DB_PORT=1521              # Database port
DB_SERVICE=XEPDB1         # Oracle service name
DB_USER=pos_user          # Database username
DB_PASSWORD=pos_password  # Database password
JWT_SECRET=...            # Secret key for JWT tokens
PORT=8080                 # Backend server port
GIN_MODE=debug            # Development mode
```

**Important**: Never commit this to Git! (It's in .gitignore)

---

#### `backend/start-backend.bat`

**What it is**: Windows batch script to start the backend

**Purpose**: Automates backend startup

- Adds Go to PATH
- Downloads dependencies
- Starts the server

**Usage**: Just double-click or run `.\start-backend.bat`

---

### **Directory: `backend/cmd/`**

This is the **entry point** of your application.

#### `backend/cmd/server/main.go`

**What it is**: The main file that starts everything

**Purpose**:

1. Loads configuration from `.env`
2. Connects to Oracle database
3. Sets up HTTP routes
4. Starts the web server

**How it works**:

```go
func main() {
    // 1. Load config
    config.LoadConfig()

    // 2. Connect to database
    database.InitDB()

    // 3. Setup routes
    router := gin.Default()
    // ... register routes

    // 4. Start server
    router.Run(":8080")
}
```

**This is what you run**: `go run cmd/server/main.go`

---

### **Directory: `backend/internal/`**

Contains all the **business logic** of your application. "Internal" means these packages can only be imported by this project.

---

#### `backend/internal/config/`

**Purpose**: Manages application configuration

**File**: `config.go`

**What it does**:

```go
// Loads .env file
LoadConfig()

// Provides config to other parts
AppConfig.DBHost
AppConfig.JWTSecret

// Builds database connection string
GetDSN() // Returns: oracle://user:pass@host:port/service
```

**Key functions**:

- `LoadConfig()` - Reads .env file
- `GetDSN()` - Creates Oracle connection string
- `getEnv()` - Gets environment variable with default value

---

#### `backend/internal/database/`

**Purpose**: Manages database connection

**File**: `oracle.go`

**What it does**:

```go
// Global database connection
var DB *sql.DB

// Initialize connection
InitDB()

// Close connection
CloseDB()

// Get connection
GetDB()
```

**How it works**:

1. Opens connection to Oracle using `go-ora` driver
2. Configures connection pool (max 25 connections)
3. Tests connection with Ping()
4. Stores in global `DB` variable

**Connection pool settings**:

- `MaxOpenConns: 25` - Maximum 25 simultaneous connections
- `MaxIdleConns: 5` - Keep 5 idle connections ready
- `ConnMaxLifetime: 5 min` - Recycle connections every 5 minutes

---

#### `backend/internal/models/`

**Purpose**: Defines data structures (like database tables)

**Files**:

- `user.go` - User structure
- `product.go` - Product structure
- `stock.go` - Stock log structure
- `request.go` - API request structures
- `response.go` - API response structures

**Example** (`user.go`):

```go
type User struct {
    ID           int64     `json:"id"`
    Username     string    `json:"username"`
    PasswordHash string    `json:"-"`              // "-" means don't send in JSON
    FullName     string    `json:"full_name"`
    Role         string    `json:"role"`
    Status       string    `json:"status"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

**Why separate files?**

- Keeps code organized
- Each file focuses on one entity
- Easy to find and modify

---

#### `backend/internal/repository/`

**Purpose**: Database queries (SQL operations)

**Files**:

- `user_repository.go` - User database operations
- `product_repository.go` - Product database operations
- `stock_repository.go` - Stock database operations

**What it does**:

- Executes SQL queries
- Maps database rows to Go structs
- Handles database errors

**Example** (`user_repository.go`):

```go
type UserRepository struct {
    db *sql.DB
}

// Find user by username
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
    query := `SELECT id, username, password_hash, full_name, role, status
              FROM users WHERE username = :1`

    var user models.User
    err := r.db.QueryRow(query, username).Scan(
        &user.ID, &user.Username, &user.PasswordHash,
        &user.FullName, &user.Role, &user.Status,
    )

    return &user, err
}
```

**Key concepts**:

- `:1, :2, :3` - Oracle bind parameters (prevents SQL injection)
- `QueryRow()` - Execute query, expect one row
- `Scan()` - Map database columns to struct fields
- `&user.ID` - Pointer to struct field

---

#### `backend/internal/service/`

**Purpose**: Business logic (rules and validation)

**Files**:

- `auth_service.go` - Login logic
- `product_service.go` - Product business rules
- `stock_service.go` - Stock management logic

**What it does**:

- Validates input
- Applies business rules
- Calls repository for database operations
- Returns results to handlers

**Example** (`auth_service.go`):

```go
type AuthService struct {
    userRepo *repository.UserRepository
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
    // 1. Find user in database
    user, err := s.userRepo.FindByUsername(req.Username)
    if err != nil {
        return nil, fmt.Errorf("invalid username or password")
    }

    // 2. Verify password
    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
    if err != nil {
        // Fallback to plain text (for testing)
        if user.PasswordHash != req.Password {
            return nil, fmt.Errorf("invalid username or password")
        }
    }

    // 3. Generate JWT token
    token, err := jwt.GenerateToken(user)

    // 4. Return response
    return &models.LoginResponse{
        Token: token,
        User:  *user,
    }, nil
}
```

**Why separate from repository?**

- Repository = database operations only
- Service = business logic + validation
- Keeps concerns separated

---

#### `backend/internal/handler/`

**Purpose**: HTTP request/response handling

**Files**:

- `auth_handler.go` - Login endpoint
- `product_handler.go` - Product CRUD endpoints
- `stock_handler.go` - Stock management endpoints

**What it does**:

- Receives HTTP requests
- Validates input
- Calls service layer
- Returns JSON responses

**Example** (`auth_handler.go`):

```go
type AuthHandler struct {
    authService *service.AuthService
}

func (h *AuthHandler) Login(c *gin.Context) {
    // 1. Parse request body
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request")
        return
    }

    // 2. Call service
    result, err := h.authService.Login(&req)
    if err != nil {
        response.Error(c, http.StatusUnauthorized, err.Error())
        return
    }

    // 3. Return success response
    response.Success(c, "Login successful", result)
}
```

**Key concepts**:

- `c *gin.Context` - Contains request/response data
- `ShouldBindJSON()` - Parse JSON body into struct
- `response.Success()` - Send success JSON response
- `response.Error()` - Send error JSON response

---

#### `backend/internal/middleware/`

**Purpose**: Code that runs **before** your handlers

**Files**:

- `auth.go` - JWT authentication middleware
- `cors.go` - CORS (Cross-Origin) configuration

**What it does**:

- Validates JWT tokens
- Checks user permissions
- Allows/blocks requests

**Example** (`auth.go`):

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Get token from header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            response.Error(c, http.StatusUnauthorized, "Missing token")
            c.Abort()  // Stop request
            return
        }

        // 2. Extract token (remove "Bearer " prefix)
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        // 3. Validate token
        claims, err := jwt.ValidateToken(tokenString)
        if err != nil {
            response.Error(c, http.StatusUnauthorized, "Invalid token")
            c.Abort()
            return
        }

        // 4. Store user info in context
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)

        // 5. Continue to handler
        c.Next()
    }
}
```

**How middleware works**:

```
Request → Middleware → Handler → Response
          ↓
     Check token
     If invalid → Stop (401 error)
     If valid → Continue
```

**Usage in routes**:

```go
// Public route (no middleware)
router.POST("/api/auth/login", authHandler.Login)

// Protected route (requires auth)
protected := router.Group("/api")
protected.Use(middleware.AuthMiddleware())
protected.GET("/products", productHandler.GetProducts)
```

---

### **Directory: `backend/pkg/`**

Contains **reusable utilities** that could be used in other projects.

---

#### `backend/pkg/jwt/`

**Purpose**: JWT token generation and validation

**File**: `jwt.go`

**What it does**:

```go
// Generate token when user logs in
GenerateToken(user *models.User) (string, error)

// Validate token on every API request
ValidateToken(tokenString string) (*Claims, error)
```

**Token structure**:

```go
type Claims struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}
```

**How it works**:

1. **Generate**: Creates token with user info, signs with JWT_SECRET
2. **Validate**: Checks signature, expiration, returns user info

**Token expiration**: 24 hours

---

#### `backend/pkg/response/`

**Purpose**: Standardized JSON responses

**File**: `response.go`

**What it does**:

```go
// Success response
Success(c *gin.Context, message string, data interface{})
// Returns: {"success": true, "message": "...", "data": {...}}

// Error response
Error(c *gin.Context, statusCode int, message string)
// Returns: {"success": false, "message": "..."}
```

**Why use this?**

- Consistent response format
- Less code duplication
- Easy to change format later

---

## 🎯 Request Flow Example

Let's trace a **login request** through the entire backend:

### 1. **User submits login form**

```
POST /api/auth/login
Body: {"username": "admin", "password": "admin123"}
```

### 2. **Request hits `cmd/server/main.go`**

```go
router.POST("/api/auth/login", authHandler.Login)
```

Routes to `authHandler.Login()`

### 3. **Handler** (`internal/handler/auth_handler.go`)

```go
func (h *AuthHandler) Login(c *gin.Context) {
    // Parse JSON body
    var req models.LoginRequest
    c.ShouldBindJSON(&req)

    // Call service
    result, err := h.authService.Login(&req)

    // Return response
    response.Success(c, "Login successful", result)
}
```

### 4. **Service** (`internal/service/auth_service.go`)

```go
func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
    // Get user from database
    user, err := s.userRepo.FindByUsername(req.Username)

    // Verify password
    bcrypt.CompareHashAndPassword(user.PasswordHash, req.Password)

    // Generate JWT token
    token, err := jwt.GenerateToken(user)

    return &models.LoginResponse{Token: token, User: *user}, nil
}
```

### 5. **Repository** (`internal/repository/user_repository.go`)

```go
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
    query := `SELECT id, username, password_hash, full_name, role, status
              FROM users WHERE username = :1`

    var user models.User
    r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, ...)

    return &user, nil
}
```

### 6. **Database** (Oracle)

```sql
SELECT id, username, password_hash, full_name, role, status
FROM users WHERE username = 'admin'
```

### 7. **Response flows back**

```
Database → Repository → Service → Handler → JSON Response
```

### 8. **Frontend receives**

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "username": "admin",
      "full_name": "System Administrator",
      "role": "ADMIN"
    }
  }
}
```

---

## 📊 Architecture Layers

```
┌─────────────────────────────────────────┐
│         HTTP Request (JSON)             │
└─────────────────┬───────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│  HANDLER (internal/handler/)            │
│  - Receives HTTP request                │
│  - Validates input                      │
│  - Returns JSON response                │
└─────────────────┬───────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│  SERVICE (internal/service/)            │
│  - Business logic                       │
│  - Validation rules                     │
│  - Orchestrates operations              │
└─────────────────┬───────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│  REPOSITORY (internal/repository/)      │
│  - SQL queries                          │
│  - Database operations                  │
│  - Data mapping                         │
└─────────────────┬───────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│  DATABASE (Oracle)                      │
│  - Stores data                          │
│  - Executes queries                     │
└─────────────────────────────────────────┘
```

**Why this structure?**

- **Separation of Concerns**: Each layer has one job
- **Testable**: Easy to test each layer independently
- **Maintainable**: Changes in one layer don't affect others
- **Scalable**: Easy to add new features

---

## 🔑 Key Golang Concepts Used

### 1. **Structs**

```go
type User struct {
    ID       int64  `json:"id"`
    Username string `json:"username"`
}
```

Like classes in other languages, but simpler.

### 2. **Pointers**

```go
func (r *UserRepository) FindByUsername(username string) (*models.User, error)
```

- `*UserRepository` - Method on pointer (can modify)
- `*models.User` - Returns pointer to User (efficient)

### 3. **Error Handling**

```go
user, err := r.userRepo.FindByUsername(username)
if err != nil {
    return nil, err
}
```

Go doesn't have try/catch. Always check `err`.

### 4. **Interfaces** (Implicit)

```go
type UserRepository interface {
    FindByUsername(username string) (*models.User, error)
}
```

Any struct with this method implements the interface.

### 5. **Packages**

```go
package handler

import (
    "pos-backoffice/internal/service"
    "pos-backoffice/pkg/response"
)
```

- `package` - Declares package name
- `import` - Imports other packages

### 6. **Dependency Injection**

```go
type AuthHandler struct {
    authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}
```

Pass dependencies through constructor.

---

## 🚀 Common Commands

### Development

```bash
# Install dependencies
go mod tidy

# Run server
go run cmd/server/main.go

# Build executable
go build -o server.exe cmd/server/main.go

# Run executable
./server.exe
```

### Testing

```bash
# Run all tests
go test ./...

# Test specific package
go test ./internal/service

# Test with coverage
go test -cover ./...
```

### Dependencies

```bash
# Add dependency
go get github.com/some/package

# Update all dependencies
go get -u ./...

# Remove unused dependencies
go mod tidy
```

---

## 📝 Best Practices

### 1. **Error Handling**

```go
// ✅ Good
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// ❌ Bad
if err != nil {
    panic(err)  // Don't panic in production code
}
```

### 2. **SQL Injection Prevention**

```go
// ✅ Good - Use bind parameters
query := `SELECT * FROM users WHERE username = :1`
db.QueryRow(query, username)

// ❌ Bad - String concatenation
query := `SELECT * FROM users WHERE username = '` + username + `'`
```

### 3. **Resource Cleanup**

```go
// ✅ Good - Always close resources
rows, err := db.Query(query)
if err != nil {
    return err
}
defer rows.Close()  // Ensures cleanup
```

### 4. **Struct Tags**

```go
type User struct {
    ID       int64  `json:"id" db:"id"`
    Password string `json:"-"`  // Don't send in JSON
}
```

---

## 🎓 Learning Resources

### Official Go Documentation

- https://go.dev/tour/ - Interactive Go tutorial
- https://go.dev/doc/effective_go - Best practices

### Libraries Used

- **Gin**: https://gin-gonic.com/docs/
- **go-ora**: https://github.com/sijms/go-ora
- **JWT**: https://github.com/golang-jwt/jwt

---

## 💡 Summary

**This project follows Clean Architecture**:

1. **Handler** - HTTP layer (receives requests)
2. **Service** - Business logic layer
3. **Repository** - Database layer
4. **Models** - Data structures

**Key files to understand**:

1. `cmd/server/main.go` - Entry point
2. `internal/handler/*` - HTTP endpoints
3. `internal/service/*` - Business logic
4. `internal/repository/*` - Database queries

**Start learning from**:

1. Read `cmd/server/main.go` to see how it starts
2. Follow a request through handler → service → repository
3. Understand how JWT authentication works
4. Experiment by adding new endpoints

**You're ready to start coding in Go!** 🚀
