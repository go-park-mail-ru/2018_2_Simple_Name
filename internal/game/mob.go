package game

import (
	"math"

	uuid "github.com/satori/go.uuid"
)

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
type Area struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}
type Mob struct {
	ID         string   `json:"id"`
	Type       string   `json:"type"`
	HP         int      `json:"hp"`
	Pos        Position `json:"pos"`
	Speed      int      `json:"speed"`
	Area       Area     `json:"area"`
	Force      int      `json:"-"`
	Status     string   `json:"status"`
	KillPoints int      `json:"killpoints"`
	Price      int      `json:"price"`
}

func CreateMob(Mobtype string, Pos Position) *Mob {
	mob := Mob{}
	switch Mobtype {
	case "mob1":
		id := uuid.NewV4().String()
		mob = Mob{
			ID:    id,
			Type:  "mob1",
			HP:    2,
			Speed: 3,
			Force: 3,
			Area:  Area{Width: 50, Height: 50},
			Pos: Position{
				X: Pos.X,
				Y: Pos.Y,
			},
			KillPoints: 10,
			Price:      10,
			Status:     "run",
		}
	case "mob2":
		id := uuid.NewV4().String()
		mob = Mob{
			ID:    id,
			Type:  "mob2",
			HP:    3,
			Speed: 2,
			Force: 2,
			Area:  Area{Width: 50, Height: 50},
			Pos: Position{
				X: Pos.X,
				Y: Pos.Y,
			},
			KillPoints: 20,
			Price:      20,
			Status:     "run",
		}
	case "mob3":
		id := uuid.NewV4().String()
		mob = Mob{
			ID:    id,
			Type:  "mob3",
			HP:    10,
			Speed: 1,
			Force: 1,
			Area:  Area{Width: 50, Height: 50},
			Pos: Position{
				X: Pos.X,
				Y: Pos.Y,
			},
			KillPoints: 30,
			Price:      40,
			Status:     "run",
		}
	}
	return &mob
}

func CheckMobType(Mobtype string) bool {
	Mobtypes := []string{"mob1", "mob2", "mob3"}

	for _, t := range Mobtypes {
		if t == Mobtype {
			return true
		}
	}
	return false
}

func (mob *Mob) CheckKillPos(clickpos Position) bool {
	xcheck := clickpos.X <= mob.Pos.X+mob.Area.Width && clickpos.X >= mob.Pos.X-mob.Area.Width
	ycheck := clickpos.Y <= mob.Pos.Y+mob.Area.Height && clickpos.Y >= mob.Pos.Y-mob.Area.Height
	if xcheck && ycheck {
		return true
	}
	return false
}

func (mob *Mob) CheckTargetPos(tar Target) bool {

	xcheck := mob.Pos.X+mob.Area.Width/2 >= tar.Pos.X-tar.Area.Width/2 && mob.Pos.X <= tar.Pos.X+tar.Area.Width/2
	ycheck := mob.Pos.Y+mob.Area.Width/2 >= tar.Pos.Y-tar.Area.Height/2 && mob.Pos.Y <= tar.Pos.Y+tar.Area.Height/2
	if xcheck && ycheck {
		return true
	}
	return false
}

func (m *Mob) ProgressState(tar Target, area Area) {
	// fmt.Println("old x:", m.Pos.X, " y: ", m.Pos.Y)
	step := 2.
	kx := 0.
	ky := 0.

	if tar.Pos.X-m.Pos.X != 0 {
		ky = (tar.Pos.Y - m.Pos.Y) / math.Abs(tar.Pos.X-m.Pos.X)
	} else {
		ky = tar.Pos.Y - m.Pos.Y
	}
	if ky > 1 {
		ky = ky / math.Abs(ky)
	}
	stepY := ky * step * float64(m.Speed)

	m.Pos.Y += stepY

	if tar.Pos.Y-m.Pos.Y != 0 {
		kx = (tar.Pos.X - m.Pos.X) / math.Abs(tar.Pos.Y-m.Pos.Y)
	} else {
		kx = tar.Pos.X - m.Pos.X
	}
	if kx > 1 {
		kx = kx / math.Abs(kx)
	}

	stepX := kx * step * float64(m.Speed)
	m.Pos.X += stepX
	// fmt.Println("new x:", m.Pos.X, " y: ", m.Pos.Y)
}

func (m *Mob) SetDead() {
	m.Status = "dead"
}

func (m *Mob) SetAttack() {
	m.Status = "attack"
}
