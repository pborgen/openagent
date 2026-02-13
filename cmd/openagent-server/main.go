package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type RunRequest struct {
	Workflow string            `json:"workflow"`
	Params   map[string]string `json:"params"`
}

type RunResponse struct {
	RunID string `json:"runId"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req RunRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid json"))
			return
		}
		resp := RunResponse{RunID: "run-" + "0001"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	addr := ":7341"
	if v := os.Getenv("OPENAGENT_ADDR"); v != "" {
		addr = v
	}
	log.Printf("OpenAgent backend listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
