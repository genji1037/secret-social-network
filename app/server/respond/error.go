package respond

import "net/http"

// BizError represent business error.
type BizError struct {
	Code int
	Msg  string
}

var (
	// AlreadyConfirmed represent consensus order already confirmed error.
	AlreadyConfirmed = BizError{
		Code: 1001,
		Msg:  "consensus order already confirmed",
	}
	// NoRelation represent no relation between specified two user error.
	NoRelation = BizError{
		Code: 1002,
		Msg:  "no relation between specified two user",
	}

	// InternalServerError represent internal server error.
	InternalServerError = BizError{
		Code: 500,
		Msg:  "internal server error",
	}
)

// BadRequest used to generate customer bad request error.
func BadRequest(msg string) BizError {
	return BizError{
		Code: http.StatusBadRequest,
		Msg:  msg,
	}
}
