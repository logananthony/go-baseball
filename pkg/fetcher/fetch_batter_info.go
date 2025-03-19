
package fetcher

import (
	"database/sql"
	//"fmt"
	"log"

	_ "github.com/lib/pq"
	//"github.com/logananthony/go-baseball/pkg/config"
)

// FetchPitcherFrequencies queries and prints pitch data for a pitcher ID
func FetchBatterInfo(db *sql.DB, batterID int, game_year int) (Bats string) {
	query := `SELECT bats FROM batter_info WHERE batter = $1 AND game_year = $2`

	rows, err := db.Query(query, batterID, game_year)
	if err != nil {
		log.Fatal("Query error:", err)
	}
	defer rows.Close()

	//fmt.Println("Results for batter:", batterID)
	for rows.Next() {
		//var Bats string

		if err := rows.Scan(&Bats); err != nil {
			log.Fatal(err)
		}

		//fmt.Printf("Bats: %s", Bats)
	}
  
  return Bats
}

