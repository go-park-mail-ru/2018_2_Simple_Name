package game

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type Game struct {
	Rooms      map[string]*Room
	MaxRooms   int
	Connection chan *Player
	Mutex      *sync.Mutex
	RoomStat   prometheus.Gauge
}

func NewGame() *Game {
	RoomStat := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "Simple_Game",
		Subsystem: "room_stat",
		Name:      "RoomStat",
		Help:      "Number of running rooms.",
	})
	prometheus.MustRegister(RoomStat)
	mu := &sync.Mutex{}
	return &Game{
		Rooms:      make(map[string]*Room),
		MaxRooms:   10,
		Connection: make(chan *Player),
		Mutex:      mu,
		RoomStat:   RoomStat,
	}
}

func (g *Game) Run() {
	for {
		conn := <-g.Connection
		g.ProcessConn(conn)
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
	fmt.Println("Game: Find room")

	for _, r := range g.Rooms {
		fmt.Println("Game: in for")

		if len(r.Players) < r.MaxPlayers {
			return r
		}
	}
	fmt.Println("Game: after for")

	if len(g.Rooms) >= g.MaxRooms {
		return nil
	}

	fmt.Println("Game: New room")

	g.RoomStat.Inc()
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
	g.RoomStat.Dec()

	fmt.Println("Game: delete room ", r.ID)
}
