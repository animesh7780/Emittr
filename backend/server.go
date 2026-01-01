package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port        string
	hub         *Hub
	gameManager *GameManager
	db          *Database
	router      *gin.Engine
}

func NewServer(port string, hub *Hub, gameManager *GameManager, db *Database) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	server := &Server{
		port:        port,
		hub:         hub,
		gameManager: gameManager,
		db:          db,
		router:      router,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// CORS middleware
	s.router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// WebSocket endpoint
	s.router.GET("/ws", func(c *gin.Context) {
		HandleWebSocket(s.hub)(c.Writer, c.Request)
	})

	// API endpoints
	s.router.GET("/api/leaderboard", s.getLeaderboard)
	s.router.GET("/api/player/:username", s.getPlayerStats)
	s.router.GET("/api/game/:gameId", s.getGameState)
	s.router.GET("/health", s.health)
}

func (s *Server) getLeaderboard(c *gin.Context) {
	leaderboard, err := s.db.GetLeaderboard(100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"leaderboard": leaderboard,
	})
}

func (s *Server) getPlayerStats(c *gin.Context) {
	username := c.Param("username")
	
	stats, err := s.db.GetPlayerStats(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (s *Server) getGameState(c *gin.Context) {
	gameID := c.Param("gameId")
	
	s.hub.mu.RLock()
	gameState := s.hub.games[gameID]
	s.hub.mu.RUnlock()

	if gameState == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	c.JSON(http.StatusOK, gameState)
}

func (s *Server) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

func (s *Server) Start() error {
	log.Printf("Starting server on port %s\n", s.port)
	return s.router.Run(":" + s.port)
}
