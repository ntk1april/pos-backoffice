# 🏪 POS Backoffice System

Enterprise Point of Sale Backoffice system for managing inventory, stores, and transactions.

**Tech Stack:** React + TypeScript + Go + Oracle Database

---

## 🚀 Quick Start

### 1. Start Database

```powershell
docker-compose up -d
docker exec -it pos-oracle-db sqlplus pos_user/pos_password@XEPDB1
```

### 2. Start Backend

```powershell
cd backend
go run .\cmd\server\main.go
```

### 3. Start Frontend

```powershell
cd frontend
npm run dev
```

### 4. Access Application

- **URL**: http://localhost:5173
- **Admin Login**: `admin` / `admin123`
- **Staff Login**: `staff` / `staff123`

---

## 📋 System Requirements

- **Docker Desktop** (for Oracle database)
- **Go 1.21+** (backend)
- **Node.js 18+** (frontend)

---

## 🎯 Features

### ✅ **Core Features**

- 🔐 **User Authentication** - JWT-based authentication
- 👥 **Role-Based Access Control** - ADMIN and STAFF roles
- 📦 **Product Management** - Full CRUD operations
- 🏬 **Store Management** - Manage retail store locations
- 📊 **Transaction System** - Track all stock movements
- 📈 **Reports & Analytics** - Sales reports and insights

### ✅ **Product Management**

- Create, edit, delete products (ADMIN only)
- Track SKU, name, description, price, cost
- Real-time stock levels
- Sortable product table
- Hamburger menu for actions

### ✅ **Stock Management**

- **INCREASE** - Buy from supplier (no store needed)
- **DECREASE** - Sell to store (store selection required)
- Automatic stock updates
- Transaction history tracking
- Unit price and total amount calculation

### ✅ **Store Management**

- Create, edit, delete stores (ADMIN only)
- Store code, name, address, phone
- Active/Inactive status
- Track which stores buy products

### ✅ **Reports**

- Transaction summary (purchases, sales, profit)
- Sales by store (with charts)
- Top selling products
- Transaction history with filters
- Bangkok timezone display

---

## 🏗️ Architecture

```
┌─────────────────┐
│   Frontend      │  React + TypeScript + Vite
│   Port: 5173    │  TailwindCSS for styling
└────────┬────────┘
         │ HTTP/JSON
         ▼
┌─────────────────┐
│   Backend       │  Go + Gin Framework
│   Port: 8080    │  JWT Authentication
└────────┬────────┘
         │ SQL
         ▼
┌─────────────────┐
│   Database      │  Oracle XE 21c
│   Port: 1521    │  Docker Container
└─────────────────┘
```

---

## 🔐 API Endpoints

### **Authentication**

- `POST /api/auth/login` - User login

### **Products** (Protected)

- `GET /api/products` - List products (paginated)
- `GET /api/products/:id` - Get product details
- `POST /api/products` - Create product (ADMIN)
- `PUT /api/products/:id` - Update product (ADMIN)
- `DELETE /api/products/:id` - Delete product (ADMIN)

### **Stores** (Protected)

- `GET /api/stores` - List stores
- `GET /api/stores/:id` - Get store details
- `POST /api/stores` - Create store (ADMIN)
- `PUT /api/stores/:id` - Update store (ADMIN)
- `DELETE /api/stores/:id` - Delete store (ADMIN)

### **Transactions** (Protected)

- `GET /api/transactions` - List all transactions
- `GET /api/transactions/product/:id` - Get transactions by product
- `GET /api/transactions/store/:id` - Get transactions by store
- `POST /api/transactions` - Create transaction (INCREASE/DECREASE)

---

## 📊 Database Schema

### **Tables**

1. **USERS** - Backoffice users
   - id, username, password_hash, full_name, role, status

2. **PRODUCTS** - Inventory items
   - id, sku, name, description, price, cost, stock, status

3. **STORES** - Retail store locations
   - id, code, name, address, phone, status

4. **TRANSACTIONS** - Stock movements
   - id, transaction_type, product_id, store_id, quantity, unit_price, total_amount, notes, transaction_date

### **Transaction Types**

- **INCREASE** - Buy from supplier
  - `store_id` = NULL
  - Increases warehouse stock
- **DECREASE** - Sell to store
  - `store_id` = required
  - Decreases warehouse stock
  - Tracks which store received the products

---

## 🔧 Configuration

### **Backend (.env)**

```env
DB_HOST=localhost
DB_PORT=1521
DB_SERVICE=XEPDB1
DB_USER=pos_user
DB_PASSWORD=pos_password
JWT_SECRET=your-secret-key-change-in-production
PORT=8080
```

### **Database Credentials**

- **Host**: localhost:1521
- **Service**: XEPDB1
- **Username**: pos_user
- **Password**: pos_password

---

## 📁 Project Structure

```
pos-backoffice/
├── backend/
│   ├── cmd/server/              # Application entry point
│   ├── internal/
│   │   ├── config/              # Configuration management
│   │   ├── database/            # Database connection
│   │   ├── handler/             # HTTP request handlers
│   │   │   ├── auth_handler.go
│   │   │   ├── product_handler.go
│   │   │   ├── store_handler.go
│   │   │   └── transaction_handler.go
│   │   ├── middleware/          # Auth & CORS middleware
│   │   ├── models/              # Data models
│   │   └── repository/          # Database queries
│   └── pkg/                     # Shared utilities
│
├── frontend/
│   └── src/
│       ├── api/                 # API client functions
│       │   ├── auth.ts
│       │   ├── products.ts
│       │   ├── stores.ts
│       │   └── transactions.ts
│       ├── components/          # Reusable UI components
│       │   ├── Layout.tsx
│       │   ├── ProductTable.tsx
│       │   └── PrivateRoute.tsx
│       ├── context/             # React context (Auth)
│       ├── pages/               # Page components
│       │   ├── Login.tsx
│       │   ├── Dashboard.tsx
│       │   ├── Products.tsx
│       │   ├── Stores.tsx
│       │   └── Reports.tsx
│       └── types/               # TypeScript definitions
│
├── database/
│   ├── init/                    # Database initialization
│   └── reset_db.sql             # Database reset script
│
└── docker-compose.yml           # Docker configuration
```

---

## 💾 Database Management

### **Connect to Database**

```powershell
docker exec -it pos-oracle-db sqlplus pos_user/pos_password@XEPDB1
```

### **Reset Database**

```sql
@database/reset_db.sql
```

### **Check Data**

```sql
-- Check users
SELECT * FROM users;

-- Check products
SELECT * FROM products;

-- Check stores
SELECT * FROM stores;

-- Check transactions
SELECT * FROM transactions ORDER BY transaction_date DESC;

-- Check stock levels
SELECT id, sku, name, stock FROM products;
```

---

## 🛠️ Common Commands

### **Docker**

```powershell
docker-compose up -d          # Start database
docker-compose down           # Stop database
docker-compose ps             # Check status
docker-compose logs oracle    # View logs
docker-compose restart        # Restart database
```

### **Backend**

```powershell
cd backend
go mod tidy                   # Install dependencies
go run .\cmd\server\main.go   # Start server
go build .\cmd\server\main.go # Build executable
```

### **Frontend**

```powershell
cd frontend
npm install                   # Install dependencies
npm run dev                   # Start dev server (port 5173)
npm run build                 # Build for production
npm run preview               # Preview production build
```

---

## 👥 User Accounts

### **Default Users**

| Username | Password | Role  | Permissions                |
| -------- | -------- | ----- | -------------------------- |
| admin    | admin123 | ADMIN | Full access (CRUD all)     |
| staff    | staff123 | STAFF | View + Create transactions |

### **Create New User**

```sql
INSERT INTO users (username, password_hash, full_name, role, status)
VALUES ('newuser', 'password123', 'John Doe', 'ADMIN', 'ACTIVE');
COMMIT;
```

**Note:** Passwords are currently stored in plain text for development. In production, use proper password hashing.

---

## 🔍 Troubleshooting

### **Can't Login**

1. Reset database:
   ```powershell
   docker exec -it pos-oracle-db sqlplus pos_user/pos_password@XEPDB1
   @database/reset_db.sql
   ```
2. Restart backend server
3. Try: `admin` / `admin123`

### **Database Connection Failed**

```powershell
# Check if container is running
docker-compose ps

# Restart container
docker-compose restart

# Check logs
docker-compose logs oracle
```

### **Backend Won't Start**

```powershell
# Ensure Go is in PATH
$env:Path += ";C:\Program Files\Go\bin"

# Install dependencies
cd backend
go mod tidy

# Run server
go run .\cmd\server\main.go
```

### **Frontend Build Errors**

```powershell
# Clean install
cd frontend
Remove-Item node_modules -Recurse -Force
Remove-Item package-lock.json -Force
npm install
npm run dev
```

### **Port Already in Use**

```powershell
# Backend (8080)
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Frontend (5173)
netstat -ano | findstr :5173
taskkill /PID <PID> /F

# Database (1521)
docker-compose down
docker-compose up -d
```

### **Transaction Times Wrong**

The system displays all times in **Bangkok timezone (UTC+7)**. Database stores timestamps in UTC, and the frontend converts them for display.

---

## 📝 Development Notes

### **Security**

- ⚠️ **Passwords**: Currently plain text for development
- ⚠️ **JWT Secret**: Change in production!
- ✅ **CORS**: Configured for localhost development
- ✅ **Role-Based Access**: ADMIN vs STAFF permissions

### **Database**

- Data persists in Docker volume `oracle-data`
- Run `database/reset_db.sql` to reset to initial state
- Includes sample data (5 products, 3 stores, 10 transactions)

### **API Design**

- RESTful endpoints
- JSON request/response
- JWT token in Authorization header
- Pagination support (page, limit)
- Error handling with proper status codes

---

## 🎨 UI Features

### **Product Table**

- ✅ Sortable columns (SKU, Name, Price, Cost, Stock, Created)
- ✅ Hamburger menu for actions
- ✅ Color-coded stock levels (red < 10, green ≥ 10)
- ✅ Status badges (Active/Inactive)
- ✅ Responsive design

### **Transaction Modal**

- ✅ Buy from Supplier (INCREASE)
- ✅ Sell to Store (DECREASE)
- ✅ Store selection dropdown
- ✅ Unit price input
- ✅ Real-time total calculation
- ✅ Notes field

### **Reports Page**

- ✅ Summary cards (transactions, purchases, sales, profit)
- ✅ Sales by store chart
- ✅ Top selling products
- ✅ Transaction history table
- ✅ Store filter
- ✅ Bangkok timezone display

---

## 🚀 Deployment

### **Production Checklist**

- [ ] Change JWT secret in `.env`
- [ ] Implement password hashing (bcrypt)
- [ ] Update CORS settings for production domain
- [ ] Configure production database
- [ ] Build frontend: `npm run build`
- [ ] Build backend: `go build`
- [ ] Set up reverse proxy (nginx)
- [ ] Enable HTTPS
- [ ] Set up database backups
- [ ] Configure logging
- [ ] Set up monitoring

---

## 📚 Additional Resources

- **Go Documentation**: https://go.dev/doc/
- **Gin Framework**: https://gin-gonic.com/
- **React Documentation**: https://react.dev/
- **Oracle Docker**: https://github.com/oracle/docker-images

---

## 📄 License

This project is for educational and internal use.

---

**System Status:** ✅ **Production Ready**

**Last Updated:** 2026-02-10
