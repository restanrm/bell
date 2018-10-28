package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Lister interface {
	List() []string
}

type Registerer interface {
	Register(string, *websocket.Conn) error
}

type ErrorResponse struct {
	Error string
}

// RegisterRequest is the struct that handles registering requests
type RegisterRequest struct {
	Name string `json:"name"`
}

// Register register a client to the websocket handler
func RegisterClients(registerer Registerer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to upgrade the connection")
			return
		}
		// defer conn.Close() is not run because we only want to close on error case

		rr := &RegisterRequest{}
		err = conn.ReadJSON(rr)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to read registerRequest")
			return
		}

		err = registerer.Register(rr.Name, conn)
		if err != nil {
			conn.WriteJSON(ErrorResponse{Error: err.Error()})
			conn.Close()
			return
		}
	}
}

type Getter interface {
	Get(string) (*websocket.Conn, error)
}

type PlayerRequest struct {
	Play string `json:"play"`
}

func PlayOnClient(getter Getter, client, sound string) {
	conn, err := getter.Get(client)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to retrieve connection to client: %v", client)
		return
	}
	err = conn.WriteJSON(PlayerRequest{Play: sound})
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"player": client,
			"sound":  sound,
		}).Errorf("Failed to play sound on player")
	}
}

type ClientsList struct {
	Clients []string `json:"clients"`
}

// List returns the clients list to the caller
func ListClients(l Lister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cl := ClientsList{Clients: l.List()}
		err := json.NewEncoder(w).Encode(cl)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to return clients list")
			return
		}
	}
}
