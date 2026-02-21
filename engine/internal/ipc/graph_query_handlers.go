package ipc

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/nodeveil/nodeveil/engine/internal/domain"
)

func (s *Server) neighborhood(w http.ResponseWriter, r *http.Request) {
	depth, _ := strconv.Atoi(r.URL.Query().Get("depth"))
	nodeID := r.URL.Query().Get("nodeId")
	d := domain.Direction(r.URL.Query().Get("direction"))
	if d == "" {
		d = domain.DirectionBoth
	}
	filters := domain.GraphFilters{Direction: d}
	if rel := r.URL.Query().Get("relations"); rel != "" {
		filters.RelationTypes = strings.Split(rel, ",")
	}
	if ex := r.URL.Query().Get("ext"); ex != "" {
		filters.Extensions = strings.Split(ex, ",")
	}
	resp, err := s.graphQuerySvc.GetNeighborhood(r.Context(), nodeID, depth, filters, r.URL.Query().Get("cursor"))
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, resp)
}

func (s *Server) counts(w http.ResponseWriter, r *http.Request) {
	resp, err := s.graphQuerySvc.GetCounts(r.Context(), r.URL.Query().Get("nodeId"))
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, resp)
}
