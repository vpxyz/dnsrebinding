/*
Package dnsrebinding provides simple middleware that protect your services against DNS rebinding attack.
In order to prevent DNS rebinding attack, this filter check if the Host header matches the host name of the server on which the resource resides.
This middleware increases the security level of CORS filter (see https://www.w3.org/TR/cors/#resource-security).
As default, if the provvided statusCode isn't valid, returns http.StatusNotImplemented.

Example:
package main

import (
  "net/http"
  dnsr "github.com/vpxyz/dnsrebinding"
)

func main() {
        dnsr.Filters(http.StatusNotFound, "example.com", "test.com", "test.me")

        handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                if r.Method == "GET" {
		           w.Header().Set("Content-Type", "application/json")
		           w.Write([]byte("{\"hello\": \"world\"}"))
		           return
	           }
        })
        http.ListenAndServe(":3000", dnsr(handler))
}
*/
package dnsrebinding

import (
	"net/http"
)

// Filter initialize the middleware, you must provide one or more hostname of the server on which yours services resides.
// This one can be usefull if your server has multiple hostnames.
func Filter(statusCode int, hostsName ...string) (fn func(next http.Handler) http.Handler) {
	if len(hostsName) == 0 {
		panic("You must provide one or more HostName, otherwise protection against DNS rebinding doesn't work")
	}

	for _, h := range hostsName {
		if len(h) == 0 {
			panic("You must provide no empty HostName, otherwise protection against DNS rebinding doesn't work")
		}
	}

	if http.StatusText(statusCode) == "" {
		// assume as default 501
		statusCode = http.StatusNotImplemented
	}

	if len(hostsName) == 1 {
		fn = func(next http.Handler) http.Handler {
			filter := func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Host") != hostsName[0] {
					w.WriteHeader(statusCode)

					// exit chain with status HTTP 501
					return
				}
				next.ServeHTTP(w, r)
			}
			return http.HandlerFunc(filter)
		}
		return fn
	}

	hn := make(map[string]bool, len(hostsName))

	for _, h := range hostsName {
		hn[h] = true
	}

	fn = func(next http.Handler) http.Handler {
		filter := func(w http.ResponseWriter, r *http.Request) {
			if _, ok := hn[r.Header.Get("Host")]; !ok {
				w.WriteHeader(statusCode)

				// exit chain with status HTTP 501
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(filter)
	}
	return fn
}
