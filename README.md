# 4 in a Row — Real-Time Multiplayer

## Live Links

- Play the game: https://emittr.onrender.com
- REST API base: https://emittr-backend-0l70.onrender.com (health: https://emittr-backend-0l70.onrender.com/health)
- WebSocket endpoint: wss://emittr-backend-0l70.onrender.com/ws
- Note: Render cold-start may add a few seconds on the first request.

### Access (local)
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- WebSocket: ws://localhost:8080/ws

### Access (hosted)
- Frontend: https://emittr.onrender.com
- Backend API: https://emittr-backend-0l70.onrender.com
- WebSocket: wss://emittr-backend-0l70.onrender.com/ws

### Stop
```bash
docker-compose down
```

### Manual dev setup (if not using Docker)
```bash
# Backend
cd backend
go mod download
cp .env.example .env
go run main.go bot.go game.go hub.go websocket.go database.go kafka.go server.go

# Frontend
cd ../frontend
npm install
npm run dev
```

## Requirements

- **Docker** & **Docker Compose** (easy setup)
- OR manually:
  - **GoLang 1.21+**
  - **Node.js 18+**
  - **PostgreSQL 14+**
  - **Kafka 3.5+**

## Quick Start (Docker Compose)

If you just want to see it running locally, use Docker Compose. For a zero-setup experience, jump to the live link above.

### 1. Clone and Navigate
```bash
cd d:\College\Company-Assignments\Emittr
```

### 2. Start All Services
```bash
docker-compose up --build
```

### 3. Access the Application
- Local frontend: http://localhost:3000
- Local backend API: http://localhost:8080
- Local WebSocket: ws://localhost:8080/ws
- Hosted frontend: https://emittr.onrender.com
- Hosted backend API: https://emittr-backend-0l70.onrender.com

### 4. Default Connection Strings
- **PostgreSQL**: `postgres://user:password@localhost:5432/4_in_a_row`
- **Kafka**: `localhost:9092`

### 5. Stop Services
```bash
docker-compose down
```

## Development Setup

### Backend Setup

```bash
cd backend

# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Edit .env with your database credentials
# Start PostgreSQL and Kafka first

# Run backend
go run main.go bot.go game.go hub.go websocket.go database.go kafka.go server.go
```

**Backend API Endpoints:**
- `GET /health` - Health check
- `GET /api/leaderboard` - Get top 100 players
- `GET /api/player/:username` - Get player stats
- `GET /api/game/:gameId` - Get game state
- `WS /ws` - WebSocket connection

### Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

**Frontend will be available at**: http://localhost:3000

**Build for production:**
```bash
npm run build
```

### Analytics Service Setup

```bash
cd analytics

# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Run analytics consumer
go run main.go
```

The analytics service will:
- Connect to Kafka
- Consume game events
- Log metrics and analytics
- Track game duration, win rates, player performance

## Game Rules

### Board
- **7 columns** x **6 rows**
- Discs fall to the lowest empty space

### Winning
- Connect **4 discs** in a row:
  - **Horizontal** (left to right)
  - **Vertical** (top to bottom)
  - **Diagonal** (both slopes)

### Game Flow
1. Player 1 registers username
2. If no opponent within 10 seconds → Play vs Bot
3. If opponent found → Start PvP game
4. Players alternate turns
5. First to 4 in a row wins
6. If board fills up → Draw

### Disconnection
- Player can rejoin within **30 seconds** using their game ID
- After 30 seconds → Opponent wins by default

