# ⚾ Go! Baseball -  Baseball Simulator

This project is a pitch-by-pitch baseball game simulator written in Go. It uses real-world player data to simulate innings, at-bats, and outcomes with probabilistic models - the end output being a pitch-by-pitch dataframe that includes many fields provided in Baseball Savant data.

## 📦 Features

- Simulates full 9-inning baseball games
- Models pitch-by-pitch outcomes using real batter/pitcher matchups
- Tracks base state, runs, and outs dynamically
- Modular design using Go packages:
  - `sim`: simulation engine
  - `models`: shared data structures
  - `utils`: helper functions
  - `config`: database connections
  - `poster`: inserts results into database

## 🛠 Technologies

- **Baseball Savant** (for deriving probabilities)
- **Python** (for data cleaning/calculations)
- **Go (Golang)**
- **PostgreSQL** (for storing and retrieving player/game data)

## 🧪 How It Works

Each at-bat is simulated using `SimulateAtBat()` which returns a `PlateAppearanceResult`. The `ProcessPlateAppearance()` function updates the game state (score, base runners, outs) based on the result.

Key components:
- Base runners are tracked via a `[]bool` base state
- Events like singles, doubles, walks, and home runs move runners or score them
- Full game result is stored in `GameResult` structs

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- A PostgreSQL database (connection configured in `config.ConnectDB()`)

### Run the simulator

```bash
go run main.go
