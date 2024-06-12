package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/calamityesp/chirpy/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverhits int
	DB             *database.DB
	secret_Key     string
}

func main() {
	// load env variables
	godotenv.Load()
	// local variables
	const filepath = "./"
	const port = "8080"
	const databasePath = "database.json"

	// check for debug flag
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
		secret_Key:     os.Getenv("JWT_SECRET"),
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
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)
	mux.HandleFunc("POST /api/login", apiCfg.handlerUsersLogin)

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
