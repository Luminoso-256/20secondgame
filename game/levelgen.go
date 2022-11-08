package game

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/shawnridgeway/wfc"
)

const (
	W = 30
	H = 24
	M = 8
	P = 0.125
)

type Street struct {
	sx, sy, ex, ey int
}

func GenerateLevel_WFC() (*ebiten.Image, []*ebiten.Image) {

	data := wfc.SimpleTiledData{
		Unique:    false,
		TileSize:  16,
		Tiles:     []wfc.Tile{},
		Neighbors: []wfc.Neighbor{},
	}

	// Create a new model
	model := wfc.NewSimpleTiledModel(data, 20, 20, false)

	// Run the algorithm until success or contradiction
	outputImg, _ := model.Generate()

	// Run the algorithm 10 times, stopping at success or contradiction
	outputImg, _, _ = model.Iterate(10)

	return ebiten.NewImageFromImage(outputImg), []*ebiten.Image{}
}

func GenerateLevel(g *Game) (*ebiten.Image, []*ebiten.Image) {
	var debug []*ebiten.Image
	var grid [W][H]bool
	var gridSW [W][H]bool
	var gridG [W][H]bool

	img := ebiten.NewImage(W, H)

	for i := 0; i < 4; i++ {
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

		turn := false
		grid[x][y] = true
		for true {

			x1 := false
			x2 := false
			y1 := false
			y2 := false
			if !turn {
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
			}
			if rand.Float32() < P && dist >= M ||
				(dir < 2) && (x1 || x2) ||
				(dir > 1) && (y1 || y2) || turn {

				if dir > 1 {
					dir = rand.Intn(2)
				} else {
					dir = rand.Intn(2) + 2
				}
				turn = false
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
			if !(x < 0) && !(x >= W) && !(y < 0) && !(y >= H) {
				grid[x][y] = true
			}
			dist++
			if x < 0 || x >= W {
				if init < 2 {
					turn = true
				} else {
					break
				}
			}
			if y < 0 || y >= H {
				if init > 1 {
					turn = true
				} else {
					break
				}
			}
		}
	}
	g.levelOverlayLayers[1].Clear()
	g.levelOverlayLayers[0].Clear()

	for x, l := range grid {
		for y, b := range l {
			if b {
				if x < W-2 {
					gridSW[x+1][y] = b
				}
				if x > 0 {
					gridSW[x-1][y] = b
				}
				if y < H-2 {
					gridSW[x][y+1] = b
				}
				if y > 0 {
					gridSW[x][y-1] = b
				}
			}
		}
	}
	for x, l := range gridSW {
		for y, b := range l {
			if b {
				if x < W-2 {
					gridG[x+1][y] = b
				}
				if x > 0 {
					gridG[x-1][y] = b
				}
				if y < H-2 {
					gridG[x][y+1] = b
				}
				if y > 0 {
					gridG[x][y-1] = b
				}
			}
		}
	}

	for x, l := range gridG {
		for y, b := range l {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(2, 2)
			op.GeoM.Translate(float64(x)*32, float64(y)*32)
			if b {
				//the actual level
				g.levelOverlayLayers[0].DrawImage(g.Assets.Img["tile/grass"], op)
				g.Level[x][y].T = 0
				g.Level[x][y].SaltingThreshold = 3
				g.Level[x][y].CurrentSalt = 0
				g.Level[x][y].IsSalted = false
			} else {
				g.levelOverlayLayers[0].DrawImage(g.Assets.Img["tile/temptree"], op)
				g.Level[x][y].T = 10 //building
				g.Level[x][y].SaltingThreshold = 3
				g.Level[x][y].CurrentSalt = 0
				g.Level[x][y].IsSalted = false
			}
		}
	}

	for x, l := range gridSW {
		for y, b := range l {
			if b {
				//the actual level
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(2, 2)
				op.GeoM.Translate(float64(x)*32, float64(y)*32)
				g.levelOverlayLayers[0].DrawImage(g.Assets.Img["tile/dead_grass"], op)
				g.Level[x][y].T = 5 //sidewalk
				g.Level[x][y].SaltingThreshold = 3
				g.Level[x][y].CurrentSalt = 0
				g.Level[x][y].IsSalted = false
			}
		}
	}

	px := 0.
	py := 0.

	for x, l := range grid {
		for y, b := range l {
			if b {

				if px == 0 && py == 0 {
					px = float64(x)
					py = float64(y)
				}

				//the actual level
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(2, 2)
				op.GeoM.Translate(float64(x)*32, float64(y)*32)
				g.levelOverlayLayers[0].DrawImage(g.Assets.Img["tile/cobble"], op)
				img.Set(x, y, color.Black)
				//backend data
				g.Level[x][y].T = 1
				g.Level[x][y].SaltingThreshold = 3
				g.Level[x][y].CurrentSalt = 0
				g.Level[x][y].IsSalted = false
				g.BestPossibleScore++
			}
			// else {
			// 	op := &ebiten.DrawImageOptions{}
			// 	op.GeoM.Scale(2, 2)
			// 	op.GeoM.Translate(float64(x)*32, float64(y)*32)
			// 	g.levelOverlayLayers[0].DrawImage(g.Assets.Img["tile/grass"], op)
			// 	g.Level[x][y].T = 0
			// 	g.Level[x][y].SaltingThreshold = 5
			// 	g.Level[x][y].CurrentSalt = 0
			// 	g.Level[x][y].IsSalted = false
			// }
		}
	}
	g.Player.X = px * 32
	g.Player.Y = py * 32

	return img, debug
}
