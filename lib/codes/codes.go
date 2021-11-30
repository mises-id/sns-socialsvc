package codes

import (
	"fmt"
	"net/http"
)

// Code is a error of error code
type Code struct {
	HTTPStatus int    `json:"-"`
	Code       int    `json:"code"`
	Msg        string `json:"message"`
}

// Error return error message
func (e Code) Error() string {
	return e.Msg
}

// Equal return weather it is the same code
func (e Code) Equal(err error) bool {
	sourceErr, ok := err.(Code)
	if !ok {
		return false
	}

	return sourceErr.Code == e.Code
}

// New will modify the msg of Code
func (e Code) New(msg string) Code {
	return Code{
		HTTPStatus: e.HTTPStatus,
		Code:       e.Code,
		Msg:        msg,
	}
}

// Newf will modify the msg of Code with params
func (e Code) Newf(msg string, args ...interface{}) Code {
	return Code{
		HTTPStatus: e.HTTPStatus,
		Code:       e.Code,
		Msg:        fmt.Sprintf(msg, args...),
	}
}

const (
	SuccessCode             = 0
	InvalidArgumentCode     = 400000
	InvalidAuthCode         = 400001
	InvalidAuthMethodCode   = 400002
	InvalidAuthTokenCode    = 400003
	UnauthorizedCode        = 401000
	AuthorizeFailedCode     = 401001
	ForbiddenCode           = 403000
	UsernameExistedCode     = 403001
	TokenExpiredCode        = 403002
	NotFoundCode            = 404000
	UnprocessableEntityCode = 422000
	UsernameDuplicateCode   = 422001
	InternalCode            = 500000
	UnimplementedCode       = 500001
)

var (
	Success                = Code{HTTPStatus: http.StatusOK, Code: SuccessCode, Msg: "success"}
	ErrInvalidArgument     = Code{HTTPStatus: http.StatusBadRequest, Code: InvalidArgumentCode, Msg: "invalid params"}
	ErrInvalidAuth         = Code{HTTPStatus: http.StatusBadRequest, Code: InvalidAuthCode, Msg: "invalid auth params"}
	ErrInvalidAuthMethod   = Code{HTTPStatus: http.StatusBadRequest, Code: InvalidAuthMethodCode, Msg: "invalid auth method"}
	ErrInvalidAuthToken    = Code{HTTPStatus: http.StatusBadRequest, Code: InvalidAuthTokenCode, Msg: "invalid auth token"}
	ErrUnauthorized        = Code{HTTPStatus: http.StatusUnauthorized, Code: UnauthorizedCode, Msg: "unauthorized"}
	ErrAuthorizeFailed     = Code{HTTPStatus: http.StatusUnauthorized, Code: AuthorizeFailedCode, Msg: "authorize failed"}
	ErrForbidden           = Code{HTTPStatus: http.StatusForbidden, Code: ForbiddenCode, Msg: "forbidden"}
	ErrTokenExpired        = Code{HTTPStatus: http.StatusForbidden, Code: TokenExpiredCode, Msg: "authorization expired"}
	ErrUsernameExisted     = Code{HTTPStatus: http.StatusUnprocessableEntity, Code: UsernameExistedCode, Msg: "username had existed"}
	ErrNotFound            = Code{HTTPStatus: http.StatusNotFound, Code: NotFoundCode, Msg: "not found"}
	ErrUnprocessableEntity = Code{HTTPStatus: http.StatusUnprocessableEntity, Code: UnprocessableEntityCode, Msg: "unprocessable entity"}
	ErrUsernameDuplicate   = Code{HTTPStatus: http.StatusUnprocessableEntity, Code: UsernameDuplicateCode, Msg: "username duplicate"}
	ErrInternal            = Code{HTTPStatus: http.StatusInternalServerError, Code: InternalCode, Msg: "internal error"}
	ErrUnimplemented       = Code{HTTPStatus: http.StatusInternalServerError, Code: InternalCode, Msg: "internal error"}
)
