package http

// responseErrorCode is an error code returned by the API
type responseErrorCode int

const (
	errorCodeUnknown responseErrorCode = iota
	errorCodeInvalidArgument
	errorCodeUnprocessableEntity
	errorCodeInternal
)

func (e responseErrorCode) String() string {
	switch e {
	case errorCodeUnknown:
		return "unknown"
	case errorCodeInvalidArgument:
		return "invalid_argument"
	case errorCodeInternal:
		return "internal_error"
	case errorCodeUnprocessableEntity:
		return "unprocessable_entity"
	default:
		return "unknown"
	}
}
