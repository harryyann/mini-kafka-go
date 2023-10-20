package protocol

//go:generate stringer -type=ApiKey
type ApiKey int8

const (
	PRODUCE ApiKey = iota
	FETCH
	METADATA
	JOIN_GROUP
	SYNC_GROUP
	COMMIT_OFFSET
)
