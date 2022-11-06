package game

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	N = 16
	Z = 20
)

func GenerateLevel() (*ebiten.Image, []*ebiten.Image) {
	var debug []*ebiten.Image
	for i := 0; i < Z; i++ {
		debug = append(debug, ebiten.NewImage(N, N))
		debug[i].Fill(color.White)
	}
	img := ebiten.NewImage(N, N)
	img.Fill(color.White)
	x := rand.Intn(N + 1)
	y := rand.Intn(N + 1)
	for z := 0; z < Z; z++ {

		dir := rand.Intn(4)
		dist := rand.Intn(12) + 4

		nx := x
		ny := y

		switch dir {
		case 0:
			//UP
			ny -= dist
		case 1:
			//DOWN
			ny += dist
		case 2:
			//LEFT
			nx -= dist
		case 3:
			//RIGHT
			nx += dist
		}
		//bounds, with a small amount of bounceback
		if nx > 16 {
			nx = 15
		}
		if ny > 16 {
			ny = 15
		}
		if nx < 0 {
			nx = 1
		}
		if ny < 0 {
			ny = 0
		}
		ebitenutil.DrawLine(debug[z], float64(x), float64(y), float64(nx), float64(ny), color.Black)
		ebitenutil.DrawLine(img, float64(x), float64(y), float64(nx), float64(ny), color.Black)
	}

	return img, debug
}
