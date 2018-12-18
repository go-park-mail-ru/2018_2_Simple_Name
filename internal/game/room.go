package game

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Target struct {
	Pos  Position `json:"pos"`
	Area Area     `json:"area"`
}

type Room struct {
	ID                string
	Ticker            *time.Ticker
	Players           map[string]*Player
	OwnTargetParams   Target
	RivalTargetParams Target
	AreaParams        Area
	MaxPlayers        int
	Register          chan *Player
	Unregister        chan *Player
	Message           chan *PrivateMessage
	Broadcast         chan *Message
	InCommand         chan *IncommingCommand
	Status            string
	FreeManager       chan bool
	FreeRoom          chan bool
	StopRoom          chan bool
	SingleFlag        bool
}

type IncommingCommand struct {
	Nickname string
	InMsg    *IncommingMessage
}

const (
	CommandAddMob  string = "addmob"
	CommandKillMob string = "killmob"
	CommandUpdate  string = "update"
)

type GameState map[string]PlayerState

func NewRoom(SingleFlag bool) *Room {
	id := uuid.NewV4().String()
	return &Room{
		ID:                id,
		MaxPlayers:        2,
		Players:           make(map[string]*Player),
		Register:          make(chan *Player),
		Unregister:        make(chan *Player),
		Message:           make(chan *PrivateMessage),
		Broadcast:         make(chan *Message),
		InCommand:         make(chan *IncommingCommand),
		Status:            "wait",
		FreeManager:       make(chan bool),
		FreeRoom:          make(chan bool),
		StopRoom:          make(chan bool),
		OwnTargetParams:   Target{Pos: Position{X: 75, Y: 300}, Area: Area{Height: 150, Width: 150}},
		RivalTargetParams: Target{Pos: Position{X: 1125, Y: 300}, Area: Area{Height: 150, Width: 150}},
		AreaParams:        Area{Height: 600, Width: 1200},
		SingleFlag:        SingleFlag,
	}
}

func (r *Room) RoomManager() {

	fmt.Println("Started room manager " + r.ID)

Loop:
	for {
		select {
		case m := <-r.Broadcast:

			// fmt.Println("Room Manager " + r.ID + ": Broadcast")

			for _, p := range r.Players {
				if !p.BotFlag {
					if m.Status != StatusWait && m.Status != StatusInfo {
						if r.Status != StatusError {
							m.OwnState = p.State
							m.RivalState = r.GetRival(p).State
						}
					}
					p.Send(m)
				}
			}

		case m := <-r.Message:

			fmt.Println("Room Manager " + r.ID + "send message")
			if !m.Player.BotFlag {
				m.Player.Send(m.Msg)
			}

		case p := <-r.Register:

			fmt.Println("Room Manager " + r.ID + " user " + p.State.Nickname + " enter to the room")

			r.Players[p.State.Nickname] = p

			go func() {
				r.Broadcast <- &Message{Status: StatusInfo, Room: r.ID, Info: "User " + p.State.Nickname + " entered to room"}
			}()

			if r.SingleFlag {
				p2 := GetBot()
				r.InitPlayer(p2)
				r.Players[p2.State.Nickname] = p2
				p2.Room = r
			}
			if len(r.Players) == r.MaxPlayers {
				r.Status = StatusGame
				go r.Run()
			} else {
				r.Status = StatusWait
				go func() {
					r.Broadcast <- &Message{Status: r.Status, Room: r.ID}
				}()
			}

		case p := <-r.Unregister:

			fmt.Println("Room " + r.ID + ": unregister user " + p.State.Nickname)
			if r.Status == StatusGame || r.Status == StatusStartGame {
				r.Status = StatusError
				delete(r.Players, p.State.Nickname)
				go r.Stop()
			} else {

				delete(r.Players, p.State.Nickname)
				if len(r.Players) == 0 || r.SingleFlag {
					go r.Stop()
				} else {
					go func() {
						r.Broadcast <- &Message{Status: StatusInfo, Room: r.ID, Info: "User " + p.State.Nickname + " deleted from room"}
					}()
				}
			}
		case c := <-r.InCommand:
			fmt.Println("Room Manager " + r.ID + " run perform")

			r.PerformCommand(c)

		case <-r.FreeManager:
			break Loop
		}
	}

	fmt.Println("Room Manager " + r.ID + " closed")
}

func (r *Room) Run() {

	fmt.Println("Room " + r.ID + " is running")

	r.Ticker = time.NewTicker(100 * time.Millisecond)
	go func() {
		r.Broadcast <- &Message{Status: StatusStartGame, Room: r.ID, Info: "Starting of Room"}
	}()
Loop:
	for {
		select {
		case <-r.Ticker.C:

			// fmt.Println("Room " + r.ID + ": Game Tic")
			gameover := r.ProgressState()

			if gameover {
				r.Status = StatusGameOver
				fmt.Println("Room " + r.ID + ": Stop ticker")
				r.Ticker.Stop()
				go r.Stop()
			} else {
				go func() {
					r.Broadcast <- &Message{Status: StatusGame, Room: r.ID}
				}()
			}
		case <-r.StopRoom:

			fmt.Println("Room " + r.ID + ": Stop room command.")

			break Loop
		}
	}
}

func (r *Room) Stop() {

	fmt.Println("Room " + r.ID + ": Room Stopping.")

	if r.Status != StatusWait {
		r.StopRoom <- true
	}

	r.Broadcast <- &Message{Status: r.Status, Room: r.ID, Info: "Room Stoped."}
	time.Sleep(time.Millisecond * 100)

	for id, p := range r.Players {
		p.Listenflag <- false
		delete(r.Players, id)
	}
	time.Sleep(time.Millisecond * 2)
	r.FreeManager <- true

	time.Sleep(time.Millisecond * 2)
	r.FreeRoom <- true

	fmt.Println("Room " + r.ID + ": Room Stoped. Your rival ran away!")
}

func (r *Room) ProgressState() bool {

	// fmt.Println("Room " + r.ID + ": ProgressState.")

	keyover := false
	keyNMNM := true
	for _, player := range r.Players {
		hpAttack := player.ProgressState()
		rival := r.GetRival(player)
		rival.ReduceHealth(hpAttack)
		if rival.CheckZHealth() {
			keyover = true
		}
		if !player.CheckNoMobsNoMoney() {
			keyNMNM = false
		}
	}
	if keyNMNM {
		for _, player := range r.Players {
			player.State.HP = 0
		}
		keyover = true
	}
	return keyover
}

func (r *Room) PerformCommand(c *IncommingCommand) {

	fmt.Println("Room " + r.ID + ": PerformCommand " + c.InMsg.Command)

	switch c.InMsg.Command {
	case CommandAddMob:
		r.Players[c.Nickname].AddMobCommand(c.InMsg.CreateMobType)
	case CommandKillMob:
		rival := r.GetRival(r.Players[c.Nickname])
		killPoints := rival.KillMobCommand(c.InMsg.ClickPos)
		r.Players[c.Nickname].IncreasePoints(killPoints)
	}
}

func (r *Room) GetRival(player *Player) *Player {
	for _, rivalPlayer := range r.Players {
		if rivalPlayer != player {
			return rivalPlayer
		}
	}
	return nil
}

func (r *Room) InitPlayer(p *Player) {
	p.State.Points = 150
	p.State.HP = 100

	go p.Listen()
	p.Listenflag <- true

	fmt.Println("Room " + r.ID + ": Init Player.")
}
