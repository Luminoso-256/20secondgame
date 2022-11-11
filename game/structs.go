package game

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Player struct {
	X    float64
	Y    float64
	momX float64
	momY float64
}

type Ball struct {
	X            float64
	Y            float64
	DestX        float64
	DestY        float64
	Speed        float64
	CanBeRemoved bool
	HasSalted    bool
}

type MapTile struct {
	T                int //the type
	SaltingThreshold int
	CurrentSalt      int
	IsSalted         bool
}

type AABB struct {
	x, y, w, h int
}

type Game struct {
	GameState          int
	Player             Player
	Balls              []Ball
	Assets             AssetRegistry
	Score              int
	startTime          time.Time
	Level              [32][32]MapTile
	BestPossibleScore  int
	levelOverlayLayers []*ebiten.Image
	lastFrameSx        int
	lastFrameSy        int
	debugOverlay       *ebiten.Image
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
		if x >= 0 && x < 32 {
			if y >= 0 && y < 32 {
				if g.Level[x][y].T == 10 {
					b.CanBeRemoved = true

				} else if !b.HasSalted && !g.Level[x][y].IsSalted {
					g.Level[x][y].CurrentSalt += 1

					if g.Level[x][y].T == 1 {
						if g.Level[x][y].CurrentSalt >= g.Level[x][y].SaltingThreshold {
							g.Level[x][y].IsSalted = true
							g.Score++
						}
						//add onto overlay 1 (salting)
						ebitenutil.DrawRect(g.levelOverlayLayers[1], float64(x)*32, float64(y)*32, 32, 32, color.RGBA{255, 255, 255, 20})
					} else { //grass becomes mud
						op := &ebiten.DrawImageOptions{}
						op.ColorM.Scale(1, 1, 1, 0.1)
						op.GeoM.Scale(2, 2)
						op.GeoM.Translate(float64(x)*32, float64(y)*32)
						g.levelOverlayLayers[1].DrawImage(g.Assets.Img["tile/dirt"], op)
					}
					//	b.CanBeRemoved = true
					b.HasSalted = true
				}
			}
		} else {
			b.CanBeRemoved = true //what else are we gonna do with it?
		}
	}
	return b
}
