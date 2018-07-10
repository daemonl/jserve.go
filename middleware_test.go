package jserve

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {

	var h http.Handler = http.NotFoundHandler()
	h = VersionMiddleware("appName", "versionString")(h)

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/thing", nil)

	h.ServeHTTP(rw, req)

	if got := rw.Header().Get("X-Version"); got != "versionString" {
		t.Errorf("Bad Version header: %s", got)
	}
	if got := rw.Header().Get("X-Application"); got != "appName" {
		t.Errorf("Bad Application header: %s", got)
	}

}
