package game

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type Game struct {
	Rooms      map[string]*Room
	MaxRooms   int
	Connection chan *websocket.Conn
}

type IncomingMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	Player  *Player         `json:"-"`
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func NewGame() *Game {
	return &Game{
		Rooms:      make(map[string]*Room),
		MaxRooms:   2,
		Connection: make(chan *websocket.Conn),
	}
}

func (g *Game) Run() {
	for {
		conn := <-g.Connection
		g.ProcessConn(conn)
	}
}

func (g *Game) ProcessConn(conn *websocket.Conn) {
	id := uuid.NewV4().String()
	p := &Player{
		Conn: conn,
		ID:   id,
	}

	r := g.FindRoom()
	if r == nil {
		return
	}
	r.Players[p.ID] = p
	p.Room = r
	log.Printf("player %s joined room %s", p.ID, r.ID)
	go p.Listen()

	if len(r.Players) == r.MaxPlayers {
		go r.Run()
	}

}

func (g *Game) FindRoom() *Room {
	for _, r := range g.Rooms {
		if len(r.Players) < r.MaxPlayers {
			return r
		}
	}

	if len(g.Rooms) >= g.MaxRooms {
		return nil
	}
	r := NewRoom()
	go r.ListenToPlayers()
	g.Rooms[r.ID] = r
	log.Printf("room %s created", r.ID)

	return r
}
