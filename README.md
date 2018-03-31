DNSRebinding
============

DNSRebinding is a simple net/http middleware that protect yours services against DNS rebinding attack.

This middleware increases the security level of CORS filter (see https://www.w3.org/TR/cors/#resource-security).

The usage is very simple, just pass the host name of the server on which the resource resides, and the statusCode to return in case of dns rebinding attack.

As default, if the provvided statusCode isn't valid, returns http.StatusNotImplemented.

Example
-------

``` go

 package main

 import (
   "net/http"
   dnsr "github.com/vpxyz/dnsrebinding"
 )

 func main() {
         dnsr.Filter(http.StatusNotAcceptable, "example.com")

         handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                 if r.Method == "GET" {
 		           w.Header().Set("Content-Type", "application/json")
 		           w.Write([]byte("{\"hello\": \"world\"}"))
 		           return
 	           }
         })
         http.ListenAndServe(":3000", dnsr(handler))
 }
 
```

If your server has multiple hostnames:

``` go
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
```



