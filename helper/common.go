package helper

import (
	"fmt"
	"github.com/vickydk/utl/config"
	"github.com/vickydk/utl/log"
	"net/http"
)

const (
	ErrPassAlreadyUsed = 1400
	ErrOARNotFound     = 1404
)

type Response struct {
	Content interface{} `json:"content"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
}

func MappingError(err error, resCode int) (int, string) {
	if config.Env.Debug && err != nil {
		switch resCode {
		case ErrPassAlreadyUsed:
			return http.StatusBadRequest, err.Error()
		case ErrOARNotFound:
			return http.StatusBadRequest, err.Error()
		default:
			return resCode, err.Error()
		}
	} else {
		switch resCode {
		case http.StatusOK:
			return resCode, ""
		case http.StatusPaymentRequired:
			return resCode, "Payment Required"
		case http.StatusBadRequest:
			return resCode, "Missing request"
		case http.StatusConflict:
			return resCode, "Data Conflict"
		case http.StatusNotFound:
			return resCode, "Data Not Found"
		case http.StatusInsufficientStorage:
			return resCode, "Error Database"
		case http.StatusForbidden:
			return resCode, "Access Denied"
		case http.StatusUnauthorized:
			return resCode, "User is not authorized"
		case ErrPassAlreadyUsed:
			return http.StatusBadRequest, "New password cannot be the same as your old password"
		case ErrOARNotFound:
			return http.StatusBadRequest, "Can't get any fulfillment store"
		default:
			return resCode, fmt.Sprint("Not yet mapped: ", resCode)
		}
	}
}

func Respond(err error, result interface{}, code int) *Response {
	if err != nil {
		log.Errorf("error: ", err)
	}
	rC, msg := MappingError(err, code)
	resp := &Response{
		Content: result,
		Code:    rC,
		Message: msg,
	}

	return resp
}
