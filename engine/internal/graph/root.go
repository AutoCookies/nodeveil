package graph

import (
	"errors"
	"strings"
)

var ErrInvalidRootPath = errors.New("invalid root path")

type Root struct {
	Path string
}

func NewRoot(path string) (Root, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return Root{}, ErrInvalidRootPath
	}
	return Root{Path: trimmed}, nil
}
