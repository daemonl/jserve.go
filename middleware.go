package jserve

import "net/http"

func VersionMiddleware(appName, version string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("X-Version", version)
			rw.Header().Add("X-Application", appName)
			next.ServeHTTP(rw, req)
		})
	}
}
