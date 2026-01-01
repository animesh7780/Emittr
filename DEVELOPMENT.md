# Local Development Instructions

## Prerequisites
- Docker & Docker Compose installed
- Or: Go 1.21+, Node 18+, PostgreSQL 14+, Kafka 3.5+

## Using Docker Compose (Recommended)

```bash
# 1. Navigate to project root
cd d:\College\Company-Assignments\Emittr

# 2. Run setup script (Windows)
setup.bat

# Or on Mac/Linux:
chmod +x setup.sh
./setup.sh

# This will:
# - Create .env files
# - Build all Docker images
# - Start all services

# 3. Access the application
# Frontend: http://localhost:3000
# Backend: http://localhost:8080
# Swagger/Docs: http://localhost:8080/health

# 4. To stop services
docker-compose down

# 5. To see logs
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f analytics
```

## Manual Setup (Without Docker)

### 1. PostgreSQL Setup
```bash
# Create database
createdb 4_in_a_row

# PostgreSQL should be running on localhost:5432
```

### 2. Kafka Setup
```bash
# Download and extract Kafka
# Start Zookeeper
bin/zookeeper-server-start.sh config/zookeeper.properties

# Start Kafka (in another terminal)
bin/kafka-server-start.sh config/server.properties

# Create topic
bin/kafka-topics.sh --create --topic game_events --bootstrap-server localhost:9092
```

### 3. Backend Setup
```bash
cd backend

# Create .env file
cp .env.example .env

# Edit .env if needed
# nano .env

# Install and run
go mod download
go run main.go bot.go game.go hub.go websocket.go database.go kafka.go server.go

# Server will start on port 8080
```

### 4. Frontend Setup
```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev

# Frontend will be available at http://localhost:3000
```

### 5. Analytics Service Setup
```bash
cd analytics

# Create .env file
cp .env.example .env

# Install and run
go mod download
go run main.go

# Will start consuming from Kafka topic
```

## Testing the Game

1. Open browser and go to http://localhost:3000
2. Enter username and click "Play Now"
3. Wait for opponent or bot (10 seconds)
4. Game will start automatically
5. Click columns to drop discs
6. First to 4 in a row wins!

## Common Issues

### Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```

### PostgreSQL Connection Failed
- Ensure PostgreSQL is running
- Check DATABASE_URL in .env
- Try connecting manually: `psql -U user -d 4_in_a_row -h localhost`

### Kafka Connection Failed
- Ensure Kafka and Zookeeper are running
- Check KAFKA_BROKER in .env
- Verify topic exists: `bin/kafka-topics.sh --list --bootstrap-server localhost:9092`

### Node Modules Issues
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

## Database Schema

Tables created automatically on startup:

### players
- id: INTEGER PRIMARY KEY
- username: VARCHAR(255) UNIQUE
- wins: INTEGER DEFAULT 0
- losses: INTEGER DEFAULT 0
- draws: INTEGER DEFAULT 0
- created_at: TIMESTAMP
- updated_at: TIMESTAMP

### games
- id: VARCHAR(36) PRIMARY KEY
- player1: VARCHAR(255)
- player2: VARCHAR(255)
- winner: VARCHAR(255)
- is_bot: BOOLEAN
- status: VARCHAR(50)
- board_state: JSONB
- created_at: TIMESTAMP
- updated_at: TIMESTAMP
- duration_seconds: INTEGER

## API Testing

```bash
# Get leaderboard
curl http://localhost:8080/api/leaderboard

# Get player stats
curl http://localhost:8080/api/player/username

# Get game state
curl http://localhost:8080/api/game/game-id

# Health check
curl http://localhost:8080/health
```

## WebSocket Testing

Use WebSocket client or browser console:

```javascript
// Connect
const ws = new WebSocket('ws://localhost:8080/ws');

// Register
ws.send(JSON.stringify({
  type: 'register',
  payload: { username: 'testplayer' }
}));

// Make move
ws.send(JSON.stringify({
  type: 'game_move',
  payload: { column: 3 }
}));
```

## Viewing Logs

```bash
# Backend logs
docker-compose logs -f backend

# Frontend logs
docker-compose logs -f frontend

# Analytics logs
docker-compose logs -f analytics

# Postgres logs
docker-compose logs -f postgres

# All services
docker-compose logs -f
```

## Building for Production

### Frontend
```bash
cd frontend
npm run build
# Creates dist/ folder for deployment
```

### Backend
```bash
cd backend
GOOS=linux GOARCH=amd64 go build -o backend
# Creates binary for Linux deployment
```

### Docker Images
```bash
docker-compose build

# Push to registry
docker tag 4-in-a-row-backend:latest your-registry/4-in-a-row-backend:latest
docker push your-registry/4-in-a-row-backend:latest
```
