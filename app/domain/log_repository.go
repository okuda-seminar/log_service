package domain

type ILogRepository interface {
	Save(log *Log) error
}
