package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/restanrm/bell/connstore"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Lister interface {
	List() []string
}

type Registerer interface {
	Register(*websocket.Conn) error
}

type ErrorResponse struct {
	Error string
}

// Register register a client to the websocket handler
func RegisterClients(registerer Registerer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		if viper.GetBool("websocket.checkorigin.disabled") {
			upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to upgrade the connection")
			return
		}
		// defer conn.Close() is not run because we only want to close on error case

		err = registerer.Register(conn)
		if err != nil {
			conn.WriteJSON(ErrorResponse{Error: err.Error()})
			conn.Close()
			return
		}
	}
}

type Sender interface {
	Send(string, connstore.MessageType, string) error
}

func PlayOnClient(a Sender, client, sound string) error {
	err := a.Send(client, connstore.Sound, sound)
	if err != nil {
		return errors.Wrap(err, "Failed to send play order to client")
	}
	return nil
}

func SayOnClient(a Sender, client, text string) error {
	err := a.Send(client, connstore.TTS, text)
	if err != nil {
		return errors.Wrap(err, "Failed to send text to speech request to client")
	}
	return nil
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
