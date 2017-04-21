[![GoDoc](https://godoc.org/github.com/casperin/calm?status.svg)](http://godoc.org/github.com/casperin/calm)
![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)
[![CircleCI](https://circleci.com/gh/casperin/calm.svg?style=svg)](https://circleci.com/gh/casperin/calm)

# Calm - a rate limiter

`Calm` is a rate limiter for http handler functions in Go.

## Quick intro

We'll set up a server that accepts 3 requests per second on `"/"`:
```go
package main

import (
    "fmt"
    "net/http"

	"github.com/casperin/calm"
)

var calmer = calm.New(3, time.Second)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there")
}

func main() {
    http.HandleFunc("/", calmer(handler))
    http.ListenAndServe(":8080", nil)
}
```

This will allow 3 requests per second from a user through to your `HandlerFunc`
and send a `429` to all other requests.

## More features

You can also define your own rate handle function (you can find the default
function in `/limiter`), lookup function (the function that converts a request
to an ip (or any string really), as well as define which method types should be
rate limited.
```go
// Setting the handle func
var calmer = calm.New(
    3,
    time.Second,
    calm.RateHandler(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Calm down friend."))
    }),
)

// Setting the lookup func
var calmer = calm.New(
    3,
    time.Second,
    calm.Lookup(func(r *http.Request) {
        // probably not how you want to authorize users, but I digress...
        return r.FormValue("userid")
    }),
)

// Setting the methods
var calmer = calm.New(
    3,
    time.Second,
    calm.Methods("POST", "PUT"), // only POST and PUT will be rate limited
)
```

To do more than one thing, you just pass in any number of them. Order is
irrelevant.
```go
// Setting the handle func and the methods
var calmer = calm.New(
    3,
    time.Second,
    calm.RateHandler(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Calm down friend."))
    }),
    calm.Methods("POST", "PUT"),
)
```

## Shoutout

Shoutout to Didip's [tollbooth](https://github.com/didip/tollbooth) which this
library borrowed a lot from. :)
