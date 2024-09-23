package domain

import "context"

type ILogRepository interface {
	Save(ctx context.Context, log *Log) error
}
