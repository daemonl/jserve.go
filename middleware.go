package jserve

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

var hostname string

func init() {
	hostname, _ = os.Hostname()
}

func VersionMiddleware(appName, version string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("X-Version", version)
			rw.Header().Add("X-Application", appName)
			rw.Header().Add("X-Hostname", hostname)
			next.ServeHTTP(rw, req)
		})
	}
}

func UpHandler(appName, version string) http.Handler {
	boot := time.Now()
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"application": appName,
			"version":     version,
			"uptime":      time.Now().Sub(boot).Seconds(),
			"hostname":    hostname,
		})
	})
}

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(404)
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"status": "Not Found",
			"path":   req.URL.Path,
			"method": req.Method,
		})
	})
}
