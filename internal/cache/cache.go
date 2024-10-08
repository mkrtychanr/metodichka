package cache

import "context"

type Cache interface {
	Get(context.Context, string) (string, bool, error)
	Set(context.Context, string, string) error
	Close() error
}
