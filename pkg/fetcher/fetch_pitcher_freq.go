package fetcher

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	//"github.com/logananthony/go-baseball/pkg/config"
)

// FetchPitcherFrequencies queries and prints pitch data for a pitcher ID
func FetchPitcherFrequencies(db *sql.DB, pitcherID int, stand string) {
	query := `SELECT stand, pitch_type, balls, strikes, count, frequency FROM pitcher_count_pitch_freq WHERE pitcher = $1 AND stand = $2`

	rows, err := db.Query(query, pitcherID, stand)
	if err != nil {
		log.Fatal("Query error:", err)
	}
	defer rows.Close()

	fmt.Println("Results for pitcher:", pitcherID)
	for rows.Next() {
		var standType, pitchType string
		var balls, strikes, count int
		var frequency float64

		if err := rows.Scan(&standType, &pitchType, &balls, &strikes, &count, &frequency); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Pitch: %s | Balls: %d | Strikes: %d | Count: %d | Freq: %.2f\n", pitchType, balls, strikes, count, frequency)
	}
}

