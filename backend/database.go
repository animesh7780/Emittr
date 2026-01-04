package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	conn *sql.DB
}

func InitDB(dbURL string) (*Database, error) {
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		return nil, err
	}

	db := &Database{conn: conn}

	// Run migrations
	if err := db.runMigrations(); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}

func (db *Database) runMigrations() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS players (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			wins INT DEFAULT 0,
			losses INT DEFAULT 0,
			draws INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS games (
			id VARCHAR(36) PRIMARY KEY,
			player1 VARCHAR(255),
			player2 VARCHAR(255),
			winner VARCHAR(255),
			is_bot BOOLEAN DEFAULT false,
			status VARCHAR(50),
			board_state JSONB,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,
			duration_seconds INT
		)`,
		`CREATE INDEX IF NOT EXISTS idx_games_player1 ON games(player1)`,
		`CREATE INDEX IF NOT EXISTS idx_games_player2 ON games(player2)`,
		`CREATE INDEX IF NOT EXISTS idx_games_winner ON games(winner)`,
		`CREATE INDEX IF NOT EXISTS idx_games_created_at ON games(created_at)`,
	}

	for _, migration := range migrations {
		_, err := db.conn.Exec(migration)
		if err != nil {
			log.Printf("Migration error: %v\n", err)
			// Continue even if some tables exist
		}
	}

	return nil
}

func (db *Database) SaveGame(game *GameState) error {
	createdAt, _ := time.Parse(time.RFC3339, game.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, game.UpdatedAt)
	duration := int(updatedAt.Sub(createdAt).Seconds())

	boardJSON := fmt.Sprintf(`"%v"`, game.Board.Grid)

	query := `
		INSERT INTO games (id, player1, player2, winner, is_bot, status, board_state, created_at, updated_at, duration_seconds)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE SET
			winner = $4,
			status = $6,
			updated_at = $9,
			duration_seconds = $10
	`

	_, err := db.conn.Exec(
		query,
		game.ID,
		game.Player1,
		game.Player2,
		game.Winner,
		game.IsBot,
		game.Status,
		boardJSON,
		createdAt,
		updatedAt,
		duration,
	)

	if err != nil {
		log.Printf("Error saving game: %v\n", err)
		return err
	}

	// Update player stats
	if game.Winner != "" && game.Winner != "draw" {
		if game.Winner == "Bot" {
			err := db.IncrementLosses(game.Player1)
			if err != nil {
				log.Printf("Error updating losses for %s: %v\n", game.Player1, err)
			}
		} else {
			err := db.IncrementWins(game.Winner)
			if err != nil {
				log.Printf("Error updating wins for %s: %v\n", game.Winner, err)
			}
			if game.IsBot {
				// No need to update bot stats
			} else {
				if game.Player1 == game.Winner {
					err := db.IncrementLosses(game.Player2)
					if err != nil {
						log.Printf("Error updating losses for %s: %v\n", game.Player2, err)
					}
				} else {
					err := db.IncrementLosses(game.Player1)
					if err != nil {
						log.Printf("Error updating losses for %s: %v\n", game.Player1, err)
					}
				}
			}
		}
	} else if game.Winner == "draw" {
		err := db.IncrementDraws(game.Player1)
		if err != nil {
			log.Printf("Error updating draws for %s: %v\n", game.Player1, err)
		}
		if !game.IsBot {
			err := db.IncrementDraws(game.Player2)
			if err != nil {
				log.Printf("Error updating draws for %s: %v\n", game.Player2, err)
			}
		}
	}

	return nil
}

func (db *Database) IncrementWins(username string) error {
	query := `
		INSERT INTO players (username, wins) VALUES ($1, 1)
		ON CONFLICT (username) DO UPDATE SET
			wins = wins + 1,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := db.conn.Exec(query, username)
	return err
}

func (db *Database) IncrementLosses(username string) error {
	query := `
		INSERT INTO players (username, losses) VALUES ($1, 1)
		ON CONFLICT (username) DO UPDATE SET
			losses = losses + 1,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := db.conn.Exec(query, username)
	return err
}

func (db *Database) IncrementDraws(username string) error {
	query := `
		INSERT INTO players (username, draws) VALUES ($1, 1)
		ON CONFLICT (username) DO UPDATE SET
			draws = draws + 1,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := db.conn.Exec(query, username)
	return err
}

func (db *Database) GetLeaderboard(limit int) ([]map[string]interface{}, error) {
	query := `
		SELECT username, wins, losses, draws, 
			   CAST(wins AS FLOAT) / NULLIF(wins + losses + draws, 0) as win_rate
		FROM players
		WHERE wins + losses + draws > 0
		ORDER BY wins DESC
		LIMIT $1
	`

	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	leaderboard := make([]map[string]interface{}, 0) // Initialize as empty slice instead of nil
	for rows.Next() {
		var username string
		var wins, losses, draws int
		var winRate float64

		if err := rows.Scan(&username, &wins, &losses, &draws, &winRate); err != nil {
			return nil, err
		}

		leaderboard = append(leaderboard, map[string]interface{}{
			"username": username,
			"wins":     wins,
			"losses":   losses,
			"draws":    draws,
			"winRate":  fmt.Sprintf("%.2f%%", winRate*100),
		})
	}

	return leaderboard, nil
}

func (db *Database) GetPlayerStats(username string) (map[string]interface{}, error) {
	query := `
		SELECT username, wins, losses, draws, created_at
		FROM players
		WHERE username = $1
	`

	var username_db string
	var wins, losses, draws int
	var createdAt time.Time

	err := db.conn.QueryRow(query, username).Scan(&username_db, &wins, &losses, &draws, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return map[string]interface{}{
				"username": username,
				"wins":     0,
				"losses":   0,
				"draws":    0,
				"winRate":  "0.00%",
			}, nil
		}
		return nil, err
	}

	total := wins + losses + draws
	winRate := 0.0
	if total > 0 {
		winRate = float64(wins) / float64(total) * 100
	}

	return map[string]interface{}{
		"username": username_db,
		"wins":     wins,
		"losses":   losses,
		"draws":    draws,
		"winRate":  fmt.Sprintf("%.2f%%", winRate),
		"createdAt": createdAt,
	}, nil
}

func (db *Database) Close() error {
	return db.conn.Close()
}
