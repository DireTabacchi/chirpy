package main

import (
    "log"
    "net/http"

    "github.com/DireTabacchi/chirpy/internal/database"
)

type apiConfig struct {
    fileserverHits int
    db *database.DB
}

func main() {
    const port = "8080"
    const filepathRoot = "."

    db, err := database.NewDB("database.json")
    if err != nil {
        log.Fatal(err)
    }

    apiCfg := &apiConfig{
        fileserverHits: 0,
        db: db,

    }

    serveMux := http.NewServeMux()
    fsHandler := apiCfg.middlewareMetricsInc(
        http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))),
    )

    serveMux.Handle("/app/*", fsHandler)

    serveMux.HandleFunc("GET /api/healthz", readinessHandler)
    serveMux.HandleFunc("GET /api/reset", apiCfg.resetMetricsHandler)
    serveMux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)

    serveMux.HandleFunc("POST /api/chirps", apiCfg.postChirpHandler)
    serveMux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)

    srv := &http.Server{
        Addr: ":" + port,
        Handler: serveMux,
    }

    log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
    log.Fatal(srv.ListenAndServe())

}
