package game

import (
	"github.com/gorilla/websocket"
)

type PlayerState struct {
	HP   int         `json:"hp"`
	Mobs map[int]Mob `json:"position"`
}

type Player struct {
	Nickname string
	Room     *Room
	Conn     *websocket.Conn
	Data     PlayerState
}

func (p *Player) Listen() {

	for {
		c := &Command{}

		err := p.Conn.ReadJSON(c)
		if websocket.IsUnexpectedCloseError(err) {
			p.Room.Unregister <- p
			p.Conn.Close()
			return
		}
		p.Room.Command <- c
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
	p.Data.HP = 100
	//Add some mob to start
}

func (p *Player) NextPlayerState() {
	for _, mob := range p.Data.Mobs {
		mob.NextMobState()
	}
}
