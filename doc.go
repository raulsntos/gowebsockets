/*
Package gowebsockets provides an easy way to implement Web Sockets in your server

How to install

    Run go get github.com/raulsntos/gowebsockets

Example

A basic HTTPS server that uses this package under the /socket path:

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

In the example we have a FileServer to serve static files like any other regular Go server and then we create the WebSocket. Then we have some event listeners and then we use http.Handle to assign the WebSocket handler to the /socket path like we would with an http.HandlerFunc, it works very similar to an HTTP server and you can use the handler however you use http.HandlerFunc in your server. The WebSocket.Handler is a function of type websocket.Handler, more information in the package https://godoc.org/golang.org/x/net/websocket.

A basic javascript example to use with the WebSockets example:

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

Event listeners

There are four types of event listeners: WebSocket.On, WebSocket.OnConnect, WebSocket.OnDisconnect and WebSocket.OnError.

WebSocket.On takes a string (the name of the message, see the type gowebsockets.Message) and a gowebsockets.MessageListener (a method that takes a *gowebsockets.Client and a []byte). When you implement the gowebsockets.MessageListener you can use the Client which represents the connection to a client (a browser tab) to respond by sending messages (more on that below), and the []byte is the marshaled message content (gowebsockets.Message.Content), you would probably want to unmarshal the content in order to use it.

WebSocket.OnConnect takes a gowebsockets.ConnectionListener which is a method that takes a *gowebsockets.Client, it is fired when a client connects to the server.

WebSocket.OnDisconnect takes a gowebsockets.ConnectionListener which is a method that takes a *gowebsockets.Client, it is fired when a client disconnects from the server.

WebSocket.OnError takes a gowebsockets.ErrorListener which is a method that takes an error, it is fired when there is an error like failing to parse the content of a received message because it's not a gowebsockets.Message JSON encoded string or any other errors with the WebSocket.

Message

In order to send and receive messages you'll have to use the gowebsockets.Message type which is a struct that contains a string Name and an interface{} Content. The Name will be used by the WebSocket.On event listener to fire the appropiate gowebsockets.MessageListener but you can always use the same name and handle every message in the same method. The Content contains the message data, you can assign data of any type and it will be encoded as JSON to send it, when you receive a message the WebSocket.On gives you access to the Content of the Message in []bytes so you'll have to unmarshal it using the encoding/json package, more information about the package in https://golang.org/pkg/encoding/json/.

Client

The event listeners give you access to a *gowebsockets.Client struct that contains the connection to a client. You can use the Client.ID to identify the client, it's a randomly generated UUID so it's unique. You can use the client's methods to send messages. See the type gowebsockets.Client for information about each method.

Rooms

By default, every Client joins a room with the same ID as the Client.ID, this room is called the personal room and the client cannot leave it. You can create rooms by using the gowebsockets.Client.Join method passing the ID of a room that doesn't exist. See the type type gowebsockets.Client for information about the methods Join, Leave and In and how they allow you to use rooms for sending messages.
*/
package gowebsockets
