package aal

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Server struct {
	findingsChan     <-chan *Finding
	observationsChan <-chan *Observation
	logger           *logrus.Logger
	upgrader         *websocket.Upgrader

	// clients is a map of all connected clients.
	clients      map[string]*websocket.Conn
	clientsMutex sync.Mutex
}

type Package struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func NewServer(logger *logrus.Logger, findingsChan <-chan *Finding, observationsChan <-chan *Observation) *Server {
	return &Server{
		findingsChan:     findingsChan,
		observationsChan: observationsChan,
		logger:           logger,
		upgrader:         &websocket.Upgrader{},
		clients:          map[string]*websocket.Conn{},
		clientsMutex:     sync.Mutex{},
	}
}

// serve an HTTP server that receives findings from the AAL system.
// The findings are received from the findings channel.
// The findings are served via a websocket connection.
func (s *Server) Run() error {
	go s.watchAndSendFindings()
	go s.watchAndSendObservations()

	s.logger.Info("Serving findings on port 8080")
	http.HandleFunc("/ws", s.websocketHandler)
	return http.ListenAndServe(":8080", nil)
}

// websocketHandler handles the websocket connection.
func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	// Allow all origins
	s.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// Upgrade HTTP connection to websocket connection
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade HTTP connection to websocket connection:", err)
		return
	}
	s.clientsMutex.Lock()
	s.clients[r.RemoteAddr] = conn
	s.clientsMutex.Unlock()

	s.logger.Info("New websocket client connected: ", r.RemoteAddr)
}

// watchAndSendObservations sends observations to all connected clients.
func (s *Server) watchAndSendObservations() {
	s.logger.Info("Waiting for observations")
	for observation := range s.observationsChan {
		s.logger.Debugf("Sending observation %s to %d clients", observation.Name, len(s.clients))
		for _, conn := range s.clients {
			err := conn.WriteJSON(&Package{
				Type: "observation",
				Data: observation,
			})
			if err != nil {
				client := conn.RemoteAddr().String()
				s.logger.Warnf("Removing client %s - %s", client, err.Error())
				s.clientsMutex.Lock()
				conn.Close()
				delete(s.clients, client)
				s.clientsMutex.Unlock()
				continue
			}
			s.logger.Debugf("Sent observation: %s to %s", observation.Name, conn.RemoteAddr())
		}
	}
	s.logger.Info("Observations channel closed")
}

// watchAndSendFindings sends findings to all connected clients.
func (s *Server) watchAndSendFindings() {
	s.logger.Info("Waiting for findings")
	for finding := range s.findingsChan {
		s.logger.Debugf("Sending finding %s to %d clients", finding.Name, len(s.clients))
		for _, conn := range s.clients {
			err := conn.WriteJSON(&Package{
				Type: "finding",
				Data: finding,
			})
			if err != nil {
				client := conn.RemoteAddr().String()
				s.logger.Warnf("Removing client %s - %s", client, err.Error())
				s.clientsMutex.Lock()
				conn.Close()
				delete(s.clients, client)
				s.clientsMutex.Unlock()
				continue
			}
			s.logger.Debugf("Sent finding: %s to %s", finding.Name, conn.RemoteAddr())
		}
	}
	s.logger.Info("Findings channel closed")
}
