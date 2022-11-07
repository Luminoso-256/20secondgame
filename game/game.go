package game

import (
	"fmt"
	"image/color"
	_ "image/png"
	"math"
	"math/rand"

	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	USE_TIMECHECK = true
)

var (
	DEBUG_moreThan20Secs = false

	levelPrototype *ebiten.Image
	levelProtoDbg  []*ebiten.Image
)

func (g *Game) Init() {
	g.startTime = time.Now()
	levelPrototype, levelProtoDbg = GenerateLevel()
}

func (g *Game) Update() error {
	var newBalls []Ball
	for _, ball := range g.Balls {
		b := BallMoveTick(ball, g)
		if !b.CanBeRemoved {
			newBalls = append(newBalls, b)
		}
	}
	g.Balls = newBalls

	mx, my := ebiten.CursorPosition()
	dx := math.Abs(float64(mx - g.Player.X))
	dy := math.Abs(float64(my - g.Player.Y))
	netDist := math.Sqrt((dx * dx) + (dy * dy))
	accuracyRadius := (20 * (netDist / 100))

	//biased random point (sqrt random would make it uniform)
	accOffset := accuracyRadius * rand.Float64()
	accTheta := rand.Float64() * 2 * math.Pi
	adjustedPx := float64(mx) + (accOffset * math.Cos(accTheta))
	adjustedPy := float64(my) + (accOffset * math.Sin(accTheta))

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.Balls = append(g.Balls, Ball{
			X:     float64(g.Player.X),
			Y:     float64(g.Player.Y),
			DestX: adjustedPx,
			DestY: adjustedPy,
			Speed: 5,
		})
	}
	//clear
	if len(g.Balls) > 1600 {
		g.Balls = g.Balls[5:]
	}

	if ebiten.IsKeyPressed(ebiten.Key0) {
		levelPrototype, levelProtoDbg = GenerateLevel()
	}

	for _, key := range inpututil.PressedKeys() {
		switch key {
		case ebiten.KeyW:
			g.Player.Y -= 5
		case ebiten.KeyS:
			g.Player.Y += 5
		case ebiten.KeyA:
			g.Player.X -= 5
		case ebiten.KeyD:
			g.Player.X += 5
		}
	}

	if time.Now().Sub(g.startTime).Seconds() > 20 {
		DEBUG_moreThan20Secs = true
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{40, 40, 40, 255})

	//terrain
	for x, r := range g.Level {
		for y, t := range r {
			if t.IsSalted {
				ebitenutil.DrawRect(screen, float64(x)*32, float64(y)*32, 32, 32, color.RGBA{0, 255, 0, 100})
			}
		}
	}

	ebitenutil.DrawRect(screen, float64(g.Player.X), float64(g.Player.Y), 16, 16, color.RGBA{200, 10, 200, 255})
	mx, my := ebiten.CursorPosition()

	//time for math:tm!
	dx := math.Abs(float64(mx - g.Player.X))
	dy := math.Abs(float64(my - g.Player.Y))
	netDist := math.Sqrt((dx * dx) + (dy * dy))
	radius := (20 * (netDist / 100))
	//debug print it
	//ebitenutil.DrawCircle(screen, float64(mx), float64(my), radius, color.RGBA{0, 255, 255, 255})
	radius /= 8
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(radius, radius)
	op.GeoM.Translate(float64(mx)-8*radius, float64(my)-8*radius)
	//ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%v ", radius), mx, my)
	if radius < 3.2 {
		screen.DrawImage(g.Assets.Img["ui/aimmarker_sm2"], op)
	} else if radius < 6 {
		screen.DrawImage(g.Assets.Img["ui/aimmarker_sm3"], op)
	} else {
		screen.DrawImage(g.Assets.Img["ui/aimmarker_sm1"], op)
	}

	//ebitenutil.DrawLine(screen, float64(g.Player.X), float64(g.Player.Y), float64(mx), float64(my), color.RGBA{255, 0, 0, 255})

	for _, ball := range g.Balls {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Translate(ball.X, ball.Y)

		screen.DrawImage(g.Assets.Img["golfball"], op)
		//ebitenutil.DrawCircle(screen, ball.X, ball.Y, 2, color.RGBA{255, 10, 255, 255})

	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %v (ball count %v) | TPS: %v", int(ebiten.CurrentFPS()), len(g.Balls), int(ebiten.CurrentTPS())), 430, 10)

	//== UI ==//
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 10)
	screen.DrawImage(g.Assets.Img["ui/bars/time"], op)
	if time.Now().Sub(g.startTime).Seconds() <= 20 {
		ebitenutil.DrawRect(screen, 48, 22, 80-(time.Now().Sub(g.startTime).Seconds()*4), 8, color.RGBA{76, 76, 129, 255})
	} else {
		screen.DrawImage(g.Assets.Img["ui/bars/timeover"], op)
	}
	op.GeoM.Translate(0, 37)
	screen.DrawImage(g.Assets.Img["ui/bars/score"], op)
	//todo: actual font
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: ##"), 40, 54)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(6, 6)
	screen.DrawImage(levelPrototype, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}
