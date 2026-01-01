@echo off
REM Setup script for 4 in a Row game (Windows)
REM This script helps with local development setup

echo.
echo 🎮 Setting up 4 in a Row Game...
echo.

REM Check prerequisites
echo 📋 Checking prerequisites...

where docker >nul 2>nul
if errorlevel 1 (
    echo ❌ Docker is not installed. Please install Docker Desktop for Windows.
    pause
    exit /b 1
)

where docker-compose >nul 2>nul
if errorlevel 1 (
    echo ❌ Docker Compose is not installed. Please install Docker Desktop for Windows.
    pause
    exit /b 1
)

echo ✅ Docker and Docker Compose found
echo.

REM Create environment files if they don't exist
echo 📄 Setting up environment files...

if not exist "backend\.env" (
    copy backend\.env.example backend\.env
    echo ✅ Created backend\.env
) else (
    echo ⚠️  backend\.env already exists
)

if not exist "analytics\.env" (
    copy analytics\.env.example analytics\.env
    echo ✅ Created analytics\.env
) else (
    echo ⚠️  analytics\.env already exists
)

REM Build and start services
echo.
echo 🐳 Building and starting Docker services...
echo This may take a few minutes on first run...
echo.

docker-compose up --build

echo.
echo ✅ Setup complete!
echo.
echo 🎯 Access the application:
echo    Frontend: http://localhost:3000
echo    Backend API: http://localhost:8080
echo    API Docs: http://localhost:8080/health
echo.
echo To stop the services, run: docker-compose down
echo.
pause
