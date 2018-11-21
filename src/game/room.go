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
	InCommand  chan *IncommingCommand
}

type IncommingCommand struct {
	Nickname string
	InMsg    IncommingMessage
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
				if m.Status == StatusStartGame || m.Status == StatusGame || m.Status == StatusEndGame {
					m.OwnState = p.State
					RState := PlayerState{}
					for _, templ := range r.Players {
						if templ != p {
							RState = templ.State
						}
					}
					m.RivalState = RState
				}
				p.Send(m)
			}
		case m := <-r.Message:
			m.Player.Send(m.Msg)

		case p := <-r.Register:
			mu := &sync.Mutex{}

			mu.Lock()
			r.Players[p.State.Nickname] = p
			mu.Unlock()
			r.Broadcast <- &Message{Status: StatusInfo, Room: r.ID, Info: "User " + p.State.Nickname + " entered to room"}

		case p := <-r.Unregister:
			delete(r.Players, p.State.Nickname)
			r.Broadcast <- &Message{Status: StatusInfo, Room: r.ID, Info: "User " + p.State.Nickname + " deleted from room"}

		case c := <-r.InCommand:
			r.PerformCommand(c)
		}

	}
}

func (r *Room) Run() {
	r.Ticker = time.NewTicker(time.Second)
	r.GameState("init")
	r.Broadcast <- &Message{Status: StatusStartGame, Room: r.ID, Info: "Starting of Room"}
	for {
		<-r.Ticker.C
		end := r.GameState("next")

		r.Broadcast <- &Message{Status: StatusGame, Room: r.ID}
		if end {
			break
		}
	}
	r.Ticker.Stop()
	r.Broadcast <- &Message{Status: StatusEndGame, Room: r.ID, Info: "Room Stop"}
	r.Stop()
}

func (r *Room) Stop() {
}

func (r *Room) GameState(key string) bool {
	for _, p := range r.Players {
		switch key {
		case "init":
			p.InitPlayerState()
		case "next":
			p.NextPlayerState()
		}
	}

	endflag := false
	if key == "next" {
		endflag = r.ProcessGameState()
	}
	return endflag
}

func (r *Room) ProcessGameState() bool {
	return false
}

func (r *Room) PerformCommand(c *IncommingCommand) {

}
