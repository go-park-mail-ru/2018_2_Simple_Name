package game

const (
	StatusError     string = "error"
	StatusInfo      string = "info"
	StatusWait      string = "wait"
	StatusStartGame string = "startgame"
	StatusGame      string = "game"
	StatusGameOver  string = "gameover"
)

type Message struct {
	Status     string      `json:"status"`
	Room       string      `json:"room"`
	OwnState   PlayerState `json:"ownstate"`
	RivalState PlayerState `json:"rivalstate"`
	Info       string      `json:"info"`
}

type IncommingMessage struct {
	Command       string   `json:"command"`
	Info          string   `json:"info"`
	ClickPos      Position `json:"clickpos"`
	CreateMobType string   `json:"createmobtype"`
	// OwnTarget     Target   `json:"owntarget"`
	// RivalTarget   Target   `json:"rivaltarget"`
	// Area          Area     `json:"area"`
}

type PrivateMessage struct {
	Player *Player
	Msg    *Message
}
