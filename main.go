package main

import (
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"minigame/game"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	USE_TIMECHECK = true
)

var (
	DEBUG_moreThan20Secs = false
	//todo: pull in image registry code
	golfball *ebiten.Image
	aimGood  *ebiten.Image
	aimMeh   *ebiten.Image
	aimBad   *ebiten.Image

	levelPrototype *ebiten.Image
	levelProtoDbg  []*ebiten.Image
)

func init() {
	golfball, _, _ = ebitenutil.NewImageFromFile("assets/images/golfball.png")
	aimBad, _, _ = ebitenutil.NewImageFromFile("assets/images/ui/aimmarker_sm1.png")
	aimMeh, _, _ = ebitenutil.NewImageFromFile("assets/images/ui/aimmarker_sm3.png")
	aimGood, _, _ = ebitenutil.NewImageFromFile("assets/images/ui/aimmarker_sm2.png")

	levelPrototype, levelProtoDbg = game.GenerateLevel()
}

type Game struct {
	player game.Player
	balls  []game.Ball
}

func (g *Game) Update() error {
	for x, ball := range g.balls {
		g.balls[x] = game.BallMoveTick(ball)
	}

	mx, my := ebiten.CursorPosition()
	dx := math.Abs(float64(mx - g.player.X))
	dy := math.Abs(float64(my - g.player.Y))
	netDist := math.Sqrt((dx * dx) + (dy * dy))
	accuracyRadius := (20 * (netDist / 100))

	//biased random point (sqrt random would make it uniform)
	accOffset := accuracyRadius * rand.Float64()
	accTheta := rand.Float64() * 2 * math.Pi
	adjustedPx := float64(mx) + (accOffset * math.Cos(accTheta))
	adjustedPy := float64(my) + (accOffset * math.Sin(accTheta))

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.balls = append(g.balls, game.Ball{
			X:     float64(g.player.X),
			Y:     float64(g.player.Y),
			DestX: adjustedPx,
			DestY: adjustedPy,
			Speed: 5,
		})
	}
	//clear
	if len(g.balls) > 512 {
		g.balls = g.balls[30:]
	}

	if ebiten.IsKeyPressed(ebiten.Key0) {
		levelPrototype, levelProtoDbg = game.GenerateLevel()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{40, 40, 40, 255})
	if DEBUG_moreThan20Secs {
		screen.Fill(color.Black)
		ebitenutil.DebugPrint(screen, "Exceeded Time Budget!!")
	}
	ebitenutil.DrawRect(screen, float64(g.player.X), float64(g.player.Y), 16, 16, color.RGBA{200, 10, 200, 255})
	mx, my := ebiten.CursorPosition()

	//time for math:tm!
	dx := math.Abs(float64(mx - g.player.X))
	dy := math.Abs(float64(my - g.player.Y))
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
		screen.DrawImage(aimGood, op)
	} else if radius < 6 {
		screen.DrawImage(aimMeh, op)
	} else {
		screen.DrawImage(aimBad, op)
	}

	//ebitenutil.DrawLine(screen, float64(g.player.X), float64(g.player.Y), float64(mx), float64(my), color.RGBA{255, 0, 0, 255})

	for _, ball := range g.balls {
		// ballDx := math.Abs(ball.DestX - ball.X)
		// ballDy := math.Abs(ball.DestY - ball.Y)
		// ballCx := math.Sqrt((ballDx * ballDx) + (ballDy * ballDy))
		// ballTx := math.Sqrt((ball.DestX * ball.DestX) + (ball.DestY * ball.DestY))

		// progress := ballCx / ballTx

		// //quadratic in line.
		// modifier := (-1 * (math.Pow((2*progress)-1, 2) + 2)) / 2

		// for i := 0.0; i < 1; i += 0.1 {
		// 	res := (-1 * (math.Pow((2*i)-1, 2) + 2)) / 2
		// 	ebitenutil.DrawRect(screen, i*10, res*10, 2, 2, color.Black)
		// }
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(ball.X, ball.Y)
		screen.DrawImage(golfball, op)
		//ebitenutil.DrawCircle(screen, ball.X, ball.Y, 2, color.RGBA{255, 10, 255, 255})

	}
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(4, 4)
	for i := 0; i < 10; i++ {
		screen.DrawImage(levelProtoDbg[i], op)
		op.GeoM.Translate(64, 0)
	}
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(4, 4)
	op.GeoM.Translate(0, 128)
	screen.DrawImage(levelPrototype, op)

	//ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %v (ball count %v) | TPS: %v", ebiten.CurrentFPS(), len(g.balls), ebiten.CurrentTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("20-Second Game")

	//20 second dbg
	if USE_TIMECHECK {
		dbgTimer := time.NewTimer(time.Second * 20)
		go func() {
			<-dbgTimer.C
			DEBUG_moreThan20Secs = true
		}()
	}

	if err := ebiten.RunGame(&Game{
		player: game.Player{
			X: 640 / 2, Y: 480 / 2,
		},
	}); err != nil {
		log.Fatal(err)
	}
}
