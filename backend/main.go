package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = "game_events"
	}

	// Initialize database
	db, err := InitDB(dbURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize Kafka producer (optional, skip if not available)
	var producer *KafkaProducer
	if kafkaBroker != "" {
		var err error
		producer, err = NewKafkaProducer(kafkaBroker, kafkaTopic)
		if err != nil {
			log.Println("Warning: Kafka not available, analytics disabled:", err)
		}
	}
	if producer != nil {
		defer producer.Close()
	}

	// Initialize game manager
	gameManager := NewGameManager(db, producer)

	// Initialize WebSocket hub
	hub := NewHub(gameManager)
	go hub.Run()

	// Start server
	server := NewServer(port, hub, gameManager, db)
	log.Printf("Server starting on port %s\n", port)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
