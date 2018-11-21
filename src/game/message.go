package game

const (
	StatusError     string = "error"
	StatusInfo      string = "info"
	StatusWait      string = "wait"
	StatusStartGame string = "startgame"
	StatusGame      string = "game"
	StatusEndGame   string = "endgame"
)

type Message struct {
	Status     string      `json:"status"`
	Room       string      `json:"room"`
	OwnState   PlayerState `json:"ownstate"`
	RivalState PlayerState `json:"rivalstate"`
	Info       string      `json:"info"`
}

type IncommingMessage struct {
	Command string   `json:"command"`
	Info    string   `json:"info"`
	Pos     Position `json:"pos"`
	MobType string   `json:"mobtype"`
}

type PrivateMessage struct {
	Player *Player
	Msg    *Message
}
