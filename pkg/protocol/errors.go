package protocol

import "errors"

var UNKNOWN_TOPIC_OR_PARTITION = errors.New("This server does not host this topic-partition.")
