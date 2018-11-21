package game

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Mob struct {
	ID     int      `json:"id"`
	Type   string   `json:"type"`
	HP     int      `json:"hp"`
	Pos    Position `json:"pos"`
	Speed  int      `json:"speed"`
	Area   int      `json:"-"`
	IsDead bool     `json:"isdead"`
}

func CreateSimpleMob() Mob {
	return Mob{
		Type:  "Simple",
		Speed: 1,
		HP:    1,
		Area:  5,
		Pos: Position{
			X: 20,
			Y: 20,
		},
	}
}

func (m *Mob) NextMobState() {

}
