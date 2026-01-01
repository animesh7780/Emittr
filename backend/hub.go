package main

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Hub struct {
	clients      map[*Client]bool
	broadcast    chan interface{}
	register     chan *Client
	unregister   chan *Client
	games        map[string]*GameState
	gameManager  *GameManager
	matchmaking  map[string]*MatchmakeRequest
	mu           sync.RWMutex
}

type Client struct {
	hub       *Hub
	conn      *WSConnection
	send      chan interface{}
	username  string
	gameID    string
	closedAt  time.Time
}

type MatchmakeRequest struct {
	Username  string
	Timestamp time.Time
	Client    *Client
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type GameMoveMessage struct {
	Column int `json:"column"`
}

type GameStartMessage struct {
	GameID   string `json:"gameId"`
	Player1  string `json:"player1"`
	Player2  string `json:"player2"`
	IsBot    bool   `json:"isBot"`
	YourTurn bool   `json:"yourTurn"`
}

type GameMoveEventMessage struct {
	GameID   string `json:"gameId"`
	Column   int    `json:"column"`
	Row      int    `json:"row"`
	Player   int    `json:"player"`
	Board    [ROWS][COLS]int `json:"board"`
}

type GameResultMessage struct {
	GameID string `json:"gameId"`
	Winner string `json:"winner"` // "player1", "player2", "draw"
	WinRow int    `json:"winRow,omitempty"`
	WinCol int    `json:"winCol,omitempty"`
}

func NewHub(gameManager *GameManager) *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		broadcast:   make(chan interface{}, 256),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		games:       make(map[string]*GameState),
		gameManager: gameManager,
		matchmaking: make(map[string]*MatchmakeRequest),
	}
}

func (h *Hub) Run() {
	matchmakingTicker := time.NewTicker(1 * time.Second)
	defer matchmakingTicker.Stop()

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client registered: %s\n", client.username)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("Client unregistered: %s\n", client.username)

			// Handle player disconnection
			if client.gameID != "" {
				h.handlePlayerDisconnect(client)
			}

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Client's send channel is full, skip
				}
			}
			h.mu.RUnlock()

		case <-matchmakingTicker.C:
			h.processMatchmaking()
		}
	}
}

func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

func (h *Hub) BroadcastMessage(msg interface{}) {
	h.broadcast <- msg
}

func (h *Hub) RequestMatchmaking(username string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.matchmaking[username] = &MatchmakeRequest{
		Username:  username,
		Timestamp: time.Now(),
		Client:    client,
	}

	log.Printf("Matchmaking request from %s\n", username)
}

func (h *Hub) processMatchmaking() {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := time.Now()
	var pendingRequests []*MatchmakeRequest

	// Collect all pending requests
	for username, req := range h.matchmaking {
		if now.Sub(req.Timestamp) < 10*time.Second {
			pendingRequests = append(pendingRequests, req)
		} else {
			// Timeout - pair with bot
			if req.Client != nil && req.Client.send != nil {
				h.createGameWithBot(username, req.Client)
			}
			delete(h.matchmaking, username)
		}
	}

	// Try to pair players
	for i := 0; i < len(pendingRequests)-1; i++ {
		for j := i + 1; j < len(pendingRequests); j++ {
			req1 := pendingRequests[i]
			req2 := pendingRequests[j]

			if req1.Client != nil && req2.Client != nil {
				h.createGame(req1.Username, req1.Client, req2.Username, req2.Client)
				delete(h.matchmaking, req1.Username)
				delete(h.matchmaking, req2.Username)
				
				// Remove from pending list
				pendingRequests = append(pendingRequests[:i], pendingRequests[i+1:]...)
				j--
				break
			}
		}
	}
}

func (h *Hub) createGame(username1 string, client1 *Client, username2 string, client2 *Client) {
	gameID := uuid.New().String()
	
	gameState := &GameState{
		ID:            gameID,
		Board:         NewBoard(),
		Player1:       username1,
		Player2:       username2,
		CurrentPlayer: PLAYER1,
		Status:        "active",
		Winner:        "",
		IsBot:         false,
		CreatedAt:     time.Now().Format(time.RFC3339),
	}

	h.games[gameID] = gameState
	client1.gameID = gameID
	client2.gameID = gameID

	// Notify both players
	startMsg := &Message{
		Type: "game_start",
		Payload: GameStartMessage{
			GameID:   gameID,
			Player1:  username1,
			Player2:  username2,
			IsBot:    false,
			YourTurn: true, // Player1 goes first
		},
	}

	client1.send <- startMsg

	startMsg.Payload = GameStartMessage{
		GameID:   gameID,
		Player1:  username1,
		Player2:  username2,
		IsBot:    false,
		YourTurn: false, // Player2 goes second
	}

	client2.send <- startMsg

	log.Printf("Game created: %s between %s and %s\n", gameID, username1, username2)
}

func (h *Hub) createGameWithBot(username string, client *Client) {
	gameID := uuid.New().String()
	
	gameState := &GameState{
		ID:            gameID,
		Board:         NewBoard(),
		Player1:       username,
		Player2:       "Bot",
		CurrentPlayer: PLAYER1,
		Status:        "active",
		Winner:        "",
		IsBot:         true,
		CreatedAt:     time.Now().Format(time.RFC3339),
	}

	h.games[gameID] = gameState
	client.gameID = gameID

	// Notify player
	startMsg := &Message{
		Type: "game_start",
		Payload: GameStartMessage{
			GameID:   gameID,
			Player1:  username,
			Player2:  "Bot",
			IsBot:    true,
			YourTurn: true,
		},
	}

	client.send <- startMsg

	log.Printf("Game created with bot: %s for %s\n", gameID, username)
}

func (h *Hub) HandleGameMove(client *Client, column int) {
	h.mu.Lock()
	gameState := h.games[client.gameID]
	h.mu.Unlock()

	if gameState == nil {
		client.send <- &Message{
			Type:    "error",
			Payload: map[string]string{"message": "Game not found"},
		}
		return
	}

	// Determine which player made the move
	var player int
	if gameState.Player1 == client.username {
		player = PLAYER1
	} else {
		player = PLAYER2
	}

	// Validate it's the player's turn
	if gameState.CurrentPlayer != player {
		client.send <- &Message{
			Type:    "error",
			Payload: map[string]string{"message": "Not your turn"},
		}
		return
	}

	// Make the move
	row, err := gameState.Board.DropDisc(column, player)
	if err != nil {
		client.send <- &Message{
			Type:    "error",
			Payload: map[string]string{"message": err.Error()},
		}
		return
	}

	// Broadcast the move
	moveMsg := &Message{
		Type: "game_move",
		Payload: GameMoveEventMessage{
			GameID: client.gameID,
			Column: column,
			Row:    row,
			Player: player,
			Board:  gameState.Board.Grid,
		},
	}

	h.broadcastToGame(client.gameID, moveMsg)

	// Check for win
	if gameState.Board.CheckWin(row, column, player) {
		gameState.Status = "finished"
		gameState.Winner = client.username
		gameState.UpdatedAt = time.Now().Format(time.RFC3339)

		resultMsg := &Message{
			Type: "game_result",
			Payload: GameResultMessage{
				GameID: client.gameID,
				Winner: client.username,
				WinRow: row,
				WinCol: column,
			},
		}

		h.broadcastToGame(client.gameID, resultMsg)
		h.gameManager.SaveGame(gameState)
		return
	}

	// Check for draw
	if gameState.Board.IsBoardFull() {
		gameState.Status = "finished"
		gameState.Winner = "draw"
		gameState.UpdatedAt = time.Now().Format(time.RFC3339)

		resultMsg := &Message{
			Type: "game_result",
			Payload: GameResultMessage{
				GameID: client.gameID,
				Winner: "draw",
			},
		}

		h.broadcastToGame(client.gameID, resultMsg)
		h.gameManager.SaveGame(gameState)
		return
	}

	// Switch turn
	if gameState.CurrentPlayer == PLAYER1 {
		gameState.CurrentPlayer = PLAYER2
	} else {
		gameState.CurrentPlayer = PLAYER1
	}

	// If opponent is bot, make bot move
	if gameState.IsBot && gameState.CurrentPlayer == PLAYER2 {
		time.Sleep(1 * time.Second) // Simulate thinking time
		h.makeBotMove(gameState, client)
	}
}

func (h *Hub) makeBotMove(gameState *GameState, playerClient *Client) {
	bot := NewBot("medium")
	column := bot.GetBotMove(gameState.Board, PLAYER2, PLAYER1)

	if column == -1 {
		// No valid moves (shouldn't happen)
		return
	}

	row, err := gameState.Board.DropDisc(column, PLAYER2)
	if err != nil {
		return
	}

	// Broadcast the move
	moveMsg := &Message{
		Type: "game_move",
		Payload: GameMoveEventMessage{
			GameID: gameState.ID,
			Column: column,
			Row:    row,
			Player: PLAYER2,
			Board:  gameState.Board.Grid,
		},
	}

	h.broadcastToGame(gameState.ID, moveMsg)

	// Check for win
	if gameState.Board.CheckWin(row, column, PLAYER2) {
		gameState.Status = "finished"
		gameState.Winner = "Bot"
		gameState.UpdatedAt = time.Now().Format(time.RFC3339)

		resultMsg := &Message{
			Type: "game_result",
			Payload: GameResultMessage{
				GameID: gameState.ID,
				Winner: "Bot",
				WinRow: row,
				WinCol: column,
			},
		}

		h.broadcastToGame(gameState.ID, resultMsg)
		h.gameManager.SaveGame(gameState)
		return
	}

	// Check for draw
	if gameState.Board.IsBoardFull() {
		gameState.Status = "finished"
		gameState.Winner = "draw"
		gameState.UpdatedAt = time.Now().Format(time.RFC3339)

		resultMsg := &Message{
			Type: "game_result",
			Payload: GameResultMessage{
				GameID: gameState.ID,
				Winner: "draw",
			},
		}

		h.broadcastToGame(gameState.ID, resultMsg)
		h.gameManager.SaveGame(gameState)
		return
	}

	gameState.CurrentPlayer = PLAYER1
}

func (h *Hub) handlePlayerDisconnect(client *Client) {
	h.mu.Lock()
	gameState := h.games[client.gameID]
	h.mu.Unlock()

	if gameState == nil || gameState.Status == "finished" {
		return
	}

	// Mark disconnection time
	client.closedAt = time.Now()

	// Wait 30 seconds for reconnection
	go func() {
		time.Sleep(30 * time.Second)

		h.mu.Lock()
		defer h.mu.Unlock()

		gameState := h.games[client.gameID]
		if gameState == nil || gameState.Status == "finished" {
			return
		}

		// Still disconnected - forfeit
		gameState.Status = "finished"
		gameState.UpdatedAt = time.Now().Format(time.RFC3339)

		if gameState.IsBot {
			gameState.Winner = "Bot"
		} else {
			if gameState.Player1 == client.username {
				gameState.Winner = gameState.Player2
			} else {
				gameState.Winner = gameState.Player1
			}
		}

		h.gameManager.SaveGame(gameState)
		log.Printf("Game %s forfeited due to player disconnect\n", client.gameID)
	}()
}

func (h *Hub) broadcastToGame(gameID string, msg interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.gameID == gameID {
			select {
			case client.send <- msg:
			default:
				// Client's send channel is full
			}
		}
	}
}
