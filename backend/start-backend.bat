@echo off
echo ========================================
echo POS Backoffice - Backend Startup
echo ========================================
echo.

REM Add Go to PATH
set PATH=%PATH%;C:\Program Files\Go\bin

echo Step 1: Downloading dependencies...
go mod tidy

echo.
echo Step 2: Copying environment configuration...
if not exist .env (
    copy .env.docker .env
    echo Created .env file
) else (
    echo .env file already exists
)

echo.
echo Step 3: Starting backend server...
echo.
echo ========================================
echo Server will start on http://localhost:8080
echo Press Ctrl+C to stop
echo ========================================
echo.

go run cmd/server/main.go

pause
