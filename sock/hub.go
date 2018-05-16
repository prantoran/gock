package sock

// Hub maintains the set of active clients and broadcasts messages to the
// clients.

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	send chan *MsgBody

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// UserIDs stores the list of userIDs to allow communications
	tokens map[string]bool

	userCache map[string]*Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		send:       make(chan *MsgBody),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		tokens:     make(map[string]bool),
		userCache:  make(map[string]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case body := <-h.send:
			client := h.userCache[body.Token]
			client.send <- []byte(body.Msg)
		case client := <-h.register:
			h.clients[client] = true
			h.userCache[client.Token] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
