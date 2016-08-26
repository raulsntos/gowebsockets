package gowebsockets

// Room contains clients that have joined the room.
// ID is the identifier of the room.
// To use rooms to manage messages see the gowebsockets.Client type and the methods Join, Leave and In.
type Room struct {
	ID      string
	clients map[string]bool
}

// NewRoom creates a room with the given room ID.
// You won't use this method directly.
func NewRoom(id string) *Room {
	r := &Room{ID: id}
	r.clients = make(map[string]bool)
	return r
}

// IsClient returns if the client with the given ID is in the room.
func (r *Room) IsClient(clientID string) bool {
	_, ok := r.clients[clientID]
	return ok
}

// GetClients returns an array with the IDs of every client that are currently in the room.
func (r *Room) GetClients() []string {
	clients := []string{}
	for clientID := range r.clients {
		clients = append(clients, clientID)
	}
	return clients
}
