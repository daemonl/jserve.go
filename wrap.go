package jserve

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.ecal.com/logger"
)

// UserError is an error to be encoded and sent to the API consumer
type UserError interface {
	HTTPStatus() int
	ResponseObject() interface{}
}

// Wrap returns a http handler for a given handler func. The handler returns an
// interface and error, Errors implementing UserError are shown to the user,
// other errors respond 500 and log.
// Responses implementing http.Handler will be chained to the handler, unmodified
// Otherwise, the response is written back the the response writer in JSON format
func Wrap(handler func(*http.Request) (interface{}, error)) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		response, err := handler(req)

		if err != nil {
			SendError(rw, req, err)
		} else {
			SendObject(rw, req, 200, response)
		}
	})
}

func SendError(rw http.ResponseWriter, req *http.Request, err error) {
	userError, ok := err.(UserError)
	if !ok {
		entry := logger.FromContext(req.Context())
		entry.WithField("error", err.Error()).
			Error("Unhandled server error")
		rw.WriteHeader(500)
		if _, err := rw.Write([]byte(`{"error":"Unknown Server Error"}`)); err != nil {
			entry.WithField("error", err.Error()).Error("Writing Resp")
		}
		return
	}
	rw.WriteHeader(userError.HTTPStatus())
	json.NewEncoder(rw).Encode(userError.ResponseObject())
	return

}

func SendObject(rw http.ResponseWriter, req *http.Request, status int, response interface{}) {

	if handlerResponse, ok := response.(http.Handler); ok {
		handlerResponse.ServeHTTP(rw, req)
		return
	}

	if response == nil {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	rw.WriteHeader(status)
	json.NewEncoder(rw).Encode(response)
}

type simpleError struct {
	Status  int    `json:"-"`
	Message string `json:"error"`
}

func (err simpleError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", err.Status, err.Message)
}
func (err simpleError) HTTPStatus() int {
	return err.Status
}
func (err simpleError) ResponseObject() interface{} {
	return err
}

// GenericUserError has only an 'error' field, set to the message.
func GenericUserError(status int, message string, params ...interface{}) error {
	return &simpleError{
		Status:  status,
		Message: fmt.Sprintf(message, params...),
	}
}

// JSONError returns a HTTP error for a JSON request parser error
func JSONError(err error) error {
	switch jErr := err.(type) {
	case *json.SyntaxError:
		return simpleError{
			Status:  http.StatusBadRequest,
			Message: jErr.Error(),
		}
	case *json.UnmarshalTypeError:
		return simpleError{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("Expecting %s for field %s, got %s", jErr.Type.Name(), jErr.Field, jErr.Value),
		}
	}
	return err
}
