package aal

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Server struct {
	findingsChan <-chan *Finding
	logger       *logrus.Logger
	upgrader     *websocket.Upgrader

	clients map[string]*websocket.Conn
}

func NewServer(logger *logrus.Logger, findingsChan <-chan *Finding) *Server {
	return &Server{
		findingsChan: findingsChan,
		logger:       logger,
		upgrader:     &websocket.Upgrader{},
		clients:      map[string]*websocket.Conn{},
	}
}

// serve an HTTP server that receives findings from the AAL system.
// The findings are received from the findings channel.
// The findings are served via a websocket connection.
func (s *Server) Run() error {
	go s.watchAndSendFindings()

	s.logger.Info("Serving findings on port 8080")
	http.HandleFunc("/ws", s.findingsHandler)
	return http.ListenAndServe(":8080", nil)
}

// findingsHandler handles the findings websocket connection.
func (s *Server) findingsHandler(w http.ResponseWriter, r *http.Request) {
	// Allow all origins
	s.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// Upgrade HTTP connection to websocket connection
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade HTTP connection to websocket connection:", err)
		return
	}
	s.clients[r.RemoteAddr] = conn

	s.logger.Info("New websocket client connected: ", r.RemoteAddr)
}

// watchAndSendFindings sends findings to all connected clients.
func (s *Server) watchAndSendFindings() {
	s.logger.Info("Waiting for findings")
	for finding := range s.findingsChan {
		s.logger.Debugf("Sending finding: %v", finding)
		for _, conn := range s.clients {
			err := conn.WriteJSON(finding)
			if err != nil {
				client := conn.RemoteAddr().String()
				s.logger.Errorf("Failed to send finding: %s. Removing client %s", err, client)
				delete(s.clients, client)
				continue
			}
			s.logger.Debugf("Sent finding: %s to %s", finding.Name, conn.RemoteAddr())
		}
	}
	s.logger.Info("Findings channel closed")
}
