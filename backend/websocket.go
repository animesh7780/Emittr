package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type WSConnection struct {
	conn *websocket.Conn
}

func NewWSConnection(w http.ResponseWriter, r *http.Request) (*WSConnection, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return &WSConnection{conn: conn}, nil
}

func (wsc *WSConnection) ReadMessage() (string, []byte, error) {
	_, data, err := wsc.conn.ReadMessage()
	if err != nil {
		return "", nil, err
	}

	var msg Message
	err = json.Unmarshal(data, &msg)
	if err != nil {
		return "", nil, err
	}

	payload, _ := json.Marshal(msg.Payload)
	return msg.Type, payload, nil
}

func (wsc *WSConnection) WriteMessage(msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return wsc.conn.WriteMessage(websocket.TextMessage, data)
}

func (wsc *WSConnection) Close() error {
	return wsc.conn.Close()
}

func (wsc *WSConnection) SetReadDeadline(t time.Time) error {
	return wsc.conn.SetReadDeadline(t)
}

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wsConn, err := NewWSConnection(w, r)
		if err != nil {
			log.Println("WebSocket upgrade error:", err)
			return
		}
		defer wsConn.Close()

		client := &Client{
			hub:      hub,
			conn:     wsConn,
			send:     make(chan interface{}, 256),
			username: "",
			gameID:   "",
		}

		hub.RegisterClient(client)
		defer hub.UnregisterClient(client)

		// Send goroutine
		go func() {
			for msg := range client.send {
				if err := wsConn.WriteMessage(msg); err != nil {
					log.Println("Write error:", err)
					return
				}
			}
		}()

		// Read goroutine
		for {
			msgType, payload, err := wsConn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}

			log.Printf("Received message type: %s from %s\n", msgType, client.username)

			switch msgType {
			case "register":
				var registerMsg struct {
					Username string `json:"username"`
				}
				json.Unmarshal(payload, &registerMsg)
				client.username = registerMsg.Username
				hub.RequestMatchmaking(registerMsg.Username, client)
				log.Printf("Player registered: %s\n", registerMsg.Username)

			case "game_move":
				var moveMsg GameMoveMessage
				json.Unmarshal(payload, &moveMsg)
				hub.HandleGameMove(client, moveMsg.Column)

			case "rejoin":
				var rejoinMsg struct {
					GameID string `json:"gameId"`
				}
				json.Unmarshal(payload, &rejoinMsg)
				client.gameID = rejoinMsg.GameID
				log.Printf("Player %s rejoining game %s\n", client.username, rejoinMsg.GameID)

			default:
				log.Printf("Unknown message type: %s\n", msgType)
			}
		}
	}
}
