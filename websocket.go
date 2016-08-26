package gowebsockets

import "golang.org/x/net/websocket"

// WebSocket contains the websocket handler.
type WebSocket struct {
	Handler websocket.Handler

	rooms   map[string]*Room
	clients map[string]*Client

	sendCh      chan *Message
	broadcastCh chan *Message
	doneCh      chan *Client
	errCh       chan error

	listeners     map[string]MessageListener
	connListeners map[string]ConnectionListener
	errListener   ErrorListener
}

// MessageListener is a function that will be called for the message that it listens to.
type MessageListener func(c *Client, msg []byte)

// ConnectionListener is a function that will be called for connect/disconnect events.
type ConnectionListener func(c *Client)

// ErrorListener is a function that will be called for the event that it listens to.
type ErrorListener func(err error)

const bufferSize = 100

// NewWebSocket creates a WebSocket handler.
func NewWebSocket() *WebSocket {
	wss := &WebSocket{}

	// Channels
	wss.sendCh = make(chan *Message, bufferSize)
	wss.broadcastCh = make(chan *Message, bufferSize)
	wss.doneCh = make(chan *Client, bufferSize)
	wss.errCh = make(chan error, bufferSize)

	// Maps
	wss.rooms = make(map[string]*Room)
	wss.clients = make(map[string]*Client)

	// Listeners
	wss.listeners = make(map[string]MessageListener)
	wss.connListeners = make(map[string]ConnectionListener)
	wss.errListener = nil

	// Setup
	wss.setHandler()
	go wss.Listen()
	return wss
}

func (wss *WebSocket) addClient(c *Client) {
	wss.clients[c.ID] = c

	room := NewRoom(c.ID)
	wss.rooms[c.ID] = room

	c.Join(room.ID)

	listener, ok := wss.connListeners["connect"]
	if ok {
		listener(c)
	}
}

func (wss *WebSocket) deleteClient(c *Client) {
	listener, ok := wss.connListeners["disconnect"]
	if ok {
		listener(c)
	}

	for roomID := range c.rooms {
		wss.leaveRoom(c, roomID)
	}

	delete(wss.clients, c.ID)

	// TODO: Close Client and delete it
}

func (wss *WebSocket) joinRoom(c *Client, roomID string) {
	room, exists := wss.rooms[roomID]
	if !exists {
		room = NewRoom(roomID)
		wss.rooms[roomID] = room
	}
	room.clients[c.ID] = true
	c.rooms[roomID] = true
}

func (wss *WebSocket) leaveRoom(c *Client, roomID string) {
	delete(c.rooms, roomID)

	r := wss.rooms[roomID]
	delete(r.clients, c.ID)

	// If the room is empty, get rid of it
	if len(r.clients) == 0 {
		delete(wss.rooms, r.ID)
	}
}

func (wss *WebSocket) setHandler() {
	wss.Handler = websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		c := NewClient(ws, wss)
		wss.addClient(c)
		c.Listen()
	})
}

// Listen waits for messages to send to the connected clients.
// You won't use this method directly.
func (wss *WebSocket) Listen() {
	for {
		select {

		case msg := <-wss.sendCh:
			c, ok := wss.clients[msg.to]
			if ok {
				websocket.JSON.Send(c.ws, msg)
			}

		case msg := <-wss.broadcastCh:
			clientsSent := make(map[string]bool)
			clientsSent[msg.from] = true
			for _, roomID := range msg.rooms {
				for clientID := range wss.rooms[roomID].clients {
					_, sent := clientsSent[clientID]
					if !sent {
						c := wss.clients[clientID]
						websocket.JSON.Send(c.ws, msg)
						clientsSent[clientID] = true
					}
				}
			}

		case err := <-wss.errCh:
			if wss.errListener != nil {
				wss.errListener(err)
			}

		case c := <-wss.doneCh:
			wss.deleteClient(c)
		}
	}
}

// On subscribes a WSListener function to a specific event.
func (wss *WebSocket) On(event string, fn MessageListener) {
	wss.listeners[event] = fn
}

// OnConnect subscribes a channel function to the connect event.
func (wss *WebSocket) OnConnect(fn ConnectionListener) {
	wss.connListeners["connect"] = fn
}

// OnDisconnect subscribes a channel function to the error event.
func (wss *WebSocket) OnDisconnect(fn ConnectionListener) {
	wss.connListeners["disconnect"] = fn
}

// OnError subscribes a channel function to the error event.
func (wss *WebSocket) OnError(fn ErrorListener) {
	wss.errListener = fn
}

// GetClient checks if the client exists and returns the client and a boolean that can be used to check if the client exists (since it might be nil).
func (wss *WebSocket) GetClient(clientID string) (*Client, bool) {
	c, ok := wss.clients[clientID]
	return c, ok
}

// GetClients returns every client connected at the time.
func (wss *WebSocket) GetClients() []string {
	clients := []string{}
	for clientID := range wss.clients {
		clients = append(clients, clientID)
	}
	return clients
}

// GetRoom checks if the room exists and returns the room and a boolean that can be used to check if the room exists (since it might be nil).
func (wss *WebSocket) GetRoom(roomID string) (*Room, bool) {
	r, ok := wss.rooms[roomID]
	return r, ok
}

// GetRooms returns every room that currently exists.
func (wss *WebSocket) GetRooms() []string {
	rooms := []string{}
	for roomID := range wss.rooms {
		rooms = append(rooms, roomID)
	}
	return rooms
}
