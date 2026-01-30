package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// --- 1. STRUKTURA DANYCH ---
// Używamy API Binance, bo jest najprostsze i zawsze działa.
type BinanceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// --- 2. METRYKI ---
var bitcoinRequests = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "bitcoin_price_requests_total",
	Help: "Całkowita liczba zapytań o cenę Bitcoina",
})

func main() {
	// Konfiguracja portu
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Opcjonalne logowanie klucza API (tylko dla testu zmiennych env)
	apiKey := os.Getenv("API_KEY")
	if len(apiKey) > 3 {
		log.Printf("Start z kluczem API: %s***", apiKey[0:3])
	}

	// --- 3. REJESTR PROMETHEUSA ---
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewGoCollector())
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(bitcoinRequests)

	// --- 4. ENDPOINTY ---
	http.HandleFunc("/bitcoin", getBitcoinPrice)
	// Podpinamy nasz rejestr pod /metrics
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.HandleFunc("/healthz/live", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    })
	http.HandleFunc("/healthz/ready", readinessHandler)

	log.Printf("Aplikacja startuje na porcie %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func getBitcoinPrice(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	// Podbijamy licznik odwiedzin
	bitcoinRequests.Inc()

	// Pobieramy dane z Binance (proste API, bez kluczy)
	url := "https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT"
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Błąd połączenia z API", http.StatusServiceUnavailable)
		log.Printf("Błąd API: %v", err)
		return
	}
	defer resp.Body.Close()

	// Dekodujemy JSON do struktury BinanceResponse
	var data BinanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, "Błąd odczytu danych", http.StatusInternalServerError)
		log.Printf("Błąd JSON: %v", err)
		return
	}

	duration := time.Since(start)
	log.Printf("Pobrano cenę: %s USD w czasie %v", data.Price, duration)

	w.Header().Set("Content-Type", "application/json")
	// Zwracamy czysty JSON do przeglądarki
	fmt.Fprintf(w, `{"currency": "USD", "price": "%s", "source": "Binance"}`, data.Price)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	client := http.Client{Timeout: 2 * time.Second}
	_, err := client.Get("https://www.google.com")
	
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("NOT READY"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("READY"))
}