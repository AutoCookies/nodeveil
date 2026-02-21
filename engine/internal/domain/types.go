package domain

import "time"

type NodeKind string

const (
	KindFile    NodeKind = "file"
	KindDir     NodeKind = "dir"
	KindSymlink NodeKind = "symlink"
)

type Root struct {
	ID        int64
	Path      string
	AddedAt   time.Time
	RemovedAt *time.Time
}

type Node struct {
	ID        string
	RootID    int64
	AbsPath   string
	RelPath   string
	Kind      NodeKind
	Ext       string
	SizeBytes int64
	MtimeUnix int64
	CtimeUnix *int64
	Mode      uint32
	CreatedAt int64
	UpdatedAt int64
	DeletedAt *int64
}

type IndexStatus struct {
	ProgressPercent float64 `json:"progressPercent"`
	CurrentRoot     string  `json:"currentRoot"`
	FilesIndexed    int64   `json:"filesIndexed"`
	QueueDepth      int     `json:"queueDepth"`
	Errors          int64   `json:"errors"`
	ScannedFiles    int64   `json:"scannedFilesTotal"`
	WatchEvents     int64   `json:"watchEventsTotal"`
	DBUpserts       int64   `json:"dbUpsertsTotal"`
	LastScanUnix    int64   `json:"lastScanUnix"`
}

type IndexEventType string

const (
	EventUpsert IndexEventType = "upsert"
	EventDelete IndexEventType = "delete"
	EventMove   IndexEventType = "move"
)

type IndexEvent struct {
	Type    IndexEventType
	RootID  int64
	Path    string
	OldPath string
	Node    *Node
}
