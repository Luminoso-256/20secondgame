package game

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	N = 16
	Z = 20
)

type Street struct {
	sx, sy, ex, ey int
}

func GenerateLevel() (*ebiten.Image, []*ebiten.Image) {
	var debug []*ebiten.Image
	for i := 0; i < Z; i++ {
		debug = append(debug, ebiten.NewImage(N, N))
		debug[i].Fill(color.White)
	}
	img := ebiten.NewImage(N, N)
	img.Fill(color.White)

	rX := 8
	rY := 8

	var avenues []Street
	numAvenues := rand.Intn(3) + 5
	numAvFromRoot := rand.Intn(numAvenues) + 1
	fmt.Printf("total aves: %v | from root: %v\n", numAvenues, numAvFromRoot)
	for i := 0; i < numAvFromRoot; i++ {
		avenues = append(avenues, Street{
			rX, rY, rand.Intn(N + 1), rY + rand.Intn(N+1),
		})
	}
	// for i := 0; i < (numAvenues - numAvFromRoot); i++ {
	// 	toIntersect := avenues[rand.Intn(len(avenues))]
	// 	avenues = append(avenues, Street{
	// 		toIntersect.sy, -1 * toIntersect.sx,
	// 		toIntersect.sy, -1 * toIntersect.ex,
	// 	})
	// }

	for _, street := range avenues {
		ebitenutil.DrawLine(img, float64(street.sx), float64(street.sy), float64(street.ex), float64(street.ey), color.Black)
	}

	return img, debug
}
