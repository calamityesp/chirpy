package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileServerHits int
}

func (apiconf *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiconf.fileServerHits++
		next.ServeHTTP(w, r)
	})
}

// handler to display numbe rof hits
func (apiconf *apiConfig) metricsHandler(rw http.ResponseWriter, r *http.Request) {
	response := fmt.Sprintf("Hits: %d", apiconf.fileServerHits)
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(response))
}

// reset handler
func (apiconf *apiConfig) resetHandler(rw http.ResponseWriter, r *http.Request) {
	apiconf.fileServerHits = 0
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Reset successful"))
}

type apiHandler struct{}

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func readinessHandler(rw http.ResponseWriter, rq *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("OK"))
}

func (apiconf *apiConfig) adminMetricsHandler(rw http.ResponseWriter, r *http.Request) {
	visitCount := apiconf.fileServerHits

	// Set the Content-Type header to "text/html"
	rw.Header().Set("Content-Type", "text/html")

	// Create the HTML response
	html := fmt.Sprintf(`
        <html>
            <body>
                <h1>Welcome, Chirpy Admin</h1>
                <p>Chirpy has been visited %d times!</p>
            </body>
        </html>`, visitCount)

	// Write the HTML response
	rw.Write([]byte(html))
}

func main() {
	mux := http.NewServeMux()
	apiconf := apiConfig{}

	// handle the root route with a fileserver
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/app/", apiconf.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))

	// server static files from subdirectoyr "assets"
	mux.Handle("assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	//register readinessHandler
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}
	// Register the reset handler
	mux.HandleFunc("/api/reset", apiconf.resetHandler)

	// Register the metrics handler
	mux.HandleFunc("GET /api/metrics", apiconf.metricsHandler)

	// Register the Adminmetrics handler
	mux.HandleFunc("GET /admin/metrics", apiconf.adminMetricsHandler)

	// Start web server and log errors
	fmt.Printf("Starting server on %s\n", server.Addr)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}
