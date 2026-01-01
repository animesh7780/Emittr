package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

type GameEvent struct {
	EventType   string    `json:"eventType"`
	GameID      string    `json:"gameId"`
	Player      string    `json:"player"`
	Opponent    string    `json:"opponent"`
	Action      string    `json:"action"`
	Column      int       `json:"column,omitempty"`
	Row         int       `json:"row,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
	IsBot       bool      `json:"isBot"`
	GameResult  string    `json:"gameResult,omitempty"`
	Duration    int       `json:"duration,omitempty"`
}

func NewKafkaProducer(broker string, topic string) (*KafkaProducer, error) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{broker},
		Topic:        topic,
		RequiredAcks: int(kafka.RequireOne),
	})

	return &KafkaProducer{writer: writer}, nil
}

func (kp *KafkaProducer) PublishEvent(event GameEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling event: %v\n", err)
		return err
	}

	msg := kafka.Message{
		Key:   []byte(event.GameID),
		Value: data,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = kp.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("Error publishing event to Kafka: %v\n", err)
		return err
	}

	return nil
}

func (kp *KafkaProducer) Close() error {
	return kp.writer.Close()
}

// GameManager manages games and publishes events
type GameManager struct {
	db       *Database
	producer *KafkaProducer
}

func NewGameManager(db *Database, producer *KafkaProducer) *GameManager {
	return &GameManager{
		db:       db,
		producer: producer,
	}
}

func (gm *GameManager) SaveGame(game *GameState) error {
	// Save to database
	err := gm.db.SaveGame(game)
	if err != nil {
		log.Printf("Error saving game to database: %v\n", err)
		return err
	}

	createdAt, _ := time.Parse(time.RFC3339, game.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, game.UpdatedAt)
	duration := int(updatedAt.Sub(createdAt).Seconds())

	// Publish event
	event := GameEvent{
		EventType:  "game_completed",
		GameID:     game.ID,
		Player:     game.Player1,
		Opponent:   game.Player2,
		Timestamp:  time.Now(),
		IsBot:      game.IsBot,
		GameResult: game.Winner,
		Duration:   duration,
	}

	if err := gm.producer.PublishEvent(event); err != nil {
		log.Printf("Error publishing game event: %v\n", err)
		// Don't fail - game is already saved
	}

	return nil
}

func (gm *GameManager) PublishMoveEvent(gameID string, player string, opponent string, column int, row int, isBot bool) error {
	event := GameEvent{
		EventType: "game_move",
		GameID:    gameID,
		Player:    player,
		Opponent:  opponent,
		Action:    "move",
		Column:    column,
		Row:       row,
		Timestamp: time.Now(),
		IsBot:     isBot,
	}

	return gm.producer.PublishEvent(event)
}
