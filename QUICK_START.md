# 🚀 Quick Reference Guide

## ⚡ Start Playing in 5 Minutes

### Using Docker Compose (Windows)
```bash
cd d:\College\Company-Assignments\Emittr
setup.bat
# Wait 2-3 minutes for services to start
# Open browser: http://localhost:3000
```

### Using Docker Compose (Mac/Linux)
```bash
cd /path/to/Emittr
./setup.sh
# Wait 2-3 minutes for services to start
# Open browser: http://localhost:3000
```

---

## 📝 Game Rules (30 seconds)

| Aspect | Details |
|--------|---------|
| Board | 7 columns × 6 rows |
| Goal | Connect 4 discs in a row |
| Directions | Horizontal, vertical, diagonal |
| Turns | Players alternate |
| Win | First to 4 in a row wins |
| Draw | Board full with no winner |
| Opponent | Player or Bot (after 10 sec wait) |

---

## 🌐 Access Points

| Service | URL | Port |
|---------|-----|------|
| Frontend | http://localhost:3000 | 3000 |
| Backend API | http://localhost:8080 | 8080 |
| WebSocket | ws://localhost:8080/ws | 8080 |
| PostgreSQL | localhost:5432 | 5432 |
| Kafka | localhost:9092 | 9092 |

---

## 🎮 How to Play

1. **Go to** http://localhost:3000
2. **Enter** your username
3. **Click** "Play Now"
4. **Wait** for opponent (10 seconds) or play bot
5. **Click** column ↓ buttons to drop discs
6. **Connect** 4 discs to win!
7. **View** leaderboard anytime

---

## 💻 Key Commands

### Start Services
```bash
docker-compose up --build
```

### Stop Services
```bash
docker-compose down
```

### View Logs
```bash
docker-compose logs -f backend      # Backend logs
docker-compose logs -f frontend     # Frontend logs
docker-compose logs -f analytics    # Analytics logs
```

### Database Access
```bash
# Connect to PostgreSQL
psql -U user -d 4_in_a_row -h localhost

# Inside psql:
SELECT * FROM players;              # View players
SELECT * FROM games;                # View games
\dt                                 # List tables
\q                                  # Exit
```

### Clear Everything
```bash
docker-compose down -v              # Remove all data
rm -rf node_modules                 # Clean frontend
```

---

## 🔧 Configuration

### Backend Environment (.env)
```
DATABASE_URL=postgres://user:password@postgres:5432/4_in_a_row
KAFKA_BROKER=kafka:9092
KAFKA_TOPIC=game_events
PORT=8080
```

### Frontend Environment (.env)
```
VITE_API_URL=http://localhost:8080
```

### Analytics Environment (.env)
```
DATABASE_URL=postgres://user:password@postgres:5432/4_in_a_row
KAFKA_BROKER=kafka:9092
KAFKA_TOPIC=game_events
KAFKA_GROUP=analytics_group
```

---

## 📱 API Quick Reference

### Leaderboard
```bash
curl http://localhost:8080/api/leaderboard
```

### Player Stats
```bash
curl http://localhost:8080/api/player/username
```

### Game State
```bash
curl http://localhost:8080/api/game/game-id
```

### Health Check
```bash
curl http://localhost:8080/health
```

---

## 🐛 Common Issues & Fixes

| Issue | Solution |
|-------|----------|
| Port already in use | Change port in .env or kill process |
| Database connection failed | Ensure PostgreSQL is running |
| Kafka connection failed | Check Kafka is running |
| Can't connect to frontend | Verify backend is running first |
| Slow bot moves | Normal - intentional 1s delay |
| Game won't load | Check browser console (F12) for errors |

---

## 🎯 Project Files Quick Guide

| File | Purpose |
|------|---------|
| README.md | Complete documentation |
| DEVELOPMENT.md | Local setup guide |
| DEPLOYMENT.md | Deployment instructions |
| PROJECT_SUMMARY.md | Project overview |
| docker-compose.yml | Service orchestration |
| backend/main.go | Backend entry point |
| frontend/src/App.jsx | Frontend entry point |
| analytics/main.go | Analytics entry point |

---

## 📊 Game Flow Diagram

```
┌─────────────┐
│   Lobby     │
│  Enter Name │
└──────┬──────┘
       │
       ▼
┌──────────────────┐
│  Matchmaking     │
│  Wait 10 seconds │
└──────┬───────┬──────┘
       │       │
    Player   Bot
       │       │
       ▼       ▼
┌──────────────────┐
│   Game Start     │
│  Player 1 First  │
└──────┬───────────┘
       │
       ▼
   ┌────────────────────────┐
   │  Player Makes Move     │
   │  Drop Disc in Column   │
   │  Board Updates         │
   └────────┬───────────────┘
            │
     ┌──────┴─────────┐
     │                │
     ▼                ▼
  WIN?              DRAW?
     │                │
   YES              YES
     │                │
     ▼                ▼
 ┌──────┐       ┌──────┐
 │WIN   │       │DRAW  │
 └──────┘       └──────┘
     │                │
     └────┬───────────┘
          │
          ▼
    ┌──────────────┐
    │  Game Over   │
    │  View Result │
    │  Play Again  │
    └──────────────┘
```

---

## 🎓 Learning Resources

- **Go Programming**: https://tour.golang.org
- **React Hooks**: https://react.dev/reference/react/hooks
- **WebSocket**: https://developer.mozilla.org/en-US/docs/Web/API/WebSocket
- **PostgreSQL**: https://www.postgresql.org/docs
- **Kafka**: https://kafka.apache.org/intro
- **Docker**: https://docs.docker.com

---

## 🚀 Deployment Checklist

- [ ] All code pushed to GitHub
- [ ] README.md reviewed
- [ ] docker-compose.yml tested locally
- [ ] Environment variables configured
- [ ] Database schema verified
- [ ] Frontend build tested
- [ ] Backend APIs tested
- [ ] Kafka connection verified
- [ ] Docker images built
- [ ] Choose deployment platform
- [ ] Set up cloud resources
- [ ] Deploy services
- [ ] Test live application
- [ ] Monitor logs

---

## 💡 Pro Tips

1. **WebSocket Debugging**
   - Open DevTools (F12)
   - Go to Network tab
   - Filter by WS (WebSocket)
   - Watch messages in real-time

2. **Database Debugging**
   - Use `psql` command line tool
   - Query tables directly
   - Check player stats in real-time

3. **Log Watching**
   ```bash
   docker-compose logs -f --tail=50
   ```

4. **Clean Rebuild**
   ```bash
   docker-compose down -v
   docker-compose build --no-cache
   docker-compose up
   ```

5. **Test Bot AI**
   - Play multiple games
   - Note bot strategic moves
   - Bot should block your wins

---

## 📞 Quick Help

**Error in logs?**
```bash
docker-compose logs -f backend | grep -i error
```

**Need to reset database?**
```bash
docker-compose down -v
docker-compose up
```

**Want to rebuild one service?**
```bash
docker-compose build backend --no-cache
docker-compose up backend
```

**Check all running containers?**
```bash
docker ps
```

**View resource usage?**
```bash
docker stats
```

---

## 🎉 That's it!

You're ready to:
1. ✅ Play 4 in a Row
2. ✅ Compete with bot or friends
3. ✅ View leaderboards
4. ✅ Deploy to production

**Happy gaming!** 🎮

---

**For more details**, check:
- README.md (full documentation)
- DEVELOPMENT.md (setup guide)
- DEPLOYMENT.md (deployment guide)
