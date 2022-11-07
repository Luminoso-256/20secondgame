package game

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	W = 32
	H = 32
	M = 5
	P = 0.125

	Z = 20
)

type Street struct {
	sx, sy, ex, ey int
}

func GenerateLevel() (*ebiten.Image, []*ebiten.Image) {
	var debug []*ebiten.Image
	var grid [W][H]bool
	for i := 0; i < Z; i++ {
		debug = append(debug, ebiten.NewImage(W, H))
		debug[i].Fill(color.White)
	}
	img := ebiten.NewImage(W, H)
	img.Fill(color.White)

	for i := 0; i < 6; i++ {
		x := 0
		y := 0
		dir := rand.Intn(4)
		dist := 0
		init := dir

		switch init {
		case 0: //Down
			x = rand.Intn(W)
			y = H - 1
		case 1: //Up
			x = rand.Intn(W)
		case 2: //Left
			y = rand.Intn(H)
			x = W - 1
		case 3: //Right
			y = rand.Intn(H)
		}
		if i == 0 {
			//we need a **FOR SURE YOU CAN PLACE PLAYER HERE!!**  point
			x = W / 2
			y = H / 2
		}

		for !(x < 0) && !(x >= W) && !(y < 0) && !(y >= H) {
			grid[x][y] = true

			x1 := false
			x2 := false
			y1 := false
			y2 := false
			if x-1 >= 0 {
				x1 = grid[x-1][y]
			}
			if x+1 < W {
				x2 = grid[x+1][y]
			}
			if y-1 >= 0 {
				y1 = grid[x][y-1]
			}
			if y+1 < H {
				y2 = grid[x][y+1]
			}

			if rand.Float32() < P && dist >= M ||
				(dir < 2) && (x1 || x2) ||
				(dir > 1) && (y1 || y2) {

				if dir > 1 {
					dir = rand.Intn(2)
				} else {
					dir = rand.Intn(2) + 2
				}
				dist = 0
			}
			switch dir {
			case 0:
				y--
			case 1:
				y++
			case 2:
				x--
			case 3:
				x++
			}
			dist++
		}
	}

	for x, l := range grid {
		for y, b := range l {
			if b {
				img.Set(x, y, color.Black)
			}
		}
	}

	return img, debug
}
