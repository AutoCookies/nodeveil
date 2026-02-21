package store

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/nodeveil/nodeveil/engine/internal/domain"
)

type Store struct{ dbPath string }

func Open(ctx context.Context, dbPath string) (*Store, error) {
	s := &Store{dbPath: dbPath}
	if err := s.Migrate(ctx); err != nil {
		return nil, err
	}
	return s, nil
}
func (s *Store) Close() error { return nil }

func (s *Store) exec(ctx context.Context, sql string) (string, error) {
	cmd := exec.CommandContext(ctx, "sqlite3", s.dbPath, sql)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("sqlite3: %w: %s", err, string(b))
	}
	return string(b), nil
}
func esc(v string) string { return strings.ReplaceAll(v, "'", "''") }

func (s *Store) Migrate(ctx context.Context) error {
	_, err := s.exec(ctx, `PRAGMA journal_mode=WAL; PRAGMA busy_timeout=5000; PRAGMA foreign_keys=ON;
CREATE TABLE IF NOT EXISTS roots(id INTEGER PRIMARY KEY AUTOINCREMENT, path TEXT UNIQUE NOT NULL, added_at INTEGER NOT NULL, removed_at INTEGER);
CREATE TABLE IF NOT EXISTS nodes(id TEXT PRIMARY KEY, root_id INTEGER NOT NULL, abs_path TEXT NOT NULL, rel_path TEXT NOT NULL, kind TEXT NOT NULL, ext TEXT, size_bytes INTEGER NOT NULL, mtime_unix INTEGER NOT NULL, ctime_unix INTEGER, mode INTEGER NOT NULL, created_at INTEGER NOT NULL, updated_at INTEGER NOT NULL, deleted_at INTEGER, UNIQUE(root_id,abs_path));
CREATE TABLE IF NOT EXISTS scan_state(root_id INTEGER PRIMARY KEY, last_scan_started_at INTEGER, last_scan_completed_at INTEGER, last_scan_error TEXT, cursor_hint TEXT);
CREATE TABLE IF NOT EXISTS events_log(id INTEGER PRIMARY KEY AUTOINCREMENT, ts INTEGER NOT NULL, type TEXT NOT NULL, root_id INTEGER, abs_path TEXT, detail_json TEXT);
CREATE TABLE IF NOT EXISTS edges(id TEXT PRIMARY KEY, from_id TEXT NOT NULL, to_id TEXT NOT NULL, relation_type TEXT NOT NULL, note TEXT, created_at INTEGER NOT NULL, updated_at INTEGER NOT NULL, deleted_at INTEGER, UNIQUE(from_id,to_id,relation_type));
CREATE INDEX IF NOT EXISTS idx_nodes_root_rel ON nodes(root_id, rel_path);
CREATE INDEX IF NOT EXISTS idx_nodes_ext ON nodes(ext);
CREATE INDEX IF NOT EXISTS idx_nodes_updated ON nodes(updated_at);
CREATE INDEX IF NOT EXISTS idx_nodes_deleted ON nodes(deleted_at);
CREATE INDEX IF NOT EXISTS idx_edges_from_live ON edges(from_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_edges_to_live ON edges(to_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_edges_type_live ON edges(relation_type) WHERE deleted_at IS NULL;`)
	return err
}

func NormalizePath(p string) string { return filepath.Clean(p) }
func makeID(rootID int64, path string) string {
	h := fnv.New64a()
	_, _ = h.Write([]byte(strconv.FormatInt(rootID, 10) + "|" + path))
	return fmt.Sprintf("%x", h.Sum64())
}
func nullableInt(v *int64) string {
	if v == nil {
		return "NULL"
	}
	return strconv.FormatInt(*v, 10)
}

func (s *Store) AddRoot(ctx context.Context, path string) error {
	p := NormalizePath(path)
	_, err := s.exec(ctx, fmt.Sprintf("INSERT INTO roots(path,added_at) VALUES('%s',%d) ON CONFLICT(path) DO UPDATE SET removed_at=NULL;", esc(p), time.Now().Unix()))
	return err
}
func (s *Store) RemoveRoot(ctx context.Context, path string) error {
	_, err := s.exec(ctx, fmt.Sprintf("UPDATE roots SET removed_at=%d WHERE path='%s';", time.Now().Unix(), esc(NormalizePath(path))))
	return err
}
func (s *Store) ListRoots(ctx context.Context) ([]domain.Root, error) {
	out, err := s.exec(ctx, "SELECT id||'|'||path||'|'||added_at FROM roots WHERE removed_at IS NULL ORDER BY path;")
	if err != nil {
		return nil, err
	}
	roots := []domain.Root{}
	for _, ln := range strings.Split(strings.TrimSpace(out), "\n") {
		if ln == "" {
			continue
		}
		p := strings.SplitN(ln, "|", 3)
		id, _ := strconv.ParseInt(p[0], 10, 64)
		ts, _ := strconv.ParseInt(p[2], 10, 64)
		roots = append(roots, domain.Root{ID: id, Path: p[1], AddedAt: time.Unix(ts, 0)})
	}
	return roots, nil
}

func (s *Store) BulkUpsert(ctx context.Context, nodes []*domain.Node) error {
	for _, n := range nodes {
		if err := s.UpsertNode(ctx, n); err != nil {
			return err
		}
	}
	return nil
}
func (s *Store) UpsertNode(ctx context.Context, n *domain.Node) error {
	if n.ID == "" {
		n.ID = makeID(n.RootID, n.AbsPath)
	}
	now := time.Now().Unix()
	if n.CreatedAt == 0 {
		n.CreatedAt = now
	}
	if n.UpdatedAt == 0 {
		n.UpdatedAt = now
	}
	_, err := s.exec(ctx, fmt.Sprintf(`INSERT INTO nodes(id,root_id,abs_path,rel_path,kind,ext,size_bytes,mtime_unix,ctime_unix,mode,created_at,updated_at,deleted_at)
VALUES('%s',%d,'%s','%s','%s','%s',%d,%d,%s,%d,%d,%d,NULL)
ON CONFLICT(root_id,abs_path) DO UPDATE SET rel_path=excluded.rel_path,kind=excluded.kind,ext=excluded.ext,size_bytes=excluded.size_bytes,mtime_unix=excluded.mtime_unix,mode=excluded.mode,updated_at=excluded.updated_at,deleted_at=NULL;`,
		esc(n.ID), n.RootID, esc(n.AbsPath), esc(n.RelPath), n.Kind, esc(n.Ext), n.SizeBytes, n.MtimeUnix, nullableInt(n.CtimeUnix), n.Mode, n.CreatedAt, n.UpdatedAt))
	return err
}
func (s *Store) SoftDeleteByPath(ctx context.Context, rootID int64, absPath string) error {
	now := time.Now().Unix()
	_, err := s.exec(ctx, fmt.Sprintf("BEGIN; UPDATE nodes SET deleted_at=%d,updated_at=%d WHERE root_id=%d AND abs_path='%s'; UPDATE edges SET deleted_at=%d,updated_at=%d WHERE (from_id IN (SELECT id FROM nodes WHERE root_id=%d AND abs_path='%s') OR to_id IN (SELECT id FROM nodes WHERE root_id=%d AND abs_path='%s')) AND deleted_at IS NULL; COMMIT;", now, now, rootID, esc(absPath), now, now, rootID, esc(absPath), rootID, esc(absPath)))
	return err
}

func parseNode(line string) domain.Node {
	p := strings.Split(line, "|")
	rid, _ := strconv.ParseInt(p[1], 10, 64)
	sz, _ := strconv.ParseInt(p[6], 10, 64)
	mt, _ := strconv.ParseInt(p[7], 10, 64)
	md, _ := strconv.ParseUint(p[8], 10, 32)
	up, _ := strconv.ParseInt(p[9], 10, 64)
	return domain.Node{ID: p[0], RootID: rid, AbsPath: p[2], RelPath: p[3], Kind: domain.NodeKind(p[4]), Ext: p[5], SizeBytes: sz, MtimeUnix: mt, Mode: uint32(md), UpdatedAt: up}
}
func (s *Store) MoveNodePath(ctx context.Context, rootID int64, oldPath, newPath string) error {
	rel := filepath.Base(newPath)
	if roots, err := s.ListRoots(ctx); err == nil {
		for _, r := range roots {
			if r.ID == rootID {
				if rp, rerr := filepath.Rel(r.Path, newPath); rerr == nil {
					rel = filepath.ToSlash(rp)
				}
			}
		}
	}
	_, err := s.exec(ctx, fmt.Sprintf("UPDATE nodes SET abs_path='%s', rel_path='%s', ext='%s', updated_at=%d WHERE root_id=%d AND abs_path='%s' AND deleted_at IS NULL;", esc(newPath), esc(rel), esc(strings.ToLower(filepath.Ext(newPath))), time.Now().Unix(), rootID, esc(oldPath)))
	return err
}

func (s *Store) SearchFiles(ctx context.Context, q, ext string, rootID int64, kind string) ([]domain.Node, error) {
	where := "deleted_at IS NULL"
	if q != "" {
		lq := esc(strings.ToLower(q))
		where += fmt.Sprintf(" AND (lower(abs_path) LIKE '%%%s%%' OR lower(rel_path) LIKE '%%%s%%')", lq, lq)
	}
	if ext != "" {
		where += fmt.Sprintf(" AND ext='%s'", esc(strings.ToLower(ext)))
	}
	if rootID > 0 {
		where += fmt.Sprintf(" AND root_id=%d", rootID)
	}
	if kind != "" {
		where += fmt.Sprintf(" AND kind='%s'", esc(kind))
	}
	out, err := s.exec(ctx, "SELECT id||'|'||root_id||'|'||abs_path||'|'||rel_path||'|'||kind||'|'||ifnull(ext,'')||'|'||size_bytes||'|'||mtime_unix||'|'||mode||'|'||updated_at FROM nodes WHERE "+where+" ORDER BY updated_at DESC LIMIT 200;")
	if err != nil {
		return nil, err
	}
	res := []domain.Node{}
	for _, ln := range strings.Split(strings.TrimSpace(out), "\n") {
		if ln != "" {
			res = append(res, parseNode(ln))
		}
	}
	return res, nil
}
func (s *Store) GetFileByID(ctx context.Context, id string) (*domain.Node, error) {
	out, err := s.exec(ctx, fmt.Sprintf("SELECT id||'|'||root_id||'|'||abs_path||'|'||rel_path||'|'||kind||'|'||ifnull(ext,'')||'|'||size_bytes||'|'||mtime_unix||'|'||mode||'|'||updated_at FROM nodes WHERE id='%s' LIMIT 1;", esc(id)))
	if err != nil || strings.TrimSpace(out) == "" {
		return nil, err
	}
	n := parseNode(strings.TrimSpace(out))
	return &n, nil
}
func (s *Store) CountActiveFiles(ctx context.Context) (int, error) {
	out, err := s.exec(ctx, "SELECT COUNT(1) FROM nodes WHERE deleted_at IS NULL AND kind='file';")
	if err != nil {
		return 0, err
	}
	v, _ := strconv.Atoi(strings.TrimSpace(out))
	return v, nil
}
func (s *Store) RecentChanges(ctx context.Context, since int64) ([]domain.Node, error) {
	out, err := s.exec(ctx, fmt.Sprintf("SELECT id||'|'||root_id||'|'||abs_path||'|'||rel_path||'|'||kind||'|'||ifnull(ext,'')||'|'||size_bytes||'|'||mtime_unix||'|'||mode||'|'||updated_at FROM nodes WHERE updated_at >= %d ORDER BY updated_at DESC LIMIT 500;", since))
	if err != nil {
		return nil, err
	}
	res := []domain.Node{}
	for _, ln := range strings.Split(strings.TrimSpace(out), "\n") {
		if ln != "" {
			res = append(res, parseNode(ln))
		}
	}
	return res, nil
}

func (s *Store) CreateLink(ctx context.Context, fromID, toID, relationType, note string) error {
	if fromID == toID {
		return fmt.Errorf("self-link is not allowed")
	}
	note = strings.TrimSpace(note)
	if len(note) > 256 {
		return fmt.Errorf("note too long")
	}
	if relationType == "" {
		relationType = string(domain.RelationRelated)
	}
	if ok, err := s.nodeExists(ctx, fromID); err != nil || !ok {
		return fmt.Errorf("from node not indexed")
	}
	if ok, err := s.nodeExists(ctx, toID); err != nil || !ok {
		return fmt.Errorf("to node not indexed")
	}
	now := time.Now().Unix()
	edgeID := makeID(0, fromID+"|"+toID+"|"+relationType)
	_, err := s.exec(ctx, fmt.Sprintf("BEGIN; INSERT INTO edges(id,from_id,to_id,relation_type,note,created_at,updated_at,deleted_at) VALUES('%s','%s','%s','%s','%s',%d,%d,NULL) ON CONFLICT(from_id,to_id,relation_type) DO UPDATE SET note=excluded.note,updated_at=%d,deleted_at=NULL; INSERT INTO events_log(ts,type,abs_path,detail_json) VALUES(%d,'link_create','',json_object('from','%s','to','%s','type','%s')); COMMIT;", esc(edgeID), esc(fromID), esc(toID), esc(relationType), esc(note), now, now, now, now, esc(fromID), esc(toID), esc(relationType)))
	return err
}
func (s *Store) BatchCreateLinks(ctx context.Context, edges []domain.Edge) error {
	if len(edges) == 0 {
		return nil
	}
	var b strings.Builder
	now := time.Now().Unix()
	b.WriteString("BEGIN;")
	for _, e := range edges {
		if e.FromID == e.ToID || e.FromID == "" || e.ToID == "" {
			continue
		}
		r := string(e.RelationType)
		if r == "" {
			r = string(domain.RelationRelated)
		}
		edgeID := makeID(0, e.FromID+"|"+e.ToID+"|"+r)
		b.WriteString(fmt.Sprintf("INSERT INTO edges(id,from_id,to_id,relation_type,note,created_at,updated_at,deleted_at) VALUES('%s','%s','%s','%s','%s',%d,%d,NULL) ON CONFLICT(from_id,to_id,relation_type) DO UPDATE SET note=excluded.note,updated_at=%d,deleted_at=NULL;", esc(edgeID), esc(e.FromID), esc(e.ToID), esc(r), esc(e.Note), now, now, now))
	}
	b.WriteString("COMMIT;")
	_, err := s.exec(ctx, b.String())
	return err
}

func (s *Store) RemoveLinkByNodes(ctx context.Context, fromID, toID, relationType string) error {
	now := time.Now().Unix()
	q := fmt.Sprintf("UPDATE edges SET deleted_at=%d,updated_at=%d WHERE from_id='%s' AND to_id='%s'", now, now, esc(fromID), esc(toID))
	if relationType != "" {
		q += fmt.Sprintf(" AND relation_type='%s'", esc(relationType))
	}
	q += " AND deleted_at IS NULL;"
	_, err := s.exec(ctx, "BEGIN; "+q+" COMMIT;")
	return err
}
func (s *Store) RemoveLink(ctx context.Context, edgeID string) error {
	now := time.Now().Unix()
	_, err := s.exec(ctx, fmt.Sprintf("UPDATE edges SET deleted_at=%d,updated_at=%d WHERE id='%s' AND deleted_at IS NULL;", now, now, esc(edgeID)))
	return err
}

func (s *Store) ListLinks(ctx context.Context, nodeID string, direction domain.Direction) ([]domain.Edge, error) {
	where := "deleted_at IS NULL"
	switch direction {
	case domain.DirectionIn:
		where += fmt.Sprintf(" AND to_id='%s'", esc(nodeID))
	case domain.DirectionOut:
		where += fmt.Sprintf(" AND from_id='%s'", esc(nodeID))
	default:
		where += fmt.Sprintf(" AND (from_id='%s' OR to_id='%s')", esc(nodeID), esc(nodeID))
	}
	out, err := s.exec(ctx, "SELECT id||'|'||from_id||'|'||to_id||'|'||relation_type||'|'||ifnull(note,'')||'|'||created_at||'|'||updated_at FROM edges WHERE "+where+" ORDER BY created_at ASC;")
	if err != nil {
		return nil, err
	}
	edges := []domain.Edge{}
	for _, ln := range strings.Split(strings.TrimSpace(out), "\n") {
		if ln == "" {
			continue
		}
		p := strings.Split(ln, "|")
		ca, _ := strconv.ParseInt(p[5], 10, 64)
		ua, _ := strconv.ParseInt(p[6], 10, 64)
		edges = append(edges, domain.Edge{ID: p[0], FromID: p[1], ToID: p[2], RelationType: domain.RelationType(p[3]), Note: p[4], CreatedAt: ca, UpdatedAt: ua})
	}
	return edges, nil
}
func (s *Store) GetNeighbors(ctx context.Context, nodeID string) ([]domain.Edge, error) {
	return s.ListLinks(ctx, nodeID, domain.DirectionBoth)
}

func (s *Store) GraphStats(ctx context.Context) (domain.GraphStats, error) {
	stats := domain.GraphStats{ByType: map[string]int64{}}
	edges, err := s.exec(ctx, "SELECT COUNT(1) FROM edges WHERE deleted_at IS NULL;")
	if err != nil {
		return stats, err
	}
	stats.Edges, _ = strconv.ParseInt(strings.TrimSpace(edges), 10, 64)
	nodes, err := s.exec(ctx, "SELECT COUNT(DISTINCT id) FROM nodes WHERE id IN (SELECT from_id FROM edges WHERE deleted_at IS NULL UNION SELECT to_id FROM edges WHERE deleted_at IS NULL);")
	if err == nil {
		stats.LinkedNodes, _ = strconv.ParseInt(strings.TrimSpace(nodes), 10, 64)
	}
	rows, err := s.exec(ctx, "SELECT relation_type||'|'||COUNT(1) FROM edges WHERE deleted_at IS NULL GROUP BY relation_type ORDER BY relation_type;")
	if err == nil {
		for _, ln := range strings.Split(strings.TrimSpace(rows), "\n") {
			if ln == "" {
				continue
			}
			p := strings.Split(ln, "|")
			c, _ := strconv.ParseInt(p[1], 10, 64)
			stats.ByType[p[0]] = c
		}
	}
	return stats, nil
}

func (s *Store) ExportGraph(ctx context.Context) (domain.GraphExport, error) {
	res := domain.GraphExport{Version: 1, ExportedAt: time.Now().Unix(), Nodes: []domain.NodeExport{}, Edges: []domain.Edge{}}
	nodesOut, err := s.exec(ctx, "SELECT id||'|'||abs_path FROM nodes WHERE deleted_at IS NULL ORDER BY abs_path,id;")
	if err != nil {
		return res, err
	}
	for _, ln := range strings.Split(strings.TrimSpace(nodesOut), "\n") {
		if ln == "" {
			continue
		}
		p := strings.SplitN(ln, "|", 2)
		res.Nodes = append(res.Nodes, domain.NodeExport{ID: p[0], AbsPath: p[1]})
	}
	edges, err := s.ListLinks(ctx, "", domain.DirectionBoth)
	if err != nil {
		return res, err
	}
	// when nodeID empty, ListLinks returns all due BOTH where OR '' = false currently; fallback direct
	if len(edges) == 0 {
		rows, e2 := s.exec(ctx, "SELECT id||'|'||from_id||'|'||to_id||'|'||relation_type||'|'||ifnull(note,'')||'|'||created_at||'|'||updated_at FROM edges WHERE deleted_at IS NULL ORDER BY created_at ASC,id ASC;")
		if e2 == nil {
			for _, ln := range strings.Split(strings.TrimSpace(rows), "\n") {
				if ln == "" {
					continue
				}
				p := strings.Split(ln, "|")
				ca, _ := strconv.ParseInt(p[5], 10, 64)
				ua, _ := strconv.ParseInt(p[6], 10, 64)
				edges = append(edges, domain.Edge{ID: p[0], FromID: p[1], ToID: p[2], RelationType: domain.RelationType(p[3]), Note: p[4], CreatedAt: ca, UpdatedAt: ua})
			}
		}
	}
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].CreatedAt == edges[j].CreatedAt {
			return edges[i].ID < edges[j].ID
		}
		return edges[i].CreatedAt < edges[j].CreatedAt
	})
	res.Edges = edges
	return res, nil
}

func (s *Store) ImportGraph(ctx context.Context, in domain.GraphExport) (domain.ImportReport, error) {
	rep := domain.ImportReport{}
	pathMap := map[string]string{}
	for _, n := range in.Nodes {
		pathMap[n.ID] = n.ID
		if cur, _ := s.FindNodeIDByPath(ctx, n.AbsPath); cur != "" {
			pathMap[n.ID] = cur
		}
	}
	for _, e := range in.Edges {
		from := pathMap[e.FromID]
		to := pathMap[e.ToID]
		if from == "" || to == "" {
			rep.Skipped++
			continue
		}
		if err := s.CreateLink(ctx, from, to, string(e.RelationType), e.Note); err != nil {
			if strings.Contains(err.Error(), "self-link") || strings.Contains(err.Error(), "not indexed") {
				rep.Skipped++
			} else {
				rep.Duplicates++
			}
			continue
		}
		rep.Added++
	}
	return rep, nil
}

func (s *Store) FindNodeIDByPath(ctx context.Context, absPath string) (string, error) {
	out, err := s.exec(ctx, fmt.Sprintf("SELECT id FROM nodes WHERE abs_path='%s' AND deleted_at IS NULL LIMIT 1;", esc(absPath)))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}
func (s *Store) nodeExists(ctx context.Context, id string) (bool, error) {
	out, err := s.exec(ctx, fmt.Sprintf("SELECT COUNT(1) FROM nodes WHERE id='%s' AND deleted_at IS NULL;", esc(id)))
	if err != nil {
		return false, err
	}
	v, _ := strconv.Atoi(strings.TrimSpace(out))
	return v > 0, nil
}

func EncodeExport(exp domain.GraphExport) ([]byte, error) { return json.MarshalIndent(exp, "", "  ") }
func DecodeExport(raw []byte) (domain.GraphExport, error) {
	var e domain.GraphExport
	err := json.Unmarshal(raw, &e)
	return e, err
}

func (s *Store) NodeMeta(ctx context.Context, id string) (*domain.Node, error) {
	return s.GetFileByID(ctx, id)
}

func (s *Store) RecentEvents(ctx context.Context, limit int) ([]map[string]any, error) {
	if limit <= 0 {
		limit = 200
	}
	out, err := s.exec(ctx, fmt.Sprintf("SELECT ts||'|'||type||'|'||ifnull(abs_path,'')||'|'||ifnull(detail_json,'') FROM events_log ORDER BY id DESC LIMIT %d;", limit))
	if err != nil {
		return nil, err
	}
	rows := []map[string]any{}
	for _, ln := range strings.Split(strings.TrimSpace(out), "\n") {
		if ln == "" {
			continue
		}
		p := strings.SplitN(ln, "|", 4)
		ts, _ := strconv.ParseInt(p[0], 10, 64)
		rows = append(rows, map[string]any{"ts": ts, "type": p[1], "abs_path": p[2], "detail_json": p[3]})
	}
	return rows, nil
}
