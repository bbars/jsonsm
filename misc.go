package jsonsm

import (
	"encoding/json"
)

////////////////////////////////

type H map[string]interface{}

////////////////////////////////

type Error struct {
	httpCode int
	message string
	extras H
}


func NewError(httpCode int, message string, extras H) *Error {
	return &Error{
		httpCode: httpCode,
		message: message,
		extras: extras,
	}
}

func WrapError(err error, extras H) *Error {
	return &Error{
		httpCode: 400,
		message: err.Error(),
		extras: extras,
	}
}

func (this *Error) Error() string {
	if this == nil {
		return ""
	}
	return this.message
}

func (this *Error) HttpCode() int {
	if this == nil || this.httpCode <= 0 {
		return 400
	}
	return this.httpCode
}

func (this *Error) MarshalJSON() ([]byte, error) {
	res := H{
		"err_code": this.message,
	}
	if this.extras != nil {
		for k, v := range this.extras {
			res[k] = v
		}
	}
	return json.Marshal(res)
}
