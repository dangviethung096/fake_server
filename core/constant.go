package core

const (
	BLANK                 = ""
	CONTENT_TYPE_TEXT     = "text/plain"
	HTTP_CLIENT_TIMEOUT   = 0
	DEFAULT_CONSUMER_TAG  = "default_consumer"
	CONTENT_TYPE_KEY      = "Content-Type"
	JSON_CONTENT_TYPE     = "application/json"
	FORMDATA_CONTENT_TYPE = "application/x-www-form-urlencoded"
)

// Kind of message queue:
// "direct", "fanout", "topic" and "headers".
const (
	MESSAGE_QUEUE_KIND_DIRECT  = "direct"
	MESSAGE_QUEUE_KIND_FANOUT  = "fanout"
	MESSAGE_QUEUE_KIND_TOPIC   = "topic"
	MESSAGE_QUEUE_KIND_HEADERS = "headers"
)

// Error code
const (
	ERROR_CODE_READ_BODY_REQUEST_FAIL  = 100
	ERROR_CODE_CLOSE_BODY_REQUEST_FAIL = 101
	ERROR_BAD_BODY_REQUEST             = 102
)
