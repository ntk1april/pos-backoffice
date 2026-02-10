# POS Backoffice System

Enterprise Point of Sale Backoffice system with Golang, React, and Oracle Database.

---

## 🚀 Quick Start

### 1. Start Database

```powershell
docker-compose up -d
```

### 2. Start Backend

```powershell
cd backend
.\start-backend.bat
```

### 3. Start Frontend

```powershell
cd frontend
npm run dev
```

### 4. Access Application

- **URL**: http://localhost:5173
- **Login**: `admin` / `admin123`

---

## 📋 Requirements

- Docker Desktop
- Go 1.21+
- Node.js 18+

---

## 🔧 Configuration

### Backend (.env)

```env
DB_HOST=localhost
DB_PORT=1521
DB_SERVICE=XEPDB1
DB_USER=pos_user
DB_PASSWORD=pos_password
JWT_SECRET=your-secret-key
PORT=8080
```

### Database Credentials

- **Host**: localhost:1521
- **Service**: XEPDB1
- **User**: pos_user
- **Password**: pos_password

---

## 📊 Database Management

### Connect to Database

```powershell
docker exec -it pos-oracle-db sqlplus pos_user/pos_password@XEPDB1
```

### Reset Database

```sql
@database/reset_db.sql
```

### Check Data

```sql
SELECT * FROM users;
SELECT * FROM products;
SELECT * FROM stock_logs;
```

---

## 👥 User Accounts

### Default Users

| Username | Password | Role  |
| -------- | -------- | ----- |
| admin    | admin123 | ADMIN |
| staff    | staff123 | STAFF |

### Create New User

```sql
INSERT INTO users (username, password_hash, full_name, role, status)
VALUES ('newuser', 'password123', 'Full Name', 'ADMIN', 'ACTIVE');
COMMIT;
```

---

## 🛠️ Common Commands

### Docker

```powershell
docker-compose up -d          # Start
docker-compose down           # Stop
docker-compose ps             # Status
docker-compose logs oracle    # Logs
```

### Backend

```powershell
cd backend
go mod tidy                   # Install dependencies
go run cmd/server/main.go     # Start server
```

### Frontend

```powershell
cd frontend
npm install                   # Install dependencies
npm run dev                   # Start dev server
npm run build                 # Build for production
```

---

## 🔐 API Endpoints

### Authentication

- `POST /api/auth/login` - User login

### Products

- `GET /api/products` - List products (paginated)
- `GET /api/products/:id` - Get product
- `POST /api/products` - Create product (ADMIN)
- `PUT /api/products/:id` - Update product (ADMIN)
- `DELETE /api/products/:id` - Delete product (ADMIN)

### Stock

- `POST /api/stock/increase` - Increase stock
- `POST /api/stock/decrease` - Decrease stock
- `GET /api/stock/logs/:product_id` - Stock history

---

## 🏗️ Architecture

```
Frontend (React) → Backend (Go/Gin) → Database (Oracle)
                      ↓
                    JWT Auth
                      ↓
              Role-Based Access
```

### Tech Stack

- **Backend**: Golang, Gin, go-ora
- **Frontend**: React, TypeScript, TailwindCSS
- **Database**: Oracle XE 21c (Docker)
- **Auth**: JWT with bcrypt

---

## 📁 Project Structure

```
pos-backoffice/
├── backend/
│   ├── cmd/server/          # Entry point
│   ├── internal/
│   │   ├── config/          # Configuration
│   │   ├── database/        # DB connection
│   │   ├── handler/         # HTTP handlers
│   │   ├── middleware/      # Auth, CORS
│   │   ├── models/          # Data models
│   │   ├── repository/      # Database queries
│   │   └── service/         # Business logic
│   └── pkg/                 # Utilities (JWT)
├── frontend/
│   └── src/
│       ├── api/             # API client
│       ├── components/      # UI components
│       ├── context/         # Auth context
│       ├── pages/           # Pages
│       └── types/           # TypeScript types
├── database/
│   ├── init/                # Init scripts
│   └── reset_db.sql         # Reset script
└── docker-compose.yml       # Docker config
```

---

## 🔍 Troubleshooting

### Can't Login

1. Reset database: `@database/reset_db.sql`
2. Restart backend
3. Try: `admin` / `admin123`

### Database Connection Failed

```powershell
docker-compose ps             # Check if running
docker-compose restart        # Restart container
```

### Backend Won't Start

```powershell
$env:Path += ";C:\Program Files\Go\bin"
go mod tidy
go run cmd/server/main.go
```

### Port Already in Use

```powershell
# Stop backend: Ctrl+C
# Stop frontend: Ctrl+C
# Stop database: docker-compose down
```

---

## 📝 Notes

- **Passwords**: Currently using plain text for testing (admin123, staff123)
- **JWT Secret**: Change in production!
- **Database**: Data persists in Docker volume
- **Reset DB**: Run `database/reset_db.sql` to start fresh

---

## 🎯 Features

✅ User authentication with JWT
✅ Role-based access control (ADMIN/STAFF)
✅ Product management (CRUD)
✅ Stock management (increase/decrease)
✅ Stock history/audit trail
✅ Pagination and search
✅ Responsive UI
✅ Docker deployment

---

**System is ready to use!** 🚀
