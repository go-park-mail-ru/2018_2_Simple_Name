package game

import (
	"github.com/gorilla/websocket"
)

type PlayerState struct {
	Nickname string      `json:"nickname"`
	HP       int         `json:"hp"`
	Mobs     map[int]Mob `json:"mobs"`
}

type Player struct {
	Room  *Room
	Conn  *websocket.Conn
	State PlayerState
}

func (p *Player) Listen() {

	for {
		msg := &IncommingMessage{}

		err := p.Conn.ReadJSON(msg)
		if websocket.IsUnexpectedCloseError(err) {
			p.Room.Unregister <- p
			p.Conn.Close()
			return
		}
		p.Room.InCommand <- &IncommingCommand{InMsg: *msg, Nickname: p.State.Nickname}
	}
}

func (p *Player) Send(msg *Message) {
	err := p.Conn.WriteJSON(msg)
	if err != nil {
		p.Conn.Close()
		p.Room.Unregister <- p
		return
	}
}

func (p *Player) InitPlayerState() {
	p.State.HP = 100
	//Add some mob to start
}

func (p *Player) NextPlayerState() {
	for _, mob := range p.State.Mobs {
		mob.NextMobState()
	}
}
