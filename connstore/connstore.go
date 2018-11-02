package connstore

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/twinj/uuid"

	"github.com/gorilla/websocket"
)

var (
	pongWait   time.Duration = 5 * time.Second
	pingPeriod               = time.Duration((9 * pongWait) / 10)
)

type MessageType int

const (
	Error MessageType = iota
	TTS
	Sound
)

// String convert MessageType to string
func (m MessageType) String() string {
	switch m {
	case 0:
		return "error"
	case 1:
		return "tts"
	case 2:
		return "sound"
	}
	return ""
}

// RegisterRequest is the struct that handles registering requests
type RegisterRequest struct {
	Name string `json:"name"`
}

type RegisterResponse struct {
	Name             string `json:"name"`
	PingPeriodSecond int    `json:"ping_period_seconds"`
}
type PlayerRequest struct {
	// Type is the type of the payload. It can be "error|tts|sound"
	Type string `json:"type"`
	Data string `json:"data"`
}

// ConnStore is the struct that holds clients
type ConnStore struct {
	store     map[string]*client
	mu        sync.RWMutex
	interrupt chan os.Signal
}

type client struct {
	conn *websocket.Conn
	send chan []byte
}

// New return a new Client object
func New() *ConnStore {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	c := &ConnStore{
		store:     make(map[string]*client),
		interrupt: interrupt,
	}
	go c.cleanup()
	return c
}

// function launched to cleanup all websocket connexions if server is closed
func (c *ConnStore) cleanup() {
	for {
		select {
		case <-c.interrupt:
			logrus.Infof("Interrupt signal received, closing all clients connection")
			c.mu.Lock()
			defer c.mu.Unlock()
			for k, v := range c.store {
				logrus.Infof("Closing connection for client %v", k)
				err := v.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, "Server is closing down"))
				if err != nil {
					logrus.WithError(err).Errorf("Failed to send closing message to client %v", k)
					continue
				}
				err = v.conn.Close()
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

// Send get a client and a payload and send content to the destined client
func (c *ConnStore) Send(dest string, t MessageType, data string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	client, ok := c.store[dest]
	if !ok {
		return fmt.Errorf("client %q isn't registered", dest)
	}

	enc, err := json.Marshal(PlayerRequest{Type: t.String(), Data: data})
	if err != nil {
		return errors.Wrapf(err, "Failed to encode request as json")
	}
	client.send <- enc
	return nil
}

// Register function is the public handler to associate new websockets store to the service.Register.
func (c *ConnStore) Register(conn *websocket.Conn) error {
	// read the wanted name from the websocket
	rr := &RegisterRequest{}
	err := conn.ReadJSON(rr)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	name := rr.Name
	_, ok := c.store[name] // if ok, the name is already registered
	if name == "" || ok {
		name = uuid.NewV4().String()
	}
	logrus.Infof("registering new client: %v", name)

	cl := &client{
		conn: conn,
		send: make(chan []byte),
	}
	c.store[name] = cl

	// readpump now looks for disconnect and close
	go func() {
		defer cl.conn.Close()
		conn.SetReadDeadline(time.Now().Add(pongWait))
		conn.SetPongHandler(func(string) error { cl.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
		for {
			_, _, err := cl.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
					logrus.WithError(err).Errorf("Clients closing")
				}
			}
			delete(c.store, name)
			break
		}
	}()

	// writePump
	go func() {
		defer cl.conn.Close()
		tick := time.Tick(pingPeriod)
		for {
			select {
			case data := <-cl.send:
				err := conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					logrus.WithError(err).Errorf("Failed to send message to client")
					return
				}
			case <-tick:
				err := conn.WriteMessage(websocket.PingMessage, nil)
				if err != nil {
					logrus.WithError(err).Errorf("Failed to send ping to client")
					return
				}
			}
		}
	}()

	resp := &RegisterResponse{Name: name}
	err = conn.WriteJSON(resp)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to send register response to client")
		delete(c.store, name)
		return err
	}

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
