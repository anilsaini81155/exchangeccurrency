package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var rdb *redis.Client

type PutRatesRequest struct {
	Username     string  `json:"username"`
	Exchangepair string  `json:"exchangepair"`
	Exchangerate float64 `json:"exchangerate"`
}

type GetRatesRequest struct {
	Username     string `json:"username"`
	Exchangepair string `json:"exchangepair"`
}

type GetRatesResponse struct {
	Exchangepair string  `json:"exchangepair"`
	Exchangerate float64 `json:"exchangerate"`
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
}

func StoreExchangeRate(w http.ResponseWriter, r *http.Request) {
	var putrates PutRatesRequest
	initRedis()
	json.NewDecoder(r.Body).Decode(&putrates)
	StoreRate(rdb, putrates.Exchangepair, putrates.Exchangerate)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Rates stores successfully.")
}

func StoreRate(client *redis.Client, currencyPair string, rate float64) {
	err := client.HSet(ctx, "exchange_rates", currencyPair, rate).Err()
	if err != nil {
		log.Fatalf("Could not set exchange rate: %v", err)
	}

}

func GetExchangeRate(w http.ResponseWriter, r *http.Request) {
	var getrates GetRatesRequest
	initRedis()
	json.NewDecoder(r.Body).Decode(&getrates)
	rate := GetExchangeRates(rdb, getrates.Exchangepair)

	var getratesresponse GetRatesResponse
	getratesresponse.Exchangepair = getrates.Exchangepair
	getratesresponse.Exchangerate = rate

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(getratesresponse)
}

func GetExchangeRates(client *redis.Client, currencyPair string) float64 {
	rate, err := client.HGet(ctx, "exchange_rates", currencyPair).Float64()
	if err != nil {
		log.Fatalf("Could not get exchange rate: %v", err)
	}
	return rate
}

// func main() {
// client := redis.NewClient(&redis.Options{
// 	Addr: "localhost:6379",
// })

// Example usage
// StoreRate(client, "USD:EUR", 0.85)
// rate := GetExchangeRate(client, "USD:EUR")
// fmt.Printf("Exchange rate USD:EUR = %f\n", rate)
// }
