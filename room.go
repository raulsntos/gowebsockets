package websockets

// Room contains clients that have joined the room
// ID is the identifier of the room
type Room struct {
	ID      string
	clients map[string]bool
}

// NewRoom creates a room with the given room ID
func NewRoom(id string) *Room {
	r := &Room{ID: id}
	r.clients = make(map[string]bool)
	return r
}

// IsClient returns if the client with the given ID is in the room
func (r *Room) IsClient(clientID string) bool {
	_, ok := r.clients[clientID]
	return ok
}

// GetClients returns an array with the IDs of every client that are currently in the room
func (r *Room) GetClients() []string {
	clients := []string{}
	for clientID := range r.clients {
		clients = append(clients, clientID)
	}
	return clients
}
