package store

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/nodeveil/nodeveil/engine/internal/domain"
)

const (
	maxGraphNodes = 10000
	maxGraphEdges = 50000
)

func (s *Store) QueryCounts(ctx context.Context, nodeID string) (domain.GraphCounts, error) {
	return s.GraphCounts(ctx, nodeID)
}

func (s *Store) QueryNeighborhood(ctx context.Context, nodeID string, depth int, filters domain.GraphFilters, cursor string) (domain.GraphNeighborhood, error) {
	if depth < 1 {
		depth = 1
	}
	if depth > 3 {
		depth = 3
	}
	offset := 0
	if cursor != "" {
		offset, _ = strconv.Atoi(cursor)
		if offset < 0 {
			offset = 0
		}
	}
	visited := map[string]struct{}{nodeID: {}}
	frontier := []string{nodeID}
	edges := make([]domain.Edge, 0, maxGraphEdges)

	for d := 0; d < depth; d++ {
		next := []string{}
		for _, id := range frontier {
			where := "deleted_at IS NULL"
			switch filters.Direction {
			case domain.DirectionIn:
				where += fmt.Sprintf(" AND to_id='%s'", esc(id))
			case domain.DirectionOut:
				where += fmt.Sprintf(" AND from_id='%s'", esc(id))
			default:
				where += fmt.Sprintf(" AND (from_id='%s' OR to_id='%s')", esc(id), esc(id))
			}
			if len(filters.RelationTypes) > 0 {
				rels := make([]string, 0, len(filters.RelationTypes))
				for _, r := range filters.RelationTypes {
					rels = append(rels, fmt.Sprintf("'%s'", esc(r)))
				}
				where += " AND relation_type IN (" + strings.Join(rels, ",") + ")"
			}
			rows, err := s.exec(ctx, "SELECT id||'|'||from_id||'|'||to_id||'|'||relation_type||'|'||ifnull(note,'')||'|'||created_at||'|'||updated_at FROM edges WHERE "+where+" ORDER BY from_id,to_id,relation_type LIMIT "+strconv.Itoa(maxGraphEdges)+" OFFSET "+strconv.Itoa(offset)+";")
			if err != nil {
				return domain.GraphNeighborhood{}, err
			}
			for _, ln := range strings.Split(strings.TrimSpace(rows), "\n") {
				if ln == "" {
					continue
				}
				p := strings.Split(ln, "|")
				if len(p) < 7 {
					continue
				}
				ca, _ := strconv.ParseInt(p[5], 10, 64)
				ua, _ := strconv.ParseInt(p[6], 10, 64)
				e := domain.Edge{ID: p[0], FromID: p[1], ToID: p[2], RelationType: domain.RelationType(p[3]), Note: p[4], CreatedAt: ca, UpdatedAt: ua}
				edges = append(edges, e)
				if _, ok := visited[e.FromID]; !ok {
					visited[e.FromID] = struct{}{}
					next = append(next, e.FromID)
				}
				if _, ok := visited[e.ToID]; !ok {
					visited[e.ToID] = struct{}{}
					next = append(next, e.ToID)
				}
				if len(edges) >= maxGraphEdges {
					break
				}
			}
			if len(edges) >= maxGraphEdges {
				break
			}
		}
		frontier = next
		if len(edges) >= maxGraphEdges {
			break
		}
	}

	ids := make([]string, 0, len(visited))
	for id := range visited {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	nodes := make([]domain.Node, 0, len(ids))
	for _, id := range ids {
		n, err := s.GetFileByID(ctx, id)
		if err != nil || n == nil {
			continue
		}
		if len(filters.Extensions) > 0 {
			ok := false
			for _, ex := range filters.Extensions {
				if strings.EqualFold(n.Ext, ex) {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}
		nodes = append(nodes, *n)
		if len(nodes) >= maxGraphNodes {
			break
		}
	}
	sort.Slice(edges, func(i, j int) bool {
		a := edges[i]
		b := edges[j]
		if a.FromID != b.FromID {
			return a.FromID < b.FromID
		}
		if a.ToID != b.ToID {
			return a.ToID < b.ToID
		}
		return a.RelationType < b.RelationType
	})
	truncated := len(edges) >= maxGraphEdges || len(nodes) >= maxGraphNodes
	var next *string
	if truncated {
		nxt := strconv.Itoa(offset + len(edges))
		next = &nxt
	}
	if len(edges) > maxGraphEdges {
		edges = edges[:maxGraphEdges]
	}
	if len(nodes) > maxGraphNodes {
		nodes = nodes[:maxGraphNodes]
	}
	return domain.GraphNeighborhood{Nodes: nodes, Edges: edges, Truncated: truncated, NextCursor: next}, nil
}
