package main

import (
    "log"
    "net/http"
)

func main() {
    const port = "8080"
    const filepathRoot = "."
    const filepathLogo = "./assets/logo.png"

    serveMux := http.NewServeMux()
    serveMux.Handle("/app/*",
        http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))),
    )
    serveMux.Handle("/app/assets/",
        http.StripPrefix("/app/assets/", http.FileServer(http.Dir("assets"))),
    )
    serveMux.HandleFunc("/healthz", readinessHandler)

    srv := &http.Server{
        Addr: ":" + port,
        Handler: serveMux,
    }

    log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
    log.Fatal(srv.ListenAndServe())

}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(200)
    _, err := w.Write([]byte("OK"))
    if err != nil {
        log.Printf("Readiness handler exprienced error while writing body:\n%v", err)
    }
}
