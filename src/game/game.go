package game

import (
	"sync"
)

type Game struct {
	Rooms      map[string]*Room
	MaxRooms   int
	Connection chan *Player
}

func NewGame() *Game {
	return &Game{
		Rooms:      make(map[string]*Room),
		MaxRooms:   10,
		Connection: make(chan *Player),
	}
}

func (g *Game) Run() {
	for {
		conn := <-g.Connection
		g.ProcessConn(conn)
	}
}

func (g *Game) ProcessConn(p *Player) {

	r := g.FindRoom()
	if r == nil {
		p.Conn.WriteJSON(Message{Type: MsgError, Data: ErrorData{Error: "All rooms are busy"}})
		p.Conn.Close()
		return
	}
	p.Room = r
	r.Register <- p

	if len(r.Players) == r.MaxPlayers {
		go r.Run()
	} else {
		r.Broadcast <- &Message{Type: MsgInfo, Data: InfoData{Status: StatusWait, Room: r.ID}}
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
	go r.RoomManager()

	mu := &sync.Mutex{}
	mu.Lock()
	g.Rooms[r.ID] = r
	mu.Unlock()

	return r
}
