# 🏗️ Understanding the Go Project Structure

This guide explains the folder structure and architectural patterns used in this Go backend. It is designed for beginners to quickly understand where files are located and how the application works.

---

## 📂 Directory Layout

The project follows the **Standard Go Project Layout**, a widely adopted convention in the Go community.

```
backend/
├── cmd/                # Main applications (entry points)
│   └── server/         # The API server application
│       └── main.go     # 🚀 Starting point of the app
│
├── internal/           # Private application code (library code tailored to this app)
│   ├── config/         # ⚙️ Configuration loading (env vars)
│   ├── database/       # 🗄️ Database connection logic
│   ├── handler/        # 🌐 HTTP Handlers (Controllers) - Handle requests & responses
│   ├── middleware/     # 🛡️ Middleware (Auth, CORS, Logging) - Runs before handlers
│   ├── models/         # 📦 Data Structures (Structs for DB and JSON)
│   ├── repository/     # 🗃️ Data Access Layer (SQL queries) - Talks to the DB
│   └── service/        # 🧠 Business Logic Layer - Complex rules (optional)
│
├── pkg/                # Public library code (can be used by external apps, but here just shared utils)
│   └── response/       # 📤 Standardized API response format
```

---

## 🔄 Request Flow (Architecture)

When a request comes in (e.g., `GET /api/products`), it flows through these layers:

1. **Router (`main.go`)**
   - Directs the request to the correct Handler.
   - Applies Middleware (like Authentication).

2. **Middleware (`internal/middleware`)**
   - Checks if the user is logged in (JWT check).
   - Handles CORS (Cross-Origin Resource Sharing).

3. **Handler (`internal/handler`)**
   - **Input:** Parses the JSON body or URL parameters.
   - **Validation:** Checks if data is correct.
   - **Call:** Calls the Service or Repository to get work done.
   - **Output:** Sends the JSON response back to the client.

4. **Service (`internal/service`)** _(Optional Layer)_
   - Contains business logic (e.g., "If stock < 0, return error").
   - Orchestrates multiple repositories if needed.
   - _Note: Simple features might skip this and go straight to Repository._

5. **Repository (`internal/repository`)**
   - **The only place that talks to the Database.**
   - Executors SQL queries (`SELECT`, `INSERT`, `UPDATE`).
   - Maps database rows to Go Structs (`models`).

6. **Database (`internal/database` & Oracle)**
   - The actual storage (Oracle Database in Docker).

---

## 📝 Key Components Explained

### 1. `main.go` (The Entry Point)

This file initializes everything. It loads config, connects to the DB, sets up the router (Gin), and registers routes.

```go
func main() {
    // 1. Load Config
    cfg := config.Load()

    // 2. Connect DB
    db := database.Connect(cfg)

    // 3. Setup Router
    r := gin.Default()

    // 4. Register Routes
    r.GET("/products", productHandler.GetProducts)
}
```

### 2. Models (`internal/models`)

Defines what data looks like. Think of these as the "Shapes" of your data.

```go
type Product struct {
    ID    int     `json:"id" db:"id"`
    Name  string  `json:"name" db:"name"`
    Price float64 `json:"price" db:"price"`
}
```

### 3. Repository (`internal/repository`)

Handles raw SQL. We use `database/sql` or a helper library.

```go
func (r *ProductRepository) GetAll() ([]models.Product, error) {
    rows, _ := r.db.Query("SELECT id, name, price FROM products")
    // ... logic to scan rows into structs ...
    return products, nil
}
```

### 4. Handler (`internal/handler`)

Handles the HTTP request.

```go
func (h *ProductHandler) GetProducts(c *gin.Context) {
    products, err := h.repo.GetAll()
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to fetch"})
        return
    }
    c.JSON(200, products)
}
```

---

## 💡 Why this structure?

- **Separation of Concerns:** DB code is not mixed with HTTP code.
- **Testability:** You can easily test business logic without a real database.
- **Scalability:** Easy to add new features without breaking existing code.

This structure is standard in the Go industry and is excellent for building robust REST APIs.
