package game

import "math"

type Player struct {
	X int
	Y int
}

type Ball struct {
	X     float64
	Y     float64
	DestX float64
	DestY float64
	Speed float64
}

func BallMoveTick(b Ball) Ball {
	X_offset := b.DestX - b.X
	Y_offset := b.DestY - b.Y
	//change := b.Speed / math.Sqrt(X_offset*X_offset+Y_offset*Y_offset)
	change := (1 - math.Cos(X_offset*math.Pi+Y_offset*math.Pi)) / 2

	b.X += change * X_offset
	b.Y += change * Y_offset
	return b
}
