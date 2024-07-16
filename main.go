package main

import (
	"log"
	"net/http"
	"os"

	"github.com/DireTabacchi/chirpy/internal/database"
    "github.com/joho/godotenv"
)

type apiConfig struct {
    fileserverHits int
    db *database.DB
    jwtSecret string
}

func main() {
    const port = "8080"
    const filepathRoot = "."

    db, err := database.NewDB("database.json")
    if err != nil {
        log.Fatal(err)
    }

    err = godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading secrets: %v\n", err)
    }

    apiCfg := &apiConfig{
        fileserverHits: 0,
        db: db,
        jwtSecret: os.Getenv("JWT_SECRET"),
    }

    serveMux := http.NewServeMux()
    fsHandler := apiCfg.middlewareMetricsInc(
        http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))),
    )

    serveMux.Handle("/app/*", fsHandler)

    // Maintenance/Metrics endpoints
    serveMux.HandleFunc("GET /api/healthz", readinessHandler)
    serveMux.HandleFunc("GET /api/reset", apiCfg.resetMetricsHandler)
    serveMux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)

    // Chirp endpoints
    serveMux.HandleFunc("POST /api/chirps", apiCfg.postChirpHandler)
    serveMux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)
    serveMux.HandleFunc("GET /api/chirps/{id}", apiCfg.getChirpHandler)

    // User endpoints
    serveMux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
    serveMux.HandleFunc("PUT /api/users", apiCfg.updateUserHandler)

    serveMux.HandleFunc("POST /api/login", apiCfg.loginHandler)
    serveMux.HandleFunc("POST /api/refresh", apiCfg.refreshHandler)
    serveMux.HandleFunc("POST /api/revoke", apiCfg.revokeHandler)

    srv := &http.Server{
        Addr: ":" + port,
        Handler: serveMux,
    }

    log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
    log.Fatal(srv.ListenAndServe())

}
