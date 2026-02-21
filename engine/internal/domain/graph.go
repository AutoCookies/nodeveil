package domain

type RelationType string

const (
	RelationRelated   RelationType = "related"
	RelationRef       RelationType = "ref"
	RelationEvidence  RelationType = "evidence"
	RelationDependsOn RelationType = "depends_on"
)

type Direction string

const (
	DirectionOut  Direction = "OUT"
	DirectionIn   Direction = "IN"
	DirectionBoth Direction = "BOTH"
)

type Edge struct {
	ID           string       `json:"id"`
	FromID       string       `json:"from_id"`
	ToID         string       `json:"to_id"`
	RelationType RelationType `json:"relation_type"`
	Note         string       `json:"note,omitempty"`
	CreatedAt    int64        `json:"created_at"`
	UpdatedAt    int64        `json:"updated_at"`
	DeletedAt    *int64       `json:"deleted_at,omitempty"`
}

type GraphStats struct {
	LinkedNodes int64            `json:"linked_nodes"`
	Edges       int64            `json:"edges"`
	ByType      map[string]int64 `json:"by_type"`
}

type GraphExport struct {
	Version    int          `json:"version"`
	ExportedAt int64        `json:"exported_at"`
	Nodes      []NodeExport `json:"nodes"`
	Edges      []Edge       `json:"edges"`
}

type NodeExport struct {
	ID      string `json:"id"`
	AbsPath string `json:"abs_path"`
}

type ImportReport struct {
	Added      int `json:"added"`
	Skipped    int `json:"skipped"`
	Duplicates int `json:"duplicates"`
}

type GraphNeighborhood struct {
	Nodes      []Node  `json:"nodes"`
	Edges      []Edge  `json:"edges"`
	Truncated  bool    `json:"truncated"`
	NextCursor *string `json:"next_cursor"`
}

type GraphFilters struct {
	RelationTypes []string  `json:"relation_types"`
	Direction     Direction `json:"direction"`
	Extensions    []string  `json:"extensions"`
}

type GraphCounts struct {
	Out int64 `json:"out"`
	In  int64 `json:"in"`
}
