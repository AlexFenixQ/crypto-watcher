package db

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq" // ← обязательно!
)

var (
	DB      *sql.DB
	mu      sync.Mutex
	tracked = map[string]bool{}
)

func Init() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	var err error
	for i := 0; i < 10; i++ {
		DB, err = sql.Open("postgres", dsn)
		if err == nil && DB.Ping() == nil {
			break
		}
		fmt.Println("⏳ Waiting for database...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		panic(err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS prices (
        coin TEXT,
        price FLOAT,
        timestamp BIGINT
    )`)
	if err != nil {
		panic(err)
	}
}

func TrackCurrency(coin string) {
	mu.Lock()
	tracked[coin] = true
	mu.Unlock()
}

func UntrackCurrency(coin string) {
	mu.Lock()
	delete(tracked, coin)
	mu.Unlock()
}

func GetTrackedCoins() []string {
	mu.Lock()
	defer mu.Unlock()
	coins := []string{}
	for c := range tracked {
		coins = append(coins, c)
	}
	return coins
}

func SavePrice(coin string, price float64, timestamp int64) {
	_, _ = DB.Exec("INSERT INTO prices (coin, price, timestamp) VALUES ($1, $2, $3)", coin, price, timestamp)
}

func FindClosestPrice(coin string, target int64) (float64, int64) {
	row := DB.QueryRow(`
        SELECT price, timestamp FROM prices
        WHERE coin=$1
        ORDER BY ABS(timestamp - $2)
        LIMIT 1
    `, coin, target)

	var price float64
	var ts int64
	_ = row.Scan(&price, &ts)
	return price, ts
}
