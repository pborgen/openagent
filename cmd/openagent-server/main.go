package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type RunRequest struct {
	Workflow string            `json:"workflow"`
	Params   map[string]string `json:"params"`
}

type RunResponse struct {
	RunID string `json:"runId"`
}

type LogResponse struct {
	Lines []string `json:"lines"`
	Done  bool     `json:"done"`
}

type RunState struct {
	Lines []string
	Done  bool
}

type Store struct {
	mu   sync.Mutex
	runs map[string]*RunState
}

func NewStore() *Store {
	return &Store{runs: map[string]*RunState{}}
}

func (s *Store) createRun() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := "run-" + strconv.Itoa(rand.Intn(100000))
	s.runs[id] = &RunState{Lines: []string{}, Done: false}
	return id
}

func (s *Store) append(id, line string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if r, ok := s.runs[id]; ok {
		r.Lines = append(r.Lines, line)
	}
}

func (s *Store) finish(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if r, ok := s.runs[id]; ok {
		r.Done = true
	}
}

func (s *Store) get(id string) (*RunState, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.runs[id]
	return r, ok
}

func main() {
	rand.Seed(time.Now().UnixNano())
	store := NewStore()
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
		runID := store.createRun()
		store.append(runID, "starting workflow: "+req.Workflow)
		go func() {
			for i := 1; i <= 5; i++ {
				store.append(runID, "step "+strconv.Itoa(i)+" complete")
				time.Sleep(800 * time.Millisecond)
			}
			store.append(runID, "done")
			store.finish(runID)
		}()

		resp := RunResponse{RunID: runID}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("runId")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing runId"))
			return
		}
		run, ok := store.get(id)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("unknown runId"))
			return
		}
		resp := LogResponse{Lines: run.Lines, Done: run.Done}
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
