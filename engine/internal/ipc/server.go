package ipc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nodeveil/nodeveil/engine/internal/db"
	"github.com/nodeveil/nodeveil/engine/internal/version"
)

type Server struct {
	httpServer *http.Server
	store      db.RootRepository
}

func New(addr string, store db.RootRepository) *Server {
	s := &Server{store: store}
	mux := http.NewServeMux()
	mux.HandleFunc("/version", s.handleVersion)
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/roots", s.handleRoots)

	s.httpServer = &http.Server{Addr: addr, Handler: mux}
	return s
}

func (s *Server) Start() error {
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("listen and serve: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}
	return nil
}

func (s *Server) handleVersion(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, version.Get())
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleRoots(w http.ResponseWriter, r *http.Request) {
	roots, err := s.store.ListRoots(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	paths := make([]string, 0, len(roots))
	for _, root := range roots {
		paths = append(paths, root.Path)
	}
	writeJSON(w, http.StatusOK, map[string][]string{"roots": paths})
}

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}
