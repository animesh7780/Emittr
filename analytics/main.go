package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

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

type AnalyticsDB struct {
	dbURL string
}

func main() {
	godotenv.Load()

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = "game_events"
	}

	kafkaGroup := os.Getenv("KAFKA_GROUP")
	if kafkaGroup == "" {
		kafkaGroup = "analytics_group"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	analyticsDB := &AnalyticsDB{dbURL: dbURL}

	// Initialize database
	if err := analyticsDB.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	log.Printf("Analytics service starting - Kafka: %s, Topic: %s, Group: %s\n", kafkaBroker, kafkaTopic, kafkaGroup)

	// Start consuming events
	if err := analyticsDB.ConsumeEvents(kafkaBroker, kafkaTopic, kafkaGroup); err != nil {
		log.Fatal("Failed to consume events:", err)
	}
}

func (a *AnalyticsDB) InitDB() error {
	// We'll store analytics in the same database
	// Just validate connection
	return nil
}

func (a *AnalyticsDB) ConsumeEvents(broker string, topic string, group string) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{broker},
		Topic:          topic,
		GroupID:        group,
		StartOffset:    kafka.LastOffset,
		CommitInterval: 1 * time.Second,
	})

	defer reader.Close()

	log.Printf("Connected to Kafka topic: %s\n", topic)

	ctx := context.Background()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			continue
		}

		var event GameEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Error unmarshaling event: %v\n", err)
			continue
		}

		log.Printf("Event received: %s - Game: %s, Player: %s\n", event.EventType, event.GameID, event.Player)

		// Process event
		a.ProcessEvent(event)
	}
}

func (a *AnalyticsDB) ProcessEvent(event GameEvent) {
	switch event.EventType {
	case "game_completed":
		a.ProcessGameCompletion(event)
	case "game_move":
		a.ProcessGameMove(event)
	default:
		log.Printf("Unknown event type: %s\n", event.EventType)
	}
}

func (a *AnalyticsDB) ProcessGameCompletion(event GameEvent) {
	log.Printf("Game completed - ID: %s, Winner: %s, Duration: %ds, IsBot: %v\n",
		event.GameID, event.GameResult, event.Duration, event.IsBot)

	// Log analytics
	if event.GameResult != "draw" && event.GameResult != "" {
		log.Printf("  Winner: %s\n", event.GameResult)
	}
	if event.IsBot {
		log.Printf("  Game was vs Bot\n")
	}

	// In production, you would:
	// 1. Store metrics in database
	// 2. Calculate aggregates (average game duration, most frequent winners, etc.)
	// 3. Track games per hour/day
	// 4. Track user-specific metrics
	a.LogGameMetrics(event)
}

func (a *AnalyticsDB) ProcessGameMove(event GameEvent) {
	log.Printf("Game move - ID: %s, Player: %s, Column: %d, Row: %d\n",
		event.GameID, event.Player, event.Column, event.Row)
}

func (a *AnalyticsDB) LogGameMetrics(event GameEvent) {
	now := time.Now()
	
	metrics := map[string]interface{}{
		"timestamp":     now.Format(time.RFC3339),
		"game_id":       event.GameID,
		"winner":        event.GameResult,
		"player":        event.Player,
		"opponent":      event.Opponent,
		"duration_sec":  event.Duration,
		"is_bot":        event.IsBot,
		"hour_of_day":   now.Hour(),
		"day_of_week":   now.Weekday(),
		"date":          now.Format("2006-01-02"),
	}

	fmt.Printf("ANALYTICS: %+v\n", metrics)

	// In production, save to analytics database
	// For now, just logging
}
