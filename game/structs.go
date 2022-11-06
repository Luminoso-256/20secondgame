package game

import (
	"math"
	"time"
)

type Player struct {
	X int
	Y int
}

type Ball struct {
	X            float64
	Y            float64
	DestX        float64
	DestY        float64
	Speed        float64
	CanBeRemoved bool
}

type MapTile struct {
	T                int //the type
	SaltingThreshold int
	CurrentSalt      int
	IsSalted         bool
}

type Game struct {
	Player    Player
	Balls     []Ball
	Assets    AssetRegistry
	startTime time.Time
	Level     [32][32]MapTile
}

func BallMoveTick(b Ball, g *Game) Ball {
	X_offset := b.DestX - b.X
	Y_offset := b.DestY - b.Y
	//change := b.Speed / math.Sqrt(X_offset*X_offset+Y_offset*Y_offset)
	change := (1 - math.Cos(X_offset*math.Pi+Y_offset*math.Pi)) / 2

	b.X += change * X_offset
	b.Y += change * Y_offset
	//check for "rest"
	if math.Abs(b.X-b.DestX) < 0.5 && math.Abs(b.Y-b.DestY) < 0.5 {
		//flag a tile and *kill*
		x := int(b.X / 32)
		y := int(b.Y / 32)
		if x >= 0 && x <= 32 {
			if y >= 0 && y <= 32 {
				g.Level[x][y].CurrentSalt += 1
				if g.Level[x][y].CurrentSalt >= g.Level[x][y].SaltingThreshold {
					g.Level[x][y].IsSalted = true
				}
				b.CanBeRemoved = true
			}
		} else {
			b.CanBeRemoved = true //what else are we gonna do with it?
		}
	}
	return b
}
