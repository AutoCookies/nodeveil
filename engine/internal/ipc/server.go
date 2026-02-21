package ipc

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/nodeveil/nodeveil/engine/internal/config"
	"github.com/nodeveil/nodeveil/engine/internal/domain"
	"github.com/nodeveil/nodeveil/engine/internal/service"
	"github.com/nodeveil/nodeveil/engine/internal/version"
)

type Server struct {
	httpServer    *http.Server
	indexSvc      *service.IndexService
	graphSvc      *service.GraphService
	graphQuerySvc *service.GraphQueryService
}

func New(addr string, indexSvc *service.IndexService, graphSvc *service.GraphService, graphQuerySvc *service.GraphQueryService) *Server {
	s := &Server{indexSvc: indexSvc, graphSvc: graphSvc, graphQuerySvc: graphQuerySvc}
	mux := http.NewServeMux()
	mux.HandleFunc("/version", func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, version.Get()) })
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, map[string]string{"status": "ok"}) })
	mux.HandleFunc("/roots", s.roots)
	mux.HandleFunc("/index/status", func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, s.indexSvc.GetIndexStatus()) })
	mux.HandleFunc("/search", s.search)
	mux.HandleFunc("/graph/links", s.links)
	mux.HandleFunc("/graph/backlinks", s.backlinks)
	mux.HandleFunc("/graph/neighbors", s.neighbors)
	mux.HandleFunc("/graph/stats", s.stats)
	mux.HandleFunc("/graph/neighborhood", s.neighborhood)
	mux.HandleFunc("/graph/counts", s.counts)
	mux.HandleFunc("/node/meta", s.nodeMeta)
	mux.HandleFunc("/diagnostics/events", s.events)
	mux.HandleFunc("/open", s.openPath)
	mux.HandleFunc("/config/info", s.configInfo)
	mux.HandleFunc("/backup", s.backup)
	s.httpServer = &http.Server{Addr: addr, Handler: mux}
	return s
}

func (s *Server) Start() error { return s.httpServer.ListenAndServe() }

func (s *Server) roots(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		roots, err := s.indexSvc.ListRoots(r.Context())
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
		if err := s.indexSvc.AddRoot(r.Context(), body.Path); err != nil {
			writeJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, map[string]bool{"ok": true})
	case http.MethodDelete:
		if err := s.indexSvc.RemoveRoot(r.Context(), r.URL.Query().Get("path")); err != nil {
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
	res, err := s.indexSvc.SearchFiles(r.Context(), q, ext, rootID, kind)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"results": res})
}
func (s *Server) links(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		nodeID := r.URL.Query().Get("nodeId")
		d := domain.Direction(r.URL.Query().Get("direction"))
		if d == "" {
			d = domain.DirectionBoth
		}
		edges, err := s.graphSvc.ListLinks(r.Context(), nodeID, d)
		if err != nil {
			writeJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, map[string]any{"links": edges})
	case http.MethodPost:
		var body struct {
			FromID       string `json:"fromId"`
			ToID         string `json:"toId"`
			RelationType string `json:"relationType"`
			Note         string `json:"note"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		if err := s.graphSvc.CreateLink(r.Context(), body.FromID, body.ToID, body.RelationType, body.Note); err != nil {
			writeJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, map[string]bool{"ok": true})
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id != "" {
			_ = s.graphSvc.RemoveLink(r.Context(), id)
			writeJSON(w, 200, map[string]bool{"ok": true})
			return
		}
		if err := s.graphSvc.RemoveLinkByNodes(r.Context(), r.URL.Query().Get("fromId"), r.URL.Query().Get("toId"), r.URL.Query().Get("relationType")); err != nil {
			writeJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, map[string]bool{"ok": true})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func (s *Server) backlinks(w http.ResponseWriter, r *http.Request) {
	edges, err := s.graphSvc.GetBacklinks(r.Context(), r.URL.Query().Get("nodeId"))
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"links": edges})
}
func (s *Server) neighbors(w http.ResponseWriter, r *http.Request) {
	edges, err := s.graphSvc.GetNeighbors(r.Context(), r.URL.Query().Get("nodeId"), 1)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"links": edges})
}
func (s *Server) stats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.graphSvc.GetGraphStats(r.Context())
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, stats)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) nodeMeta(w http.ResponseWriter, r *http.Request) {
	n, err := s.indexSvc.GetFileByID(r.Context(), r.URL.Query().Get("id"))
	if err != nil || n == nil {
		writeJSON(w, 404, map[string]string{"error": "node not found"})
		return
	}
	writeJSON(w, 200, map[string]any{"node": n})
}
func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	ev, err := s.graphSvc.Events(r.Context(), limit)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"events": ev})
}
func (s *Server) openPath(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]bool{"ok": true})
}

func (s *Server) configInfo(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]string{"configPath": config.ConfigDir(), "dbPath": config.DefaultDBPath()})
}
func (s *Server) backup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	out := filepath.Join(config.ConfigDir(), "backup-graph.json")
	if err := os.MkdirAll(config.ConfigDir(), 0o755); err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	if err := s.graphSvc.ExportGraph(r.Context(), out); err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]string{"backup": out})
}
