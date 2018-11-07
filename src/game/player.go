package game

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PlayerData struct {
	Username string `json:"username"`
	HP       string
	Position Position `json:"position"`
}

type Player struct {
	ID   string
	Room *Room
	Conn *websocket.Conn
	Data PlayerData
}

func (p *Player) Listen() {
	log.Printf("start listening messages from player %s", p.ID)

	for {
		m := &IncomingMessage{}

		err := p.Conn.ReadJSON(m)
		if websocket.IsUnexpectedCloseError(err) {
			log.Printf("player %s was disconnected", p.ID)
			p.Room.Unregister <- p
			return
		}

		m.Player = p
		p.Room.Message <- m
	}
}

func (p *Player) Send(s *Message) {
	err := p.Conn.WriteJSON(s)
	if err != nil {
		log.Printf("cannot send state to client: %s", err)
	}
}

func (r *Room) Run() {
	r.Ticker = time.NewTicker(time.Second)
	go r.RunBroadcast()

	players := []PlayerData{}
	for _, p := range r.Players {
		players = append(players, p.Data)
	}
	state := &State{
		Players: players,
	}

	r.Broadcast <- &Message{Type: "SIGNAL_START_THE_GAME", Payload: state}

	for {
		<-r.Ticker.C
		log.Printf("room %s tick with %d players", r.ID, len(r.Players))

		players := []PlayerData{}
		for _, p := range r.Players {
			players = append(players, p.Data)
		}

		state := &State{
			Players: players,
		}

		r.Broadcast <- &Message{Type: "SIGNAL_NEW_GAME_STATE", Payload: state}
	}
}

func (r *Room) RunBroadcast() {
	for {
		m := <-r.Broadcast
		for _, p := range r.Players {
			p.Send(m)
		}
	}
}
