# ⚾ Go! Baseball 

## About the Project

Two words...baseball simulator. At its core, this project started from a personal obsession—I wanted to build baseball games that *come to life* at the click of a button. But on a more practical level, I built this to create opportunities—for myself and others—to access usable, pitch-by-pitch level data for baseball analytics.

Too many problems in baseball analytics go unsolved because we just don’t have enough samples to work with. **Go! Baseball** is here to change that. It generates as much data as you’ll ever need, all in a format that’ll feel instantly familiar to anyone who’s worked with Baseball Savant (aka Statcast).

## How It Works

- Simulates full 9-inning baseball games given multiple user inputs
- Models pitch-by-pitch outcomes using real world probability distributions
- Modular design using Go packages:
  - `sim`: simulation engine
  - `models`: shared data structures
  - `utils`: helper functions
  - `config`: database connections
  - `poster`: inserts results into database

## Technologies

- ![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
- ![Baseball Savant](https://img.shields.io/badge/Baseball_Savant-0e6ba8?style=for-the-badge&logo=mlb&logoColor=white)
- ![Python](https://img.shields.io/badge/Python-3670A0?style=for-the-badge&logo=python&logoColor=ffdd54)
- ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)

## Getting Started

### Prerequisites

- Go 1.21+
- A PostgreSQL database (connection configured in `config.ConnectDB()`)

### Run the simulator

```bash
go run main.go
