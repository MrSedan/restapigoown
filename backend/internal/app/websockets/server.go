package websockets

//Server is a websockets server
type Server struct {
	Hubs   map[string]*Hub
	NewHub chan *Hub
	RemHub chan *Hub
}

//NewServer create a server for websockets
func NewServer() *Server {
	return &Server{
		Hubs:   make(map[string]*Hub),
		NewHub: make(chan *Hub, 100),
		RemHub: make(chan *Hub, 100),
	}
}

//Run running server websocket
func (s *Server) Run() {
	for {
		select {
		case hub := <-s.NewHub:
			if _, ok := s.Hubs[hub.ID]; !ok {
				s.Hubs[hub.ID] = hub
			}
		case hub := <-s.RemHub:
			if _, ok := s.Hubs[hub.ID]; ok {
				delete(s.Hubs, hub.ID)
			}
		}

	}
}
