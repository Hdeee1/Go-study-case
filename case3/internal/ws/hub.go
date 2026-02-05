package ws

type Hub struct {
	RegisteredClients 	map[*Client]bool
	Broadcast 			chan []byte
	Register 			chan *Client
	Unregister			chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:		  make(chan []byte),
		Register:		  make(chan *Client),
		Unregister:		  make(chan *Client),
		RegisteredClients: make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.RegisteredClients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.RegisteredClients[client]; ok {
				delete(h.RegisteredClients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			for client := range h.RegisteredClients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.RegisteredClients, client)
				}
			}
		}
	}
}