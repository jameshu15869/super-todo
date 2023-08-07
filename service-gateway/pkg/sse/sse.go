package sse

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SseServer struct {
	message       chan string
	newClients    chan chan string
	closedClients chan chan string
	totalClients  map[chan string]bool
}

type ClientChan chan string

type Message struct {
	MessageType string `json:"messageType"`
	Content     string `json:"content"`
}

func NewSseServer() *SseServer {
	sseServer := &SseServer{
		message:       make(chan string),
		newClients:    make(chan chan string),
		closedClients: make(chan chan string),
		totalClients:  make(map[chan string]bool),
	}

	go sseServer.listen()

	return sseServer
}

func (s *SseServer) listen() {
	for {
		select {
		case client := <-s.newClients:
			s.totalClients[client] = true
			log.Printf("Client added. %d registered clients", len(s.totalClients))
		case client := <-s.closedClients:
			delete(s.totalClients, client)
			close(client)
			log.Printf("Removed client. %d registered clients", len(s.totalClients))
		case msg := <-s.message:
			for clientChan := range s.totalClients {
				clientChan <- msg
			}
		}
	}
}

func (s *SseServer) InitializeConnection() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientChan := make(ClientChan)

		s.newClients <- clientChan

		defer func() {
			s.closedClients <- clientChan
		}()

		c.Set("clientChan", clientChan)

		c.Next()
	}
}

func (s *SseServer) PushMessage(message *Message) {
	stringMsg, err := json.Marshal(message)
	if err != nil {
		log.Printf("PushMessage - Failed to marshal to json: %v\n", err)
		return
	}
	s.message <- string(stringMsg)
}

func SseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		// c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func Stream() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("clientChan")
		if !ok {
			c.AbortWithError(http.StatusNotFound, errors.New("clientChan was null"))
			return
		}

		clientChan, ok := v.(ClientChan)
		if !ok {
			c.AbortWithError(http.StatusNotFound, errors.New("clientChan was null"))
			return
		}

		c.Stream(func(w io.Writer) bool {
			if msg, ok := <-clientChan; ok {
				c.SSEvent("message", msg)
				return true
			}
			return false
		})
	}
}
