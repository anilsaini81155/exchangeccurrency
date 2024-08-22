package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/anilsaini81155/exchangeccurrency/internal/internalredis"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type PutRatesRequest struct {
	Username     string  `json:"username"`
	Exchangepair string  `json:"exchange_pair"`
	Exchangerate float64 `json:"exchange_rate"`
}

type GetRatesResponse struct {
	Exchangepair string  `json:"exchange_pair"`
	Exchangerate float64 `json:"exchange_rate"`
	Message      string  `json:"message"`
}

// Allowed exchange value pairs
var allowedPairs = map[string]bool{
	"INR/USD": true,
	"INR/EUR": true,
	"USD/EUR": true,
}

// ValidateExchangeRate validates that the required fields are present and correct
func ValidateExchangeRate(exchangefields PutRatesRequest) error {
	if exchangefields.Exchangerate == 0 {
		return errors.New("exchange_rate is required and cannot be zero")
	}

	if exchangefields.Username == "" {
		return errors.New("username is required")
	}

	if !allowedPairs[exchangefields.Exchangepair] {
		return errors.New("invalid exchange_pair. Allowed values are INR/USD, INR/EUR, USD/EUR")
	}

	return nil
}

// ValidateGetExchangePair checks if the provided exchange pair is valid
func ValidateGetExchangePair(pair string) bool {
	return allowedPairs[pair]
}

func StoreExchangeRate(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var putrates PutRatesRequest
	redisClient := internalredis.SetupRedis()
	err := json.NewDecoder(r.Body).Decode(&putrates)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ValidateExchangeRate(putrates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	StoreRate(redisClient, putrates.Exchangepair, putrates.Exchangerate)

	redisClient.Publish(ctx, "exchange_rates", fmt.Sprintf("%s: %.2f", putrates.Exchangepair, putrates.Exchangerate))

	response := map[string]interface{}{
		"exchange_rate":  putrates.Exchangerate,
		"exchange_value": putrates.Exchangepair,
		"username":       putrates.Username,
		"message":        "Rates stores successfully.",
	}

	errencode := json.NewEncoder(w).Encode(response)
	if errencode != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func StoreRate(client *redis.Client, currencyPair string, rate float64) {

	err := client.HSet(ctx, "exchange_rates", currencyPair, rate).Err()

	if err != nil {
		log.Fatalf("Could not set exchange rate: %v", err)
	}

}

func GetExchangeRate(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	exchange_pair := r.URL.Query().Get("exchange_pair")

	if exchange_pair == "" {
		http.Error(w, "Missing required query parameter: exchange_pair", http.StatusBadRequest)
		return
	}

	// Validate the parameter
	if !ValidateGetExchangePair(exchange_pair) {
		http.Error(w, "Invalid exchange pair. Allowed values are INR/USD, INR/EUR, USD/EUR.", http.StatusBadRequest)
		return
	}

	redisClient := internalredis.SetupRedis()

	rate := GetExchangeRates(redisClient, exchange_pair)

	var getratesresponse GetRatesResponse
	getratesresponse.Exchangepair = exchange_pair
	getratesresponse.Exchangerate = rate
	getratesresponse.Message = "Exchange rates and pair fetched successfully."

	w.WriteHeader(http.StatusOK)

	errencode := json.NewEncoder(w).Encode(getratesresponse)
	if errencode != nil {
		http.Error(w, errencode.Error(), http.StatusInternalServerError)
		return
	}
}

func GetExchangeRates(client *redis.Client, currencyPair string) float64 {
	rate, err := client.HGet(ctx, "exchange_rates", currencyPair).Float64()
	if err != nil {
		log.Fatalf("Could not get exchange rate: %v", err)
	}
	return rate
}

func ExchangeRates(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		StoreExchangeRate(w, r)
	case http.MethodGet:
		GetExchangeRate(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

}
