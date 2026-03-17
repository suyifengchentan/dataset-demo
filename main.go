package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"

	"database-demo/internal/demo"
	"database-demo/internal/live"
)

//go:embed static
var staticFiles embed.FS

func main() {
	staticRoot, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticRoot))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFileFS(w, r, staticRoot, "index.html")
	})
	mux.HandleFunc("/api/demo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(demo.Build()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	analyzer := live.NewAnalyzer()
	simulator := live.NewSimulator()
	mux.HandleFunc("/api/live/analyze", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req live.AnalyzeRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "invalid request: "+err.Error(), http.StatusBadRequest)
			return
		}

		report := analyzer.Analyze(r.Context(), req)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(report); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("/api/live/simulate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req live.AnalyzeRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "invalid request: "+err.Error(), http.StatusBadRequest)
			return
		}

		report := simulator.Simulate(r.Context(), req)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(report); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	addr := ":8080"
	log.Printf("database demo listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
