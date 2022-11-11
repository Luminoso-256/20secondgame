package game

import (
	"fmt"
	"image"
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
	DEBUG_moreThan20Secs  = false
	DEBUG_showLevelLayers = false
	levelPrototype        *ebiten.Image

	emptyImage    = ebiten.NewImage(3, 3)
	emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func (g *Game) Init() {
	g.startTime = time.Now()
	emptyImage.Fill(color.White)
	for i := 0; i < 4; i++ {
		g.levelOverlayLayers = append(g.levelOverlayLayers,
			ebiten.NewImage(32*32, 32*32))
	}
	g.debugOverlay = ebiten.NewImage(32*32, 32*32)
	levelPrototype = GenerateLevel(g)
}

func (g *Game) Update() error {
	g.debugOverlay.Clear()

	var newBalls []Ball
	for _, ball := range g.Balls {
		b := BallMoveTick(ball, g)
		if !b.CanBeRemoved {
			newBalls = append(newBalls, b)
		}
	}
	g.Balls = newBalls
	mx, my := ebiten.CursorPosition()

	mx /= 2
	my /= 2
	mx += g.lastFrameSx
	my += g.lastFrameSy

	px := g.Player.X
	py := g.Player.Y

	//time for math:tm!
	dx := math.Abs(float64(mx) - px)
	dy := math.Abs(float64(my) - py)
	netDist := math.Sqrt((dx * dx) + (dy * dy))

	accuracyRadius := (20 * (netDist / 100))

	accOffset := accuracyRadius * rand.Float64()
	accTheta := rand.Float64() * 2 * math.Pi
	targetX := float64(mx) + (accOffset * math.Cos(accTheta))
	targetY := float64(my) + (accOffset * math.Sin(accTheta))

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.Balls = append(g.Balls, Ball{
			X:     px,
			Y:     py,
			DestX: targetX,
			DestY: targetY,
			Speed: 5,
		})
	}
	//clear
	if len(g.Balls) > 1600 {
		g.Balls = g.Balls[5:]
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		levelPrototype = GenerateLevel(g)
	}

	for _, key := range inpututil.PressedKeys() {
		switch key {
		case ebiten.KeyW:
			g.Player.momY -= 2
		case ebiten.KeyS:
			g.Player.momY += 2
		case ebiten.KeyA:
			g.Player.momX -= 2
		case ebiten.KeyD:
			g.Player.momX += 2
		case ebiten.KeyL:
			DEBUG_showLevelLayers = !DEBUG_showLevelLayers
		}
	}

	gLX := g.Player.X
	gLY := g.Player.Y

	g.Player.X += g.Player.momX
	g.Player.Y += g.Player.momY
	g.Player.X = math.Floor(g.Player.X)
	g.Player.Y = math.Floor(g.Player.Y)
	if g.Player.momX > 10 {
		g.Player.momX = 10
	}
	if g.Player.momX < -10 {
		g.Player.momX = -10
	}
	if g.Player.momY > 10 {
		g.Player.momY = 10
	}
	if g.Player.momY < -10 {
		g.Player.momY = -10
	}
	g.Player.momX -= 0.05 * g.Player.momX
	g.Player.momY -= 0.05 * g.Player.momY

	pTX := int(g.Player.X / 32)
	pTY := int(g.Player.Y / 32)

	if g.Level[pTX][pTY].T == 10 {
		g.Player.momX = 0
		g.Player.momY = 0
		g.Player.X = gLX
		g.Player.Y = gLY
	}

	//TODO!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	//Add A HITBOX system for the love of all that is good

	if g.Player.X < 0 {
		g.Player.momX = 0
		g.Player.X = 0
	}
	if g.Player.Y < 0 {
		g.Player.momY = 0
		g.Player.Y = 0
	}
	if g.Player.X > 960-16 {
		g.Player.momX = 0
		g.Player.X = 960 - 16
	}
	if g.Player.Y > 720-16 {
		g.Player.momY = 0
		g.Player.Y = 720 - 16
	}

	if time.Now().Sub(g.startTime).Seconds() > 20 {
		DEBUG_moreThan20Secs = true
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	//== Titlescreen (State 0)
	if g.GameState == 0 {
		screen.Fill(color.RGBA{40, 40, 40, 255})
		ebitenutil.DebugPrint(screen, "This is going to be a titlescreen")
		ebitenutil.DebugPrintAt(screen, "Press [X] to start", 0, 20)
		return //no other draw logic for title
	}
	//== Gameplay (State 1)

	if DEBUG_showLevelLayers {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.1, 0.1)
		op.GeoM.Translate(0, 32*32/10)
		for _, l := range g.levelOverlayLayers {
			g.debugOverlay.DrawImage(l, op)
			op.GeoM.Translate(0, 32*32/10)
		}
	}
	screen.Fill(color.RGBA{40, 40, 40, 255})
	g.levelOverlayLayers[3].Fill(color.RGBA{15, 15, 15, 255})

	DrawLightCircle(g.levelOverlayLayers[3], g.Player.X+16, g.Player.Y+16, 128, color.RGBA{255, 245, 245, 100})

	for i, layer := range g.levelOverlayLayers {
		if i == 3 {
			op := &ebiten.DrawImageOptions{}
			op.CompositeMode = ebiten.CompositeModeMultiply
			screen.DrawImage(layer, op)
		} else {
			screen.DrawImage(layer, nil)
		}
	}
	//terrain
	for x, r := range g.Level {
		for y, t := range r {
			if t.IsSalted {
				ebitenutil.DrawRect(screen, float64(x)*32, float64(y)*32, 32, 32, color.RGBA{0, 255, 0, 100})
			}
		}
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.Player.X), float64(g.Player.Y))
	screen.DrawImage(g.Assets.Img["player"], op)

	for _, ball := range g.Balls {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Translate(ball.X, ball.Y)

		g.levelOverlayLayers[2].DrawImage(g.Assets.Img["golfball"], op)
	}

	//now, for a magic trick!
	//scrolling: minimal pain (hopefully) edition!

	AP_X_SZ := 960.
	AP_Y_SZ := 720.

	sx := g.Player.X - (AP_X_SZ / 4)
	sy := g.Player.Y - (AP_Y_SZ / 4)
	ex := g.Player.X + (AP_X_SZ / 4)
	ey := g.Player.Y + (AP_Y_SZ / 4)

	if sx < 0 {
		sx = 0
		ex = AP_X_SZ / 2
	}
	if sy < 0 {
		sy = 0
		ey = AP_Y_SZ / 2
	}
	if ex > 960 {
		ex = AP_X_SZ
		sx = AP_X_SZ - (AP_X_SZ / 2)
	}
	if ey > 720 {
		ey = AP_Y_SZ
		sy = AP_Y_SZ - (AP_Y_SZ / 2)
	}

	g.lastFrameSx = int(sx)
	g.lastFrameSy = int(sy)

	subview := ebiten.NewImageFromImage(screen.SubImage(image.Rect(int(sx), int(sy), int(ex), int(ey))))
	op = &ebiten.DrawImageOptions{}

	op.GeoM.Scale(960/AP_X_SZ*2, 720/AP_Y_SZ*2)
	screen.Clear()
	screen.DrawImage(subview, op)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("netX %v netY %v (> %v - %v | ^ %v - %v)", ex-sx, ey-sy, sx, ex, sy, ey), 120, 100)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %v (ball count %v) | TPS: %v", int(ebiten.CurrentFPS()), len(g.Balls), int(ebiten.CurrentTPS())), 430, 10)

	//== UI ==//
	mx, my := ebiten.CursorPosition()

	px := g.Player.X - sx //- sx
	py := g.Player.Y - sy //- sy
	px *= 2
	py *= 2
	//time for math:tm!
	dx := math.Abs(float64(mx) - px)
	dy := math.Abs(float64(my) - py)
	netDist := math.Sqrt((dx * dx) + (dy * dy))
	radius := (20 * (netDist / 100))

	radius /= 8
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(radius, radius)
	op.GeoM.Translate(float64(mx)-8*radius, float64(my)-8*radius)

	if radius < 3.2 {
		screen.DrawImage(g.Assets.Img["ui/aimmarker_sm2"], op)
	} else if radius < 8 {
		screen.DrawImage(g.Assets.Img["ui/aimmarker_sm3"], op)
	} else {
		screen.DrawImage(g.Assets.Img["ui/aimmarker_sm1"], op)
	}

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
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %v (%v)", g.Score, float64(g.Score)/float64(g.BestPossibleScore)), 40, 54)
	screen.DrawImage(g.debugOverlay, nil)

	//gameover is drawn *over* the main game render
	if g.GameState == 2 {
		//TODO
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 960, 720
}
