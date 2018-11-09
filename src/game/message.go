package game

const (
	MsgError     string = "error"
	MsgGameState string = "gamestate"
	MsgGameEnd   string = "gameend"
	MsgInfo      string = "info"
)

const (
	StatusWait      string = "wait"
	StatusStartGame string = "startgame"
	StatusEndGame   string = "endgame"
	StatusGame      string = "game"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type ErrorData struct {
	Error string
}

type GameEndData struct {
	IsWin bool
}

type InfoData struct {
	Status string
	Room   string
	Msg    string
}

type PrivateMessage struct {
	Player *Player
	Msg    *Message
}
