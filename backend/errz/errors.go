package errz

import "net/http"

type Type int8

const (
	Invalid Type = iota
	NotFound
	Unauthorized
)

func New(typ Type, msg string, err error) *Error {
	return &Error{
		msg: msg,
		err: err,
		typ: typ,
	}
}

type Error struct {
	msg string
	err error
	typ Type
}

func (e Error) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

func (e Error) Type() Type {
	return e.typ
}

func (e Error) HttpStatusText() string {
	return e.msg
}

func (e Error) HttpStatusCode() int {
	switch e.typ {
	case Invalid:
		return http.StatusBadRequest
	case NotFound:
		return http.StatusNotFound
	case Unauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
