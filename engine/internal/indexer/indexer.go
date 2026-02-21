package indexer

import "context"

type Indexer interface {
	Start(context.Context) error
	Stop(context.Context) error
}
