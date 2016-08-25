package websockets

import (
	"errors"
	"io"

	"golang.org/x/net/websocket"
)

// Client represents a connection to the websocket
// ID is the identifier of the client
type Client struct {
	ID          string
	rooms       map[string]bool
	currentRoom string
	ws          *websocket.Conn
	wss         *WebSocket
}

// NewClient creates a Client
func NewClient(ws *websocket.Conn, wss *WebSocket) *Client {
	c := &Client{ID: generateUUID(), ws: ws, wss: wss, currentRoom: ""}
	c.rooms = make(map[string]bool)
	return c
}

// Listen receives messages from the connected client
func (c *Client) Listen() {
	for {
		var msg *Message
		err := websocket.JSON.Receive(c.ws, &msg)
		if err == io.EOF {
			// TODO: Does this means the client disconnected? NO??!
			c.wss.doneCh <- true
		} else if err != nil {
			c.wss.errCh <- err
		} else {
			c.receive(msg)
		}
	}
}

// GetRooms returns an array with the IDs of every room the client is currently in
func (c *Client) GetRooms() []string {
	rooms := []string{}
	for roomID := range c.rooms {
		rooms = append(rooms, roomID)
	}
	return rooms
}

// Join adds the client to a room by ID
func (c *Client) Join(roomID string) {
	c.wss.joinRoom(c, roomID)
}

// Leave removes the client from a room by ID
func (c *Client) Leave(roomID string) error {
	if c.ID == roomID {
		return errors.New("Cannot leave the personal room")
	}
	c.wss.leaveRoom(c, roomID)
	return nil
}

// Emit sends a message to the client
func (c *Client) Emit(msg *Message, clientID string) {
	msg.from = c.ID
	msg.to = clientID
	c.wss.sendCh <- msg
}

// Broadcast sends a message to every room (or every)
func (c *Client) Broadcast(msg *Message) {
	if c.currentRoom == "" {
		rooms := []string{}
		for roomID := range c.wss.rooms {
			rooms = append(rooms, roomID)
		}
		msg.rooms = rooms
	} else {
		msg.rooms = []string{c.currentRoom}
	}

	msg.from = c.ID
	msg.to = ""
	c.wss.broadcastCh <- msg
}

// In sets the client to send a broadcast to a specific room
// Use it always chaining a Broadcast method like: c.In(roomID).Broadcast(msg)
func (c *Client) In(roomID string) *Client {
	c.currentRoom = roomID
	return c
}

func (c *Client) receive(msg *Message) {
	listener, ok := c.wss.listeners[msg.Name]
	if ok {
		bytes, err := msg.toBytes()
		if err != nil {
			c.wss.errCh <- err
		} else {
			listener(c, bytes)
		}
	}
}
