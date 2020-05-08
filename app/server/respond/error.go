package respond

import "net/http"

type BizError struct {
	Code int
	Msg  string
}

var (
	AlreadyConfirmed = BizError{
		Code: 1001,
		Msg:  "consensus order already confirmed",
	}
	NoRelation = BizError{
		Code: 1002,
		Msg:  "no relation between specified two user",
	}

	InternalServerError = BizError{
		Code: 500,
		Msg:  "internal server error",
	}
)

func BadRequest(msg string) BizError {
	return BizError{
		Code: http.StatusBadRequest,
		Msg:  msg,
	}
}
