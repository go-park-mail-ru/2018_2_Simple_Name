package game

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Mob struct {
	ID    int
	Type  string
	HP    int
	Pos   Position
	Speed int
	Area  int
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
