package game

import (
	"sync"
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
	Message    chan *PrivateMessage
	Broadcast  chan *Message
	Command    chan *Command
}

type Command struct {
	Nickname string
	Command  string
	Pos      Position
	Type     string
}

const (
	CommandAddMob  string = "addmob"
	CommandKillMob string = "killmob"
)

type GameState map[string]PlayerState

func NewRoom() *Room {
	id := uuid.NewV4().String()
	return &Room{
		ID:         id,
		MaxPlayers: 2,
		Players:    make(map[string]*Player),
		Register:   make(chan *Player),
		Unregister: make(chan *Player),
		Message:    make(chan *PrivateMessage),
		Broadcast:  make(chan *Message),
	}
}

func (r *Room) RoomManager() {
	for {
		select {
		case m := <-r.Broadcast:
			for _, p := range r.Players {
				p.Send(m)
			}
		case m := <-r.Message:
			m.Player.Send(m.Msg)

		case p := <-r.Register:
			mu := &sync.Mutex{}

			mu.Lock()
			r.Players[p.Nickname] = p
			mu.Unlock()
			r.Broadcast <- &Message{Type: MsgInfo, Data: InfoData{Status: StatusWait, Room: r.ID, Msg: "User " + p.Nickname + " entered to room"}}

		case p := <-r.Unregister:
			delete(r.Players, p.Nickname)
			r.Broadcast <- &Message{Type: MsgInfo, Data: InfoData{Status: StatusWait, Room: r.ID, Msg: "User " + p.Nickname + " deleted from room"}}
		case c := <-r.Command:
			r.PerformCommand(c)
		}

	}
}

func (r *Room) Run() {
	r.Broadcast <- &Message{Type: MsgInfo, Data: InfoData{Status: StatusStartGame, Room: r.ID, Msg: "Starting of Room"}}
	r.Ticker = time.NewTicker(time.Second)
	r.GameState("init")
	for {
		<-r.Ticker.C
		gameState, end := r.GameState("next")

		r.Broadcast <- &Message{Type: MsgGameState, Data: gameState}
		if end {
			break
		}
	}
	r.Ticker.Stop()
	r.Broadcast <- &Message{Type: MsgInfo, Data: InfoData{Status: StatusEndGame, Room: r.ID, Msg: "Room Stop"}}
	r.Stop()
}

func (r *Room) Stop() {
}

func (r *Room) GameState(key string) (GameState, bool) {
	gameState := GameState{}
	for _, p := range r.Players {
		switch key {
		case "init":
			p.InitPlayerState()
		case "next":
			p.NextPlayerState()
		}
		gameState[p.Nickname] = p.Data
	}

	endflag := false
	if key == "next" {
		endflag = r.ProcessGameState()
	}
	return gameState, endflag
}

func (r *Room) ProcessGameState() bool {
	return false
}

func (r *Room) PerformCommand(c *Command) {

}
