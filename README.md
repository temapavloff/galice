# galice - golang SDK for Yandex Alice [![Build Status](https://travis-ci.com/temapavloff/galice.svg?branch=master)](https://travis-ci.com/temapavloff/galice) [![Go Report Card](https://goreportcard.com/badge/github.com/temapavloff/galice)](https://goreportcard.com/report/github.com/temapavloff/galice)

This library provides some basic helpers for Yandex Alice API:
- go structs for representing API entities;
- helper function for creating incoming webhook HTTP handler.

For more details on Alice API see [official documentation](https://yandex.ru/dev/dialogs/alice/) and package [godoc](https://godoc.org/github.com/temapavloff/galice).

> Note. Responses with image cards not supported in this SDK version!

## Installation

`go get -u github.com/temapavloff/galice`

## Usage

Setting up client:

```golang
c:= galice.New(
    // autoPings=true; automatically handle Alice API healthcheck requests (respond with pond to pings);
    // if false, you have to handle such requests in your own AliceHandler
    true,
    // autoDanderousContext=true; automatically handle requests marked as dangerous by Alice API.
    // Dangerous requests may contain suicide thoughts, hate speech, threats, etc.
    // If false, you have to handle such requests in your own AliceHandler
    true)

// Logger function will be called if some error occured while handling Alice API incoming request:
// bad requests, invalid responses, unexpected panics, etc.
// Default logger simply writes to stderr
c.SetLogger(func (err error) {
    fmt.Print(err)
})
```

Setting up request handler for simple skill which respond with user input message and closes session:

```golang
h := cli.CreateHandler(func(i InputData) (OutputData, error) {
    r := NewResponse(i.Request.OriginalUtterance, "", true)
    return NewOutput(i, r), nil
})

http.Handle("/skill", h)
log.Fatal(http.ListenAndServe(":8080", nil))
```

Adding buttons to response:

```golang
h := cli.CreateHandler(func(i InputData) (OutputData, error) {
    r := NewResponse("Hi!", "", true)
    r.AddButton("Link somewhere", false, "https://yandex.ru", nil) // Add Link button

    p := map[string]string{"key1": "val1", "key2": "val2"}
    r.AddButton("Payload", false, "", p) // Add button with some payload
    
    return NewOutput(i, r), nil
})

http.Handle("/skill", h)
log.Fatal(http.ListenAndServe(":8080", nil))
```
