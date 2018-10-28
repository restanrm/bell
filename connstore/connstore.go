package connstore

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/twinj/uuid"

	"github.com/gorilla/websocket"
)

// ConnStore is the struct that holds clients
type ConnStore struct {
	store     map[string]*websocket.Conn
	mu        sync.RWMutex
	interrupt chan os.Signal
}

// New return a new Client object
func New() *ConnStore {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	client := &ConnStore{
		store:     make(map[string]*websocket.Conn),
		interrupt: interrupt,
	}
	go client.run()
	return client
}

func (c *ConnStore) run() {
	done := make(chan struct{})
	for {
		select {
		case <-done:
			return
		case <-c.interrupt:
			logrus.Infof("Interrupt signal received, closing all clients connection")
			c.mu.Lock()
			defer c.mu.Unlock()
			for k, v := range c.store {
				logrus.Infof("Closing connection for client %v", k)
				err := v.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, "Server is closing down"))
				if err != nil {
					logrus.WithError(err).Errorf("Failed to send closing message to client %v", k)
					continue
				}
				err = v.Close()
				if err != nil {
					logrus.WithError(err).Errorf("Failed to close connection with client %v", k)
					continue
				}
			}
			logrus.Infof("All connections have been closed, quitting")
			os.Exit(0)
		}
	}
}

// Get retrieve the websocket connection and give it to the caller
func (c *ConnStore) Get(name string) (*websocket.Conn, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	conn, ok := c.store[name]
	if !ok {
		return nil, fmt.Errorf("Client %q isn't registered", name)
	}
	return conn, nil
}

// Register function is the public handler to associate new websockets store to the service.Register.
func (c *ConnStore) Register(name string, conn *websocket.Conn) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if name != "" {
		if _, ok := c.store[name]; ok {
			// Error, client is already registered
			return fmt.Errorf("client %q is already registered", name)
		}
	} else {
		name = uuid.NewV4().String()
	}
	c.store[name] = conn
	return nil
}

// List returns the list of registered clients name
func (c *ConnStore) List() []string {
	var clist []string
	c.mu.RLock()
	defer c.mu.RUnlock()
	for k := range c.store {
		clist = append(clist, k)
	}
	return clist
}
