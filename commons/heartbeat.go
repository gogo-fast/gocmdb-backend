package commons

import "gogo-cmdb/apiserver/utils"

type HeartBeatMsg struct {
	UUID      string          `json:"uuid"`
	Timestamp utils.NullInt64 `json:"timestamp"`
}
