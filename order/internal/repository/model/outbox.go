package model

import "time"

type OutboxMessage struct {
	EventUUID string
	Topic     string
	Key       string
	Payload   []byte
	Status    OutboxStatus
	CreatedAt time.Time
}

type OutboxStatus string

const (
	OutboxStatusPending   OutboxStatus = "PENDING"
	OutboxStatusPublished OutboxStatus = "PUBLISHED"
)
