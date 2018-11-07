package game

import (
	"encoding/json"
	"log"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Room struct {
	ID         string
	Ticker     *time.Ticker
	Players    map[string]*Player
	MaxPlayers int
	Register   chan *Player
	Unregister chan *Player
	Message    chan *IncomingMessage
	Broadcast  chan *Message
	Commands   []*Command
}
type Command struct {
	Player  *Player
	Command string
}
type State struct {
	Players []PlayerData `json:"players"`
}
type NewPlayer struct {
	Username string `json:"username"`
}

func NewRoom() *Room {
	id := uuid.NewV4().String()

	return &Room{
		ID:         id,
		MaxPlayers: 2,
		Players:    make(map[string]*Player),
		Register:   make(chan *Player),
		Unregister: make(chan *Player),
		Broadcast:  make(chan *Message),
		Message:    make(chan *IncomingMessage),
	}
}
func (r *Room) ListenToPlayers() {
	for {
		select {
		case m := <-r.Message:
			log.Printf("message from player %s: %v", m.Player.ID, string(m.Payload))

			switch m.Type {
			case "newPlayer":
				np := &NewPlayer{}
				json.Unmarshal(m.Payload, np)
				m.Player.Data.Username = np.Username
			}

		case p := <-r.Unregister:
			delete(r.Players, p.ID)
			log.Printf("player was deleted from room %s", r.ID)
		}

	}
}
