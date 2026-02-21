package ipc

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nodeveil/nodeveil/engine/internal/service"
	"github.com/nodeveil/nodeveil/engine/internal/version"
)

type Server struct {
	httpServer *http.Server
	svc        *service.IndexService
}

func New(addr string, svc *service.IndexService) *Server {
	s := &Server{svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("/version", func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, version.Get()) })
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, map[string]string{"status": "ok"}) })
	mux.HandleFunc("/roots", s.roots)
	mux.HandleFunc("/index/status", func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, s.svc.GetIndexStatus()) })
	mux.HandleFunc("/search", s.search)
	s.httpServer = &http.Server{Addr: addr, Handler: mux}
	return s
}
func (s *Server) Start() error { return s.httpServer.ListenAndServe() }

func (s *Server) roots(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		roots, err := s.svc.ListRoots(r.Context())
		if err != nil {
			writeJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}
		paths := make([]string, 0, len(roots))
		for _, root := range roots {
			paths = append(paths, root.Path)
		}
		writeJSON(w, 200, map[string]any{"roots": paths})
	case http.MethodPost:
		var body struct {
			Path string `json:"path"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		if err := s.svc.AddRoot(r.Context(), body.Path); err != nil {
			writeJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, map[string]bool{"ok": true})
	case http.MethodDelete:
		path := r.URL.Query().Get("path")
		if err := s.svc.RemoveRoot(r.Context(), path); err != nil {
			writeJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, map[string]bool{"ok": true})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	ext := r.URL.Query().Get("ext")
	kind := r.URL.Query().Get("kind")
	rootID, _ := strconv.ParseInt(r.URL.Query().Get("rootId"), 10, 64)
	res, err := s.svc.SearchFiles(r.Context(), q, ext, rootID, kind)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"results": res})
}
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
