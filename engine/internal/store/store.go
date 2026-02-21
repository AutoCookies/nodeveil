package store

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nodeveil/nodeveil/engine/internal/domain"
	"hash/fnv"
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
	_, err := s.exec(ctx, `PRAGMA journal_mode=WAL; PRAGMA busy_timeout=5000;
CREATE TABLE IF NOT EXISTS roots(id INTEGER PRIMARY KEY AUTOINCREMENT, path TEXT UNIQUE NOT NULL, added_at INTEGER NOT NULL, removed_at INTEGER);
CREATE TABLE IF NOT EXISTS nodes(id TEXT PRIMARY KEY, root_id INTEGER NOT NULL, abs_path TEXT NOT NULL, rel_path TEXT NOT NULL, kind TEXT NOT NULL, ext TEXT, size_bytes INTEGER NOT NULL, mtime_unix INTEGER NOT NULL, ctime_unix INTEGER, mode INTEGER NOT NULL, created_at INTEGER NOT NULL, updated_at INTEGER NOT NULL, deleted_at INTEGER, UNIQUE(root_id,abs_path));
CREATE TABLE IF NOT EXISTS scan_state(root_id INTEGER PRIMARY KEY, last_scan_started_at INTEGER, last_scan_completed_at INTEGER, last_scan_error TEXT, cursor_hint TEXT);
CREATE TABLE IF NOT EXISTS events_log(id INTEGER PRIMARY KEY AUTOINCREMENT, ts INTEGER NOT NULL, type TEXT NOT NULL, root_id INTEGER, abs_path TEXT, detail_json TEXT);
CREATE INDEX IF NOT EXISTS idx_nodes_root_rel ON nodes(root_id, rel_path);
CREATE INDEX IF NOT EXISTS idx_nodes_ext ON nodes(ext);
CREATE INDEX IF NOT EXISTS idx_nodes_updated ON nodes(updated_at);
CREATE INDEX IF NOT EXISTS idx_nodes_deleted ON nodes(deleted_at);`)
	return err
}

func NormalizePath(p string) string { return filepath.Clean(p) }

func (s *Store) AddRoot(ctx context.Context, path string) error {
	p := NormalizePath(path)
	now := time.Now().Unix()
	_, err := s.exec(ctx, fmt.Sprintf("INSERT INTO roots(path,added_at) VALUES('%s',%d) ON CONFLICT(path) DO UPDATE SET removed_at=NULL;", esc(p), now))
	return err
}
func (s *Store) RemoveRoot(ctx context.Context, path string) error {
	p := NormalizePath(path)
	now := time.Now().Unix()
	_, err := s.exec(ctx, fmt.Sprintf("UPDATE roots SET removed_at=%d WHERE path='%s';", now, esc(p)))
	return err
}
func (s *Store) ListRoots(ctx context.Context) ([]domain.Root, error) {
	out, err := s.exec(ctx, "SELECT id||'|'||path||'|'||added_at FROM roots WHERE removed_at IS NULL ORDER BY path;")
	if err != nil {
		return nil, err
	}
	res := []domain.Root{}
	for _, ln := range strings.Split(strings.TrimSpace(out), "\n") {
		if ln == "" {
			continue
		}
		p := strings.SplitN(ln, "|", 3)
		id, _ := strconv.ParseInt(p[0], 10, 64)
		ts, _ := strconv.ParseInt(p[2], 10, 64)
		res = append(res, domain.Root{ID: id, Path: p[1], AddedAt: time.Unix(ts, 0)})
	}
	return res, nil
}

func (s *Store) BulkUpsert(ctx context.Context, nodes []*domain.Node) error {
	for i := 0; i < len(nodes); i += 400 {
		j := i + 400
		if j > len(nodes) {
			j = len(nodes)
		}
		if err := s.bulkUpsertChunk(ctx, nodes[i:j]); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) bulkUpsertChunk(ctx context.Context, nodes []*domain.Node) error {
	if len(nodes) == 0 {
		return nil
	}
	var b strings.Builder
	b.WriteString("BEGIN; ")
	now := time.Now().Unix()
	for _, n := range nodes {
		if n.ID == "" {
			n.ID = makeID(n.RootID, n.AbsPath)
		}
		if n.CreatedAt == 0 {
			n.CreatedAt = now
		}
		if n.UpdatedAt == 0 {
			n.UpdatedAt = now
		}
		b.WriteString(fmt.Sprintf(`INSERT INTO nodes(id,root_id,abs_path,rel_path,kind,ext,size_bytes,mtime_unix,ctime_unix,mode,created_at,updated_at,deleted_at)
VALUES('%s',%d,'%s','%s','%s','%s',%d,%d,%s,%d,%d,%d,NULL)
ON CONFLICT(root_id,abs_path) DO UPDATE SET rel_path=excluded.rel_path,kind=excluded.kind,ext=excluded.ext,size_bytes=excluded.size_bytes,mtime_unix=excluded.mtime_unix,mode=excluded.mode,updated_at=excluded.updated_at,deleted_at=NULL;`,
			esc(n.ID), n.RootID, esc(n.AbsPath), esc(n.RelPath), n.Kind, esc(n.Ext), n.SizeBytes, n.MtimeUnix, nullableInt(n.CtimeUnix), n.Mode, n.CreatedAt, n.UpdatedAt))
	}
	b.WriteString(" COMMIT;")
	_, err := s.exec(ctx, b.String())
	return err
}

func (s *Store) UpsertNode(ctx context.Context, n *domain.Node) error {
	if n.ID == "" {
		n.ID = makeID(n.RootID, n.AbsPath)
	}
	if n.CreatedAt == 0 {
		n.CreatedAt = time.Now().Unix()
	}
	if n.UpdatedAt == 0 {
		n.UpdatedAt = time.Now().Unix()
	}
	_, err := s.exec(ctx, fmt.Sprintf(`INSERT INTO nodes(id,root_id,abs_path,rel_path,kind,ext,size_bytes,mtime_unix,ctime_unix,mode,created_at,updated_at,deleted_at)
VALUES('%s',%d,'%s','%s','%s','%s',%d,%d,%s,%d,%d,%d,NULL)
ON CONFLICT(root_id,abs_path) DO UPDATE SET rel_path=excluded.rel_path,kind=excluded.kind,ext=excluded.ext,size_bytes=excluded.size_bytes,mtime_unix=excluded.mtime_unix,mode=excluded.mode,updated_at=excluded.updated_at,deleted_at=NULL;`,
		esc(n.ID), n.RootID, esc(n.AbsPath), esc(n.RelPath), n.Kind, esc(n.Ext), n.SizeBytes, n.MtimeUnix, nullableInt(n.CtimeUnix), n.Mode, n.CreatedAt, n.UpdatedAt))
	return err
}
func nullableInt(v *int64) string {
	if v == nil {
		return "NULL"
	}
	return strconv.FormatInt(*v, 10)
}
func (s *Store) SoftDeleteByPath(ctx context.Context, rootID int64, absPath string) error {
	now := time.Now().Unix()
	_, err := s.exec(ctx, fmt.Sprintf("UPDATE nodes SET deleted_at=%d, updated_at=%d WHERE root_id=%d AND abs_path='%s';", now, now, rootID, esc(absPath)))
	return err
}
func (s *Store) GetFileByID(ctx context.Context, id string) (*domain.Node, error) {
	out, err := s.exec(ctx, fmt.Sprintf("SELECT id||'|'||root_id||'|'||abs_path||'|'||rel_path||'|'||kind||'|'||ifnull(ext,'')||'|'||size_bytes||'|'||mtime_unix||'|'||mode||'|'||updated_at FROM nodes WHERE id='%s' LIMIT 1;", esc(id)))
	if err != nil || strings.TrimSpace(out) == "" {
		return nil, err
	}
	p := strings.Split(strings.TrimSpace(out), "|")
	rid, _ := strconv.ParseInt(p[1], 10, 64)
	sz, _ := strconv.ParseInt(p[6], 10, 64)
	mt, _ := strconv.ParseInt(p[7], 10, 64)
	md, _ := strconv.ParseUint(p[8], 10, 32)
	up, _ := strconv.ParseInt(p[9], 10, 64)
	return &domain.Node{ID: p[0], RootID: rid, AbsPath: p[2], RelPath: p[3], Kind: domain.NodeKind(p[4]), Ext: p[5], SizeBytes: sz, MtimeUnix: mt, Mode: uint32(md), UpdatedAt: up}, nil
}

func (s *Store) SearchFiles(ctx context.Context, q, ext string, rootID int64, kind string) ([]domain.Node, error) {
	where := "deleted_at IS NULL"
	if q != "" {
		where += fmt.Sprintf(" AND (lower(abs_path) LIKE '%%%s%%' OR lower(rel_path) LIKE '%%%s%%')", esc(strings.ToLower(q)), esc(strings.ToLower(q)))
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
	var res []domain.Node
	for _, ln := range strings.Split(strings.TrimSpace(out), "\n") {
		if ln == "" {
			continue
		}
		p := strings.Split(ln, "|")
		rid, _ := strconv.ParseInt(p[1], 10, 64)
		sz, _ := strconv.ParseInt(p[6], 10, 64)
		mt, _ := strconv.ParseInt(p[7], 10, 64)
		md, _ := strconv.ParseUint(p[8], 10, 32)
		up, _ := strconv.ParseInt(p[9], 10, 64)
		res = append(res, domain.Node{ID: p[0], RootID: rid, AbsPath: p[2], RelPath: p[3], Kind: domain.NodeKind(p[4]), Ext: p[5], SizeBytes: sz, MtimeUnix: mt, Mode: uint32(md), UpdatedAt: up})
	}
	return res, nil
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
	return s.SearchFiles(ctx, "", "", 0, "")
}

func makeID(rootID int64, path string) string {
	h := fnv.New64a()
	_, _ = h.Write([]byte(strconv.FormatInt(rootID, 10) + "|" + path))
	return fmt.Sprintf("%x", h.Sum64())
}
