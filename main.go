package main

import (
	"log"
	"net/http"

	"github.com/calamityesp/chirpy/internal/database"
)

type apiConfig struct {
	fileserverhits int
	DB             *database.DB
}

func main() {
	// local variables
	const filepath = "./"
	const port = "8080"

	// create the database if it doesnt exists
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverhits: 0,
		DB:             db,
	}
	// setup routing multiplexer
	mux := http.NewServeMux()

	//create file system handler
	fshandler := apiCfg.middlewareMetricsinc(http.StripPrefix("/app", http.FileServer(http.Dir(filepath))))
	mux.Handle("/app/*", fshandler)

	//handlers
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{chirpId}", apiCfg.handlerChirpRetrieveById)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving file from %s on port: %s\v", filepath, port)
	log.Fatal(srv.ListenAndServe())
}
