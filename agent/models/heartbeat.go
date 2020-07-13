package models

type HeartBeatMsg struct {
	UUID      string    `json:"uuid"`
	Timestamp NullInt64 `json:"timestamp"`
}
