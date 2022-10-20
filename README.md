# :warning: Deprecated

This package is abandoned, please use a more actively maintained WebSockets package:

- https://godoc.org/github.com/gorilla/websocket
- https://godoc.org/nhooyr.io/websocket

# gowebsockets [![GoDoc](https://godoc.org/github.com/raulsntos/gowebsockets?status.svg)](https://godoc.org/github.com/raulsntos/gowebsockets)

The gowebsockets package uses the [x/net/websocket](https://godoc.org/golang.org/x/net/websocket) package and implements an easy way to implement a WebSocket in your server. It provides an easy way to listen to messages and handle them as well as sending messages.
It also implements rooms, you can make clients join and leave rooms as well as broadcast to the every client in the server or to a specific room.

## How to install

Run `go get github.com/raulsntos/gowebsockets`

## Usage

See the [documentation](https://godoc.org/github.com/raulsntos/gowebsockets).

## Example

Go server:

```go
package main

import (
    "log"
    "net/http"

    ws "github.com/raulsntos/gowebsockets"
)

func main() {
    h := http.FileServer(http.Dir("."))
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.Method, r.URL.Path)
        h.ServeHTTP(w, r)
    })

    webSocket := ws.NewWebSocket()

    webSocket.OnConnect(ws.ConnectionListener(func(c *ws.Client) {
        log.Printf("Client %s has connected\n", c.ID)
    }))

    webSocket.On("message", ws.MessageListener(func(c *ws.Client, msg []byte) {
        log.Println("Received message", strings(msg))
        c.Emit(ws.NewMessage("message", "Welcome from the Server!"), c.ID)
    }))

    webSocket.OnDisconnect(ws.ConnectionListener(func(c *ws.Client) {
        log.Printf("Client %s has disconnected\n", c.ID)
    }))

    webSocket.OnError(ws.ErrorListener(func(err error) {
        log.Printf("Error: %s\n", err.Error())
    }))

    log.Println("Serving in http://localhost:8080")

    http.Handle("/socket", webSocket.Handler)
    http.Handle("/", handler)
    http.ListenAndServe(":8080", nil)
}
```

Javascript client:

```javascript
var socket = new WebSocket('ws://localhost:8080/socket');
socket.onopen = (event) => {

    socket.onmessage = (event) => {
        console.log(event.data);
    };

    socket.send(JSON.stringify({
        name: 'message',
        content: 'hello world'
    }));
};
```
