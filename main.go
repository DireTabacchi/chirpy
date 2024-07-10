package main

import (
    "log"
    "net/http"
)

type apiConfig struct {
    fileserverHits int
}

func main() {
    const port = "8080"
    const filepathRoot = "."

    apiCfg := &apiConfig{fileserverHits: 0}

    serveMux := http.NewServeMux()
    fsHandler := apiCfg.middlewareMetricsInc(
        http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))),
    )

    serveMux.Handle("/app/*", fsHandler)

    serveMux.HandleFunc("GET /api/healthz", readinessHandler)
    serveMux.HandleFunc("GET /api/reset", apiCfg.resetMetricsHandler)
    serveMux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)

    serveMux.HandleFunc("POST /api/chirps", validateChirpHandler)

    srv := &http.Server{
        Addr: ":" + port,
        Handler: serveMux,
    }

    log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
    log.Fatal(srv.ListenAndServe())

}
