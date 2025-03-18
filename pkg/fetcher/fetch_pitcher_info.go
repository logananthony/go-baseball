package fetcher

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	//"github.com/logananthony/go-baseball/pkg/config"
)

// FetchPitcherFrequencies queries and prints pitch data for a pitcher ID
func FetchPitcherInfo(db *sql.DB, pitcherID int, game_year int) (Throws string) {
	query := `SELECT throws FROM pitcher_info WHERE pitcher = $1 AND game_year = $2`

	rows, err := db.Query(query, pitcherID, game_year)
	if err != nil {
		log.Fatal("Query error:", err)
	}
	defer rows.Close()

	//fmt.Println("Results for batter:", batterID)
	for rows.Next() {
		//var Bats string

		if err := rows.Scan(&Throws); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Throws: %s", Throws)
	}
  
  return Throws
}


