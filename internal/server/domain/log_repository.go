package domain

import "context"

type ILogRepository interface {
	Save(ctx context.Context, log *Log) error
	CTRSave(ctx context.Context, ctrLog *CTRLog) error
	List(ctx context.Context) ([]Log, error)
}
