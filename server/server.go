package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"mini-kafka/broker"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // allow all
}

type Server struct {
	Broker *broker.Broker
}

func (s *Server) Start(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/publish", s.handlePublish)
	mux.HandleFunc("/ws", s.handleWebSocket)
	log.Printf("HTTP server running at %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func (s *Server) handlePublish(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Topic string `json:"topic"`
		Data  string `json:"data"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	log.Printf("[PUBLISH] topic=%s data=%q", req.Topic, req.Data)
	s.Broker.Publish(req.Topic, req.Data)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	name := r.URL.Query().Get("name")
	if topic == "" {
		http.Error(w, "missing topic", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "upgrade failed", http.StatusInternalServerError)
		return
	}

	clientID := fmt.Sprintf("%s - %s", name, r.RemoteAddr)
	log.Printf("[WS CONNECT] %s subscribed to topic=%q", clientID, topic)

	ch := s.Broker.Subscribe(topic)
	defer func() {
		s.Broker.Unsubscribe(topic, ch)
		conn.Close()
		log.Printf("[WS DISCONNECT] %s unsubscribed from topic=%q", clientID, topic)
	}()

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			log.Printf("[SEND] to=%s topic=%s data=%q", clientID, topic, msg)
			if err = conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				return
			}
		case <-r.Context().Done():
			return
		}
	}
}
