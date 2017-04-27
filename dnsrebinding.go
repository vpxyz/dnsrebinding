/*
Package dnsrebinding provides simple middleware that protect your services against DNS rebinding attack.
In order to prevent DNS rebinding attack, this filter check if the Host header matches the host name of the server on which the resource resides.
This middleware increases the security level of CORS filter (see https://www.w3.org/TR/cors/#resource-security).

Example:
package main

import (
  "net/http"
  dnsr "github.com/vpxyz/dnsrebinding"
)

func main() {
        dnsr.Filter("example.com")

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

// Filter initialize the middleware, you must provide the HostName of the server on which yours services resides
func Filter(hostName string) func(next http.Handler) http.Handler {
	if len(hostName) == 0 {
		panic("You must provide the HostName, otherwise protection against DNS rebinding doesn't work")
	}

	return func(next http.Handler) http.Handler {
		filter := func(w http.ResponseWriter, r *http.Request) {
			// protection against DNS rebinding
			if r.Header.Get("Host") != hostName {
				// respond with 501 error code
				w.WriteHeader(http.StatusNotImplemented)
				// exit chain with status HTTP 501
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(filter)
	}
}
