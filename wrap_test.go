package jserve

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWrapperHappy(t *testing.T) {

	apiFunc := func(req *http.Request) (interface{}, error) {
		return map[string]interface{}{"key": 123}, nil
	}

	handler := Wrap(apiFunc)
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://test.com/path", nil)
	handler.ServeHTTP(rw, req)

	if rw.Code != 200 {
		t.Errorf("Code %d", rw.Code)
	}
	if bodyString := strings.TrimSpace(rw.Body.String()); bodyString != `{"key":123}` {
		t.Errorf("Bad body: '%s'", bodyString)
	}

}

func TestWrapperErrors(t *testing.T) {

	for _, testCase := range []struct {
		apiFunc func(*http.Request) (interface{}, error)
		status  int
		body    string
	}{{
		apiFunc: func(req *http.Request) (interface{}, error) {
			return nil, fmt.Errorf("Unknown Error")
		},
		status: 500,
		body:   `{"error":"Unknown Server Error"}`,
	}, {
		apiFunc: func(req *http.Request) (interface{}, error) {
			return nil, GenericUserError(400, "Special Error Message")
		},
		status: 400,
		body:   `{"error":"Special Error Message"}`,
	}} {
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://test.com/path", nil)
		Wrap(testCase.apiFunc).ServeHTTP(rw, req)

		if rw.Code != testCase.status {
			t.Errorf("Expect Code %d, got %d", testCase.status, rw.Code)
		}
		if bodyString := strings.TrimSpace(rw.Body.String()); bodyString != testCase.body {
			t.Errorf("Bad body: \n Want %s\n  Got %s", testCase.body, bodyString)
		}
	}

}

func TestJSONErrors(t *testing.T) {

	extractErrorMessage := func(input interface{}) (string, error) {
		jByte, err := json.Marshal(input)
		if err != nil {
			return "", err
		}
		output := map[string]string{}
		if err := json.Unmarshal(jByte, &output); err != nil {
			return "", err
		}
		return output["error"], nil
	}
	into := struct {
		Key int `json:"key"`
	}{}

	for _, testCase := range []struct {
		input       string
		errContains string
	}{{
		input:       `BadJSON`,
		errContains: "invalid character 'B'",
	}, {
		input:       `{"key":"string"}`,
		errContains: "int for field key, got string",
	}} {

		if err := JSONError(json.Unmarshal([]byte(testCase.input), &into)); err == nil {
			t.Errorf("No error")
		} else if uErr, ok := err.(UserError); !ok {
			t.Errorf("Not a user error")
		} else if code := uErr.HTTPStatus(); code != 400 {
			t.Errorf("Code %d, not 400", code)
		} else if msg, err := extractErrorMessage(uErr.ResponseObject()); err != nil {
			t.Errorf("Marshal error: %s", err.Error())
		} else if !strings.Contains(msg, testCase.errContains) {
			t.Errorf("Bad error message: %s", msg)
		}
	}

}
