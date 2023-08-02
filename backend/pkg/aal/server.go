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
}

func NewServer(logger *logrus.Logger, findingsChan <-chan *Finding) *Server {
	return &Server{
		findingsChan: findingsChan,
		logger:       logger,
		upgrader:     &websocket.Upgrader{},
	}
}

// serve an HTTP server that receives findings from the AAL system.
// The findings are received from the findings channel.
// The findings are served via a websocket connection.
func (s *Server) Run() error {
	s.logger.Info("Serving findings on port 8080")

	http.HandleFunc("/ws", s.findingsHandler)
	return http.ListenAndServe(":8080", nil)
}

// findingsHandler handles the findings websocket connection.
func (s *Server) findingsHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("New websocket client connected: ", r.RemoteAddr)

	// Upgrade HTTP connection to websocket connection
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade HTTP connection to websocket connection:", err)
		return
	}
	defer conn.Close()

	// Send findings to websocket connection
	for finding := range s.findingsChan {
		s.logger.Debugf("Sending finding: %v", finding)
		err := conn.WriteJSON(finding)
		if err != nil {
			s.logger.Errorf("Failed to send finding: %s", err)
			return
		}
	}
}
