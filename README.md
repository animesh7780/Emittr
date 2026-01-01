# 🎮 4 in a Row - Backend Engineering Assignment

A real-time, multiplayer **4 in a Row** (Connect Four) game with WebSocket support, competitive bot AI, and Kafka-based analytics. Built with **GoLang** backend, **React** frontend, and **PostgreSQL** database.

## 📋 Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Development Setup](#development-setup)
- [API Documentation](#api-documentation)
- [Game Rules](#game-rules)
- [Deployment](#deployment)
- [Project Structure](#project-structure)

## ✨ Features

### Core Gameplay
- ✅ **7×6 game board** with real-time updates
- ✅ **Player matchmaking** with 10-second timeout for bot fallback
- ✅ **Competitive bot AI** that:
  - Blocks opponent's winning moves
  - Tries to win when possible
  - Makes strategic placements
- ✅ **Player reconnection** within 30 seconds
- ✅ **Win detection** (horizontal, vertical, diagonal)
- ✅ **Draw detection** when board is full

### Real-time Features
- ✅ **WebSocket support** for live gameplay
- ✅ **Instant board updates** across both players
- ✅ **Automatic bot moves** with 1-second delay (for realism)
- ✅ **Game session persistence**

### Analytics & Leaderboard
- ✅ **Kafka event streaming** for game analytics
- ✅ **Player leaderboard** showing:
  - Total wins
  - Total losses
  - Total draws
  - Win rate percentage
- ✅ **Game duration tracking**
- ✅ **Player statistics** API

### Infrastructure
- ✅ **PostgreSQL database** for persistent storage
- ✅ **Docker & Docker Compose** for containerization
- ✅ **CORS enabled** for cross-origin requests
- ✅ **Health check endpoints**

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────┐
│                   Frontend (React)                       │
│            - Game Board UI (7x6 grid)                   │
│            - Lobby & Matchmaking                        │
│            - Leaderboard Display                        │
│            - Hosting: Nginx                             │
└─────────────────────────────────────────────────────────┘
                           ↕ (WebSocket)
┌─────────────────────────────────────────────────────────┐
│              Backend (GoLang) - Core Services           │
│  ┌──────────────────────────────────────────────────┐  │
│  │ Gin Web Framework + Gorilla WebSocket            │  │
│  │ - Player Matchmaking (10s wait → Bot)           │  │
│  │ - Game Management & State                        │  │
│  │ - Turn Management & Win Detection               │  │
│  │ - Player Reconnection (30s window)              │  │
│  └──────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────┐  │
│  │ Competitive Bot AI                               │  │
│  │ - Evaluates board positions                      │  │
│  │ - Prioritizes blocking opponent                 │  │
│  │ - Finds winning opportunities                    │  │
│  │ - Strategic center-heavy placement              │  │
│  └──────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────┐  │
│  │ Kafka Producer                                    │  │
│  │ - Game completion events                         │  │
│  │ - Move events                                    │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
         ↓                              ↓
┌─────────────────┐        ┌────────────────────┐
│   PostgreSQL    │        │  Kafka Consumer    │
│                 │        │  (Analytics)       │
│ Tables:         │        │                    │
│ - players       │        │ Metrics:           │
│ - games         │        │ - Avg game time    │
│ - leaderboard   │        │ - Win rates        │
│   (views)       │        │ - Games/hour       │
└─────────────────┘        └────────────────────┘
```

### Data Flow
1. **Player Registration** → Matchmaking Queue
2. **Game Pairing** → Create Game Instance
3. **Player Moves** → Validate → Update Board → Broadcast
4. **Game Events** → Kafka Producer → Event Stream
5. **Analytics Consumer** → Process Events → Store Metrics

## 📦 Prerequisites

- **Docker** & **Docker Compose** (recommended for easy setup)
- OR manually:
  - **GoLang 1.21+**
  - **Node.js 18+**
  - **PostgreSQL 14+**
  - **Kafka 3.5+**

## 🚀 Quick Start (Docker Compose)

### 1. Clone and Navigate
```bash
cd d:\College\Company-Assignments\Emittr
```

### 2. Start All Services
```bash
docker-compose up --build
```

### 3. Access the Application
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **WebSocket**: ws://localhost:8080/ws

### 4. Default Connection Strings
- **PostgreSQL**: `postgres://user:password@localhost:5432/4_in_a_row`
- **Kafka**: `localhost:9092`

### 5. Stop Services
```bash
docker-compose down
```

## 🛠️ Development Setup

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

## 📡 API Documentation

### WebSocket Messages

#### Client → Server

**Register/Matchmaking:**
```json
{
  "type": "register",
  "payload": {
    "username": "player_name"
  }
}
```

**Make a Move:**
```json
{
  "type": "game_move",
  "payload": {
    "column": 3
  }
}
```

**Rejoin Game:**
```json
{
  "type": "rejoin",
  "payload": {
    "gameId": "game-uuid"
  }
}
```

#### Server → Client

**Game Start:**
```json
{
  "type": "game_start",
  "payload": {
    "gameId": "game-uuid",
    "player1": "alice",
    "player2": "bob or Bot",
    "isBot": false,
    "yourTurn": true
  }
}
```

**Game Move:**
```json
{
  "type": "game_move",
  "payload": {
    "gameId": "game-uuid",
    "column": 3,
    "row": 5,
    "player": 1,
    "board": [[0,0,...], ...]
  }
}
```

**Game Result:**
```json
{
  "type": "game_result",
  "payload": {
    "gameId": "game-uuid",
    "winner": "alice or draw",
    "winRow": 4,
    "winCol": 3
  }
}
```

### REST API Examples

**Get Leaderboard:**
```bash
curl http://localhost:8080/api/leaderboard
```

Response:
```json
{
  "leaderboard": [
    {
      "username": "alice",
      "wins": 10,
      "losses": 2,
      "draws": 1,
      "winRate": "76.92%"
    }
  ]
}
```

**Get Player Stats:**
```bash
curl http://localhost:8080/api/player/alice
```

Response:
```json
{
  "username": "alice",
  "wins": 10,
  "losses": 2,
  "draws": 1,
  "winRate": "76.92%",
  "createdAt": "2024-01-15T10:30:00Z"
}
```

## 🎮 Game Rules

### Board
- **7 columns** × **6 rows**
- Discs fall to the lowest empty space

### Winning
- Connect **4 discs** in a row:
  - **Horizontal** ↔️
  - **Vertical** ↕️
  - **Diagonal** ↖️↘️ or ↗️↙️

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

## 🐳 Deployment

### Heroku Deployment

1. **Setup Heroku PostgreSQL and Kafka (CloudKarafka):**
```bash
heroku create your-app-name
heroku addons:create heroku-postgresql:standard-0
heroku addons:create cloudkarafka:giraffe
```

2. **Set Environment Variables:**
```bash
heroku config:set DATABASE_URL=postgres://...
heroku config:set KAFKA_BROKER=kafka-broker.cloudkarafka.com:9092
heroku config:set PORT=8080
```

3. **Deploy Backend:**
```bash
cd backend
git push heroku main
```

4. **Deploy Frontend to Vercel/Netlify:**
```bash
cd frontend
npm run build
# Push to Vercel/Netlify
```

### Railway Deployment

1. Connect GitHub repository
2. Add PostgreSQL database service
3. Add Kafka service (via Redis or Docker image)
4. Set environment variables in dashboard
5. Deploy automatically on push

## 📁 Project Structure

```
4-in-a-row/
├── backend/
│   ├── main.go              # Application entry point
│   ├── game.go              # Core game logic (Board, moves, win detection)
│   ├── bot.go               # Competitive bot AI
│   ├── hub.go               # WebSocket hub & matchmaking
│   ├── websocket.go         # WebSocket connection handler
│   ├── server.go            # HTTP server & routes
│   ├── database.go          # PostgreSQL models & queries
│   ├── kafka.go             # Kafka producer & event handling
│   ├── go.mod               # Go module file
│   ├── .env.example         # Environment template
│   └── Dockerfile           # Container configuration
│
├── frontend/
│   ├── src/
│   │   ├── main.jsx         # React entry point
│   │   ├── App.jsx          # Main app component
│   │   ├── App.css          # Global styles
│   │   ├── index.css        # CSS variables & resets
│   │   └── components/
│   │       ├── Lobby.jsx    # Login/matchmaking screen
│   │       ├── GameBoard.jsx        # Game board & controls
│   │       ├── GameResult.jsx       # Result screen
│   │       ├── Leaderboard.jsx      # Leaderboard display
│   │       └── *.css        # Component styles
│   ├── index.html           # HTML entry
│   ├── package.json         # NPM dependencies
│   ├── vite.config.js       # Vite configuration
│   ├── nginx.conf           # Nginx config for production
│   └── Dockerfile           # Container configuration
│
├── analytics/
│   ├── main.go              # Kafka consumer service
│   ├── go.mod               # Go module file
│   ├── .env.example         # Environment template
│   └── Dockerfile           # Container configuration
│
├── docker-compose.yml       # Orchestration file
└── README.md               # This file
```

## 🤖 Bot AI Strategy

The bot uses a **minimax-inspired evaluation function**:

1. **Immediate Win** (Score: 10,000) - Take winning move
2. **Blocking** (Score: 9,000) - Block opponent's win
3. **Position Score** - Based on:
   - Center column preference (more valuable)
   - Adjacent disc count (threats & opportunities)
4. **Fallback** - Random valid move if no strategic move

## 📊 Analytics Events

Events published to Kafka for analysis:

```json
{
  "eventType": "game_completed",
  "gameId": "uuid",
  "player": "alice",
  "opponent": "bob",
  "timestamp": "2024-01-15T10:30:00Z",
  "isBot": false,
  "gameResult": "alice",
  "duration": 180
}
```

**Metrics Tracked:**
- Average game duration
- Win rates per player
- Games played per hour/day
- Bot vs player statistics
- Most frequent winners

## 🐛 Troubleshooting

### PostgreSQL Connection Error
```
DATABASE_URL must be set and valid
```
**Solution**: Check `.env` file and ensure PostgreSQL is running

### Kafka Connection Timeout
```
Failed to connect to Kafka broker
```
**Solution**: Ensure Kafka is running on configured host:port

### WebSocket Connection Refused
```
WebSocket connection failed
```
**Solution**: Ensure backend is running and CORS is enabled

### Port Already in Use
```
Error: listen EADDRINUSE: address already in use :::8080
```
**Solution**: Change PORT in `.env` or kill process using port

## 📝 Environment Variables

### Backend (.env)
```
DATABASE_URL=postgres://user:password@localhost:5432/4_in_a_row
KAFKA_BROKER=localhost:9092
KAFKA_TOPIC=game_events
PORT=8080
ENVIRONMENT=development
```

### Analytics (.env)
```
DATABASE_URL=postgres://user:password@localhost:5432/4_in_a_row
KAFKA_BROKER=localhost:9092
KAFKA_TOPIC=game_events
KAFKA_GROUP=analytics_group
```

## 🔐 Security Considerations

- ✅ CORS enabled (adjust for production)
- ✅ WebSocket origins validated
- ✅ SQL injection protection (parameterized queries)
- ✅ Input validation on game moves
- ⚠️ TODO: Add authentication (JWT tokens)
- ⚠️ TODO: Add rate limiting
- ⚠️ TODO: Add HTTPS/WSS in production

## 📈 Future Enhancements

- [ ] User authentication & registration
- [ ] Player ratings/ELO system
- [ ] Chat during gameplay
- [ ] Game replay system
- [ ] Mobile app (React Native)
- [ ] Advanced bot difficulty levels
- [ ] Tournament mode
- [ ] Achievement system
- [ ] Real-time notifications

## 📄 License

MIT License - Feel free to use this project for learning purposes

## 👥 Author

Backend Engineering Intern Assignment - 4 in a Row Game

---

**Happy Gaming! 🎮**

For questions or issues, please refer to the architecture documentation or check the code comments.
