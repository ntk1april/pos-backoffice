# 🚀 First Time Setup Guide

Follow these steps to set up the POS Backoffice system for the first time.

---

## 📋 Prerequisites

Make sure you have installed:

- ✅ Docker Desktop
- ✅ Go 1.21+
- ✅ Node.js 18+

---

## 🔧 Step-by-Step Setup

### **Step 1: Start Docker Database**

```powershell
# Navigate to project directory
cd c:\Users\User\Desktop\pos-backoffice

# Start Oracle database container
docker-compose up -d
```

**Wait 2-3 minutes** for Oracle to fully start. You can check the logs:

```powershell
docker-compose logs -f oracle
```

Look for: `DATABASE IS READY TO USE!`

---

### **Step 2: Run Complete Setup**

This single script will create the user AND initialize the database tables.

Connect to Oracle as SYSTEM user:

```powershell
docker exec -it pos-oracle-db sqlplus system/oracle@XEPDB1
```

Run the setup script:

```sql
@database/init/setup.sql
```

You should see:

```
Creating pos_user...
User created successfully.
Connecting as pos_user...
Connected.
Creating tables and data...
SETUP COMPLETE SUCCESSFULLY!
```

Exit SQLPlus:

```sql
exit
```

---

### **Step 3: Start Backend Server**

Open a new PowerShell terminal:

```powershell
cd c:\Users\User\Desktop\pos-backoffice\backend

# Install dependencies (first time only)
go mod tidy

# Start server
go run .\cmd\server\main.go
```

You should see:

```
✓ Database connection established
Server running on :8080
```

**Keep this terminal open!**

---

### **Step 4: Start Frontend**

Open another PowerShell terminal:

```powershell
cd c:\Users\User\Desktop\pos-backoffice\frontend

# Install dependencies (first time only)
npm install

# Start dev server
npm run dev
```

You should see:

```
VITE v5.4.21  ready in XXX ms
➜  Local:   http://localhost:5173/
```

**Keep this terminal open!**

---

### **Step 5: Access the Application**

Open your browser and go to:

**http://localhost:5173**

Login with:

- **Username:** `admin`
- **Password:** `admin123`

---

## ✅ Verification Checklist

After setup, verify everything works:

- [ ] Docker container is running: `docker-compose ps`
- [ ] Backend server is running on port 8080
- [ ] Frontend dev server is running on port 5173
- [ ] Can login with admin/admin123
- [ ] Can see 5 products in Products page
- [ ] Can see 3 stores in Stores page
- [ ] Can see transactions in Reports page

---

## 🔍 Troubleshooting

### **Database won't start**

```powershell
# Stop and remove containers
docker-compose down -v

# Start fresh
docker-compose up -d

# Wait 2-3 minutes and check logs
docker-compose logs -f oracle
```

### **Can't connect to database**

```powershell
# Check if container is running
docker-compose ps

# Should show pos-oracle-db as "Up"
```

### **User already exists error**

The setup script will automatically drop and recreate the user. If you still get errors:

```sql
-- Connect as SYSTEM
docker exec -it pos-oracle-db sqlplus system/oracle@XEPDB1

-- Manually drop user
DROP USER pos_user CASCADE;

-- Then run setup.sql again
@database/init/setup.sql
```

### **Backend connection error**

Make sure:

1. Docker database is running
2. User was created successfully
3. Database was initialized with reset_db.sql

Check backend `.env` file:

```env
DB_HOST=localhost
DB_PORT=1521
DB_SERVICE=XEPDB1
DB_USER=pos_user
DB_PASSWORD=pos_password
```

---

## 📝 Quick Reference

### **Start Everything**

```powershell
# Terminal 1: Database
docker-compose up -d

# Terminal 2: Backend
cd backend
go run .\cmd\server\main.go

# Terminal 3: Frontend
cd frontend
npm run dev
```

### **Stop Everything**

```powershell
# Stop frontend: Ctrl+C in frontend terminal
# Stop backend: Ctrl+C in backend terminal
# Stop database:
docker-compose down
```

### **Reset Database**

```powershell
docker exec -it pos-oracle-db sqlplus pos_user/pos_password@XEPDB1
@database/reset_db.sql
exit
```

---

## 🎉 Success!

If you can login and see the dashboard, you're all set!

**Next Steps:**

- Explore the Products page
- Try creating a new product (ADMIN only)
- Create a transaction (Buy from Supplier or Sell to Store)
- Check the Reports page for analytics

---

**Need help?** Check the main README.md for detailed documentation.
