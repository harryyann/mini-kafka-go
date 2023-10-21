package protocol

import "errors"

var INVALID_REQUEST = errors.New("the request is invalid")
var UNKNOWN_TOPIC_OR_PARTITION = errors.New("this server does not host this topic-partition")
