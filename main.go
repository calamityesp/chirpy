package main

import (
	"flag"
	"github.com/calamityesp/chirpy/internal/database"
	"log"
	"net/http"
	"os"
)

type apiConfig struct {
	fileserverhits int
	DB             *database.DB
}

func main() {
	// local variables
	const filepath = "./"
	const port = "8080"
	const databasePath = "database.json"

	dbg := flag.Bool("debug", false, "Enable Debug Mode")
	flag.Parse()

	if *dbg {
		log.Print("Debug mode enabled, deleting json database")
		err := deleteDatabase(databasePath)
		if err != nil {
			log.Fatalf("Error deleting database: %s", err)
		}
	}

	// create the database if it doesnt exists
	db, err := database.NewDB(databasePath)
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
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving file from %s on port: %s\v", filepath, port)
	log.Fatal(srv.ListenAndServe())
}

func deleteDatabase(file string) error {
	err := os.Remove(file)
	if err != nil {
		return err
	}

	return nil
}
