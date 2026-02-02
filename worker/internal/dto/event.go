package dto

import "time"

type Event struct {
	Key       []byte
	Value     []byte
	Timestamp time.Time
	Topic     string
	Partition int32
	Offset    int64
}
