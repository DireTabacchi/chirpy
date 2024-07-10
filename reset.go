package main

import (
    "fmt"
    "net/http"
)

func (cfg *apiConfig) resetMetricsHandler(w http.ResponseWriter, r *http.Request) {
    cfg.fileserverHits = 0
    respBody := fmt.Sprintf("Metrics reset.\nHits: %d", cfg.fileserverHits)
    w.Header().Add("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(respBody))
}
