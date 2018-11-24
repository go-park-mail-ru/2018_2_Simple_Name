package game

import (
	"fmt"
	"sync"
)

type Game struct {
	Rooms      map[string]*Room
	MaxRooms   int
	Connection chan *Player
	Mutex      *sync.Mutex
}

func NewGame() *Game {
	return &Game{
		Rooms:      make(map[string]*Room),
		MaxRooms:   10,
		Connection: make(chan *Player),
		Mutex:      &sync.Mutex{},
	}
}

func (g *Game) Run() {
	for {
		conn := <-g.Connection
		go g.ProcessConn(conn)
	}
}

func (g *Game) ProcessConn(p *Player) {
	fmt.Println("Game: Process connection")
	r := g.FindRoom()
	if r == nil {
		p.Conn.WriteJSON(Message{Status: StatusError, Info: "All rooms are busy"})
		p.Conn.Close()
		return
	}
	p.Room = r
	r.Register <- p
	r.InitPlayer(p)
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
	go g.FreeRoom(r)

	g.Mutex.Lock()
	g.Rooms[r.ID] = r
	g.Mutex.Unlock()

	return r
}

func (g *Game) FreeRoom(r *Room) {
	<-r.FreeRoom
	delete(g.Rooms, r.ID)

	fmt.Println("Game: delete room ", r.ID)
}
