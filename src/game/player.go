package game

import (
	"fmt"
	"math/rand"

	"github.com/gorilla/websocket"
)

type PlayerState struct {
	Nickname string          `json:"nickname"`
	HP       int             `json:"hp"`
	Mobs     map[string]*Mob `json:"mobs"`
	Points   int             `json:"points"`
}

type Player struct {
	Room       *Room
	Conn       *websocket.Conn
	State      PlayerState
	Message    chan *IncommingMessage
	Listenflag chan bool
}

func NewPlayer(Nickname string, Conn *websocket.Conn) *Player {
	return &Player{
		State:      PlayerState{Nickname: Nickname, Mobs: make(map[string]*Mob)},
		Conn:       Conn,
		Listenflag: make(chan bool),
		Message:    make(chan *IncommingMessage),
	}
}

func (p *Player) Listen() {

	fmt.Println("Player " + p.State.Nickname + ": Start listening.")

Loop:
	for {
		select {
		case msg := <-p.Message: //читает когда приходит сообщение, асинхронно

			fmt.Println("Player " + p.State.Nickname + ": Translate message to room.")

			p.Room.InCommand <- &IncommingCommand{InMsg: msg, Nickname: p.State.Nickname}
		case flag := <-p.Listenflag:
			switch flag {
			case false:
				p.Room.Unregister <- p
				p.Conn.Close()
				break Loop
			case true:

				fmt.Println("Player " + p.State.Nickname + ": Wait incomming message.")

				go func() {
					msg := &IncommingMessage{}
					err := p.Conn.ReadJSON(msg)
					if websocket.IsUnexpectedCloseError(err) {
						p.Listenflag <- false
						return
					}

					fmt.Println("Player " + p.State.Nickname + ": Get incomming message.")

					p.Message <- msg
					p.Listenflag <- true
				}()
			}
		}
	}

	fmt.Println("Player " + p.State.Nickname + ": End listening.")
}

func (p *Player) Send(msg *Message) {

	// fmt.Println("Send to player "+p.State.Nickname+" ", msg.Status, " message: ")

	err := p.Conn.WriteJSON(msg)
	if err != nil {
		fmt.Println("Error send to player " + p.State.Nickname)
		p.Listenflag <- false
		return
	}
}

func (p *Player) AddMobCommand(Mobtype string) {

	fmt.Println("Player " + p.State.Nickname + ": Perform command addmob " + Mobtype)

	if CheckMobType(Mobtype) {
		mob := CreateMob(Mobtype, GetInitPos(p.Room.OwnTargetParams, p.Room.AreaParams))
		if mob.Price <= p.State.Points {
			p.State.Points -= mob.Price
			p.State.Mobs[mob.ID] = mob
		} else {
			go func() {
				p.Room.Message <- &PrivateMessage{Player: p, Msg: &Message{Status: StatusInfo, Room: p.Room.ID, Info: "Not enough points to buy."}}
			}()
		}
	}
}

func (p *Player) KillMobCommand(pos Position) int {

	fmt.Println("Player " + p.State.Nickname + ": Perform command killmob.")

	killPoints := 0
	for _, mob := range p.State.Mobs {
		if mob.Status != "dead" {
			if mob.CheckKillPos(pos) {
				killPoints += mob.KillPoints
				mob.SetDead()
			}
		}
	}

	fmt.Println("Player "+p.State.Nickname+": killPoints = ", killPoints)

	return killPoints
}

func (p *Player) ProgressState() int {
	// fmt.Println("Player " + p.State.Nickname + ": ProgressState.")

	hpAttack := 0
	for _, mob := range p.State.Mobs {
		switch mob.Status {
		case "run":
			mob.ProgressState(p.Room.RivalTargetParams, p.Room.AreaParams)
			if mob.CheckTargetPos(p.Room.RivalTargetParams) {
				mob.SetAttack()
			}
		case "attack":
			hpAttack += mob.Force
		}
	}
	return hpAttack
}

func (p *Player) IncreasePoints(count int) {
	p.State.Points += count
}

func (p *Player) ReduceHealth(hp int) {
	p.State.HP -= hp
	if p.State.HP < 0 {
		p.State.HP = 0
	}
}

func (p *Player) CheckZHealth() bool {
	return p.State.HP == 0
}

func GetInitPos(target *Target, area *Area) Position {
	y := rand.Intn(int(area.Height))
	return Position{
		X: target.Pos.X - 25,
		Y: float64(y),
	}
}
