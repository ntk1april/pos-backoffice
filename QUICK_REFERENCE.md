# 📋 Quick Reference Card

## 🚀 First Time Setup (Run Once)

```powershell
# 1. Start Docker
docker-compose up -d

# 2. Complete Setup (Create User + Tables)
# Wait 2-3 minutes for Oracle to be ready
docker exec -it pos-oracle-db sqlplus system/oracle@XEPDB1
@database/init/setup.sql
exit

# 3. Start backend (new terminal)
cd backend
go run .\cmd\server\main.go

# 4. Start frontend (new terminal)
cd frontend
npm run dev

# 5. Open browser
http://localhost:5173
Login: admin / admin123
```

---

## 🔄 Daily Use (After First Setup)

```powershell
# Terminal 1: Start database
docker-compose up -d

# Terminal 2: Start backend
cd backend
go run .\cmd\server\main.go

# Terminal 3: Start frontend
cd frontend
npm run dev

# Browser
http://localhost:5173
```

---

## 🛑 Stop Everything

```powershell
# Frontend: Ctrl+C
# Backend: Ctrl+C
# Database:
docker-compose down
```

---

## 🔧 Common Commands

### Database

```powershell
# Connect to database
docker exec -it pos-oracle-db sqlplus pos_user/pos_password@XEPDB1

# Reset database (inside SQLPlus)
@database/reset_db.sql

# Check data
SELECT * FROM users;
SELECT * FROM products;
SELECT * FROM stores;
SELECT * FROM transactions;
```

### Backend

```powershell
cd backend
go mod tidy              # Install dependencies
go run .\cmd\server\main.go  # Start server
```

### Frontend

```powershell
cd frontend
npm install              # Install dependencies
npm run dev              # Start dev server
npm run build            # Build for production
```

---

## 👥 Default Users

| Username | Password | Role  |
| -------- | -------- | ----- |
| admin    | admin123 | ADMIN |
| staff    | staff123 | STAFF |

---

## 🌐 URLs

- **Frontend:** http://localhost:5173
- **Backend API:** http://localhost:8080
- **Health Check:** http://localhost:8080/health

---

## 📊 Key Features

### Products

- View all products
- Create/Edit/Delete (ADMIN only)
- Sort by any column
- Hamburger menu for actions

### Transactions

- **Buy from Supplier** (INCREASE)
  - No store needed
  - Increases stock
- **Sell to Store** (DECREASE)
  - Store selection required
  - Decreases stock

### Reports

- Transaction summary
- Sales by store
- Top selling products
- Transaction history

---

## 🔍 Troubleshooting

### Can't login?

```powershell
docker exec -it pos-oracle-db sqlplus pos_user/pos_password@XEPDB1
@database/reset_db.sql
exit
```

### Database connection failed?

```powershell
docker-compose restart
```

### Port already in use?

```powershell
# Find and kill process
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

---

## 📁 Important Files

- `database/init/setup.sql` - First time user setup
- `database/reset_db.sql` - Reset database to initial state
- `backend/.env` - Backend configuration
- `docker-compose.yml` - Docker configuration

---

**For detailed documentation, see:**

- `SETUP_GUIDE.md` - Complete setup instructions
- `README.md` - Full system documentation
