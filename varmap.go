package jserve

import (
	"fmt"
	"net/http"
	"strconv"
)

// VarMap allows easy parsing of string types
type VarMap map[string]string

type ValidationError struct {
	Key   string
	Issue string
}

func (err ValidationError) Error() string {
	return fmt.Sprintf("Validation: Key %s - %s", err.Key, err.Issue)
}
func (err ValidationError) HTTPStatus() int {
	return http.StatusBadRequest
}
func (err ValidationError) ResponseObject() interface{} {
	return map[string]string{
		"error": fmt.Sprintf("%s %s", err.Key, err.Issue),
	}
}

type ParameterMissingError struct {
	Key string
}

func (err ParameterMissingError) Error() string {
	return fmt.Sprintf("Validation: Key %s is required", err.Key)
}
func (err ParameterMissingError) HTTPStatus() int {
	return http.StatusBadRequest
}
func (err ParameterMissingError) ResponseObject() interface{} {
	return map[string]string{
		"error": fmt.Sprintf("%s is required", err.Key),
	}
}

func (vm VarMap) String(key string) (string, error) {
	val, ok := vm[key]
	if !ok {
		return "", &ParameterMissingError{
			Key: key,
		}
	}
	return val, nil
}

func (vm VarMap) UInt64(key string) (uint64, error) {
	val, err := vm.String(key)
	if err != nil {
		return 0, err
	}
	intVal, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, &ValidationError{
			Key:   key,
			Issue: err.Error(),
		}
	}

	return intVal, nil
}
