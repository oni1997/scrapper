package main

import (
	"encoding/json"
	"log"
	"net/http"
	"scrapper/models"
	"scrapper/scraper"
)


func main() {
    s := scraper.NewScraper()

    http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        if r.Method != http.MethodPost {
            w.WriteHeader(http.StatusMethodNotAllowed)
            json.NewEncoder(w).Encode(map[string]string{
                "error": "Only POST method is allowed",
            })
            return
        }

        var req models.SearchRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{
                "error": "Invalid request body: " + err.Error(),
            })
            return
        }

        if req.SearchTerm == "" {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{
                "error": "SearchTerm cannot be empty",
            })
            return
        }

        result, err := s.ScrapeProducts(req.SearchTerm)
        if err != nil {
            log.Printf("Scraping error: %v", err)
            // Still return the result as it contains error information
            // but with a 500 status code
            w.WriteHeader(http.StatusInternalServerError)
        }

        if err := json.NewEncoder(w).Encode(result); err != nil {
            log.Printf("Error encoding response: %v", err)
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(map[string]string{
                "error": "Error encoding response",
            })
            return
        }
    })

    port := ":8080"
    log.Printf("Server starting on port %s", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}